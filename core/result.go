package core

/*
The Result struct holds the information about the rolls that went into producing
the ending total. If "Results" has a length greater than 0, then this is a
parent result, and this is not a base case.
*/
type Result struct {

	// Indicates whether or not this is a base roll or not
	Base bool `json:"base"`

	/*
		An array of all sub-results. This is only applicable if not the
		base-case. During the base case, this will be a nil value.
	*/
	Results []*Result `json:"results"`

	/*
		A list of values which contribute to the total. This is useful during
		cases where more rolls happen than those that actually count. An example
		would be during an advantage roll, two rolls happen, but only one is
		counted and added to the "Values" array.
	*/
	Values []int `json:"values"`

	/*
		The total value of all the integers in "Values". Sometimes this includes
		an additional operation to be performed outside of the scope of the
		rolls.
	*/
	Total int `json:"total"`
}
