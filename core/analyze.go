package core

import (
	"math"
	"time"
)

// The Analysis object stores information about many trials.
type Analysis struct {

	// Combined values of each roll in a standard array.
	Rolls map[int]int

	// The number of trials for this analysis. The length of the "Rolls"
	// 	attribute should be equal to this.
	N int

	// The mean value of the trial.
	Mean float64

	// The standard deviation of the trial.
	Deviation float64

	// The standard deviation up and down
	DeviationUp   float64
	DeviationDown float64

	// Timing information, in seconds
	Duration float64
}

/*
Function for determine how often a roll was at least a value
*/
func (a *Analysis) AtLeast(N int) float64 {

	atLeast := 0
	for value, count := range a.Rolls {
		if value >= N {
			atLeast += count
		}
	}
	return float64(atLeast) / float64(a.N)

}

/*
Function for determine how often a roll was at most a value
*/
func (a *Analysis) AtMost(N int) float64 {
	return 1.0 - a.AtLeast(N+1)
}

// Analyze a roll for a given amount of time.
func (d *Definition) AnalyzeTime(duration float64, threads int) *Analysis {
	return d.analyze(0, duration, threads)
}

// Analyze a roll for a given number of rolls.
func (d *Definition) AnalyzeN(N int, threads int) *Analysis {
	return d.analyze(N, 0.0, threads)
}

/*
Basic analyze function. If N is positive, it will just calculate until all rolls
are completed. If N is <= 0, then it will roll until "duration" time in seconds
have passed. (And then process the results of the roll. So in practice, it will
take a little longer than the duration requested.)
*/
func (d *Definition) analyze(N int, duration float64, threads int) *Analysis {

	// Constant for number of rolls to complete in each batch
	batchSize := 1024

	// Keep track of timing
	startTime := time.Now()

	// Determine the number of tests to perform in each thread.
	threadN := N / threads
	N = threadN * threads

	// This is where all the subresults get sent.
	A := &Analysis{
		Rolls: map[int]int{},
		N:     N,
	}

	// We need a couple control channels for handling the different process steps
	totalChannel := make(chan int)
	meanChannel := make(chan float64)
	varianceChannel := make(chan [5]float64)
	mapChannel := make(chan map[int]int)
	stopChannel := make(chan bool)

	/*
		In order to run many of these in parallel, we need to create "thread"
		copies of the original definition, each with it's own random number
		generator. We will dispatch each of these into a goroutine which will
		respond with an analysis object that will be combined later.
	*/
	for i := 0; i < threads; i++ {

		/*
			Each thread's logic is very simple. Just run a crapton of trials and
			record all the results in the analysis object. We will not keep
			track of the standard deviation here as we need to run a post
			analysis afterward.
		*/
		go func() {

			// Create the new Definition object and set a new seed value.
			definition := d.Copy()
			definition = definition.SetRandomSource()

			// Setup a separate channel for continuing processing
			processingChannel := make(chan bool, 1)
			processingChannel <- true

			// Now we just do a shitton of rolls and keep track of the total
			total := 0
			rollMap := map[int]int{}

			/*
				The control loop is the main processing loop. We'll process in
				larger groups to reduce overhead. Plus, processing rolls should
				not take too long in the first place.
			*/
		controlLoop:
			for j := 0; j < threadN/batchSize || N <= 0; j++ {

				// Use the stop channel to determine if time has elapsed
				select {
				case <-stopChannel:
					break controlLoop

				// If the stop channel hasn't been triggered, or we are just processing by
				// 	roll count, we will just continue processing rolls.
				case <-processingChannel:

					for k := 0; k < batchSize; k++ {
						result := definition.Roll()
						total += result.Total
						rollMap[result.Total] += 1
					}

					// We need to trigger the next process loop by sending another bump
					// 	to the processing channel. The outer "j" loop will stop us once
					// 	we've processed enough rolls.
					processingChannel <- true

				}

			}

			// Send the all the rolls, the calculated total, and the total channel and then await the
			// 	calculated mean response.
			mapChannel <- rollMap
			totalChannel <- total
			mean := <-meanChannel

			// Now we calculate the variance for this routine
			variance := 0.0
			varianceUp := 0.0
			upCount := 0.0
			varianceDown := 0.0
			downCount := 0.0
			for k, v := range rollMap {
				deviation := float64(v) * math.Pow(mean-float64(k), 2)
				variance += deviation

				percent := (math.Tanh(float64(k)-mean) + 1.0) / 2.0

				upCount += float64(v) * percent
				downCount += float64(v) * (1.0 - percent)
				varianceUp += deviation * percent
				varianceDown += deviation * (1.0 - percent)
			}

			/*
				Once the variance is calculated, we will send it back through
				the varianceChannel along with the calculated roll map.
			*/
			varianceChannel <- [5]float64{variance, varianceUp, upCount, varianceDown, downCount}

			// Now we safely end as all the processing we need to do has been done.

		}()

	}

	// Send a stop command if a duration is specified
	if N <= 0 {
		<-time.NewTicker(time.Duration(duration * float64(time.Second))).C
		for i := 0; i < threads; i++ {
			stopChannel <- true
		}
	}

	// Now we need to construct the total heatmap of results
	totalRolls := 0
	for j := 0; j < threads; j++ {
		M := <-mapChannel
		for k, v := range M {
			A.Rolls[k] += v
			totalRolls += v
		}
	}

	// As the completed results come in, we want to create a meta-analysis object which
	// 	combines all the information into one place. First, we will analyze and send back the mean
	total := float64(0)
	for j := 0; j < threads; j++ {
		total += float64(<-totalChannel)
	}
	A.Mean = total / float64(totalRolls)

	// We now send the means back to all the individual channels
	for j := 0; j < threads; j++ {
		meanChannel <- A.Mean
	}

	// Now we need to extract the total variance
	variance := 0.0
	varianceUp := 0.0
	upCount := 0.0
	varianceDown := 0.0
	downCount := 0.0
	for j := 0; j < threads; j++ {
		variances := <-varianceChannel
		variance += variances[0]
		varianceUp += variances[1]
		upCount += variances[2]
		varianceDown += variances[3]
		downCount += variances[4]
	}
	A.Deviation = math.Pow(variance/float64(totalRolls), 0.5)
	A.DeviationUp = math.Pow(varianceUp/upCount, 0.5)
	A.DeviationDown = math.Pow(varianceDown/downCount, 0.5)

	A.N = totalRolls
	A.Duration = float64(time.Since(startTime)) / float64(time.Second)

	// Once everything is aggregated, return the completed analysis object.
	return A

}
