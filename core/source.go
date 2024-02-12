package core

import (
	"math/rand"
	"time"
)

/*
This is the source for the whole package. You can also set the source for any
specific roll if your application calls for it by utilizing the
"(*Definition).WithSeed()" or "(*Definition).WithSource()" functions.

You can additionally set the default seed for the package by using the
"dice.SetSeed()" or "dice.SetSource()" functions.

Keep in mind that Sources are not thread-safe in golang, and if you want to do
large amounts of processing on several different Sources in parallel, you should
generate random Sources for them using the "(*Definition).RandomSource()"
function.
*/
var defaultSource = rand.New(rand.NewSource(time.Now().UnixNano()))

/////////////
// GLOBALS //
/////////////

// Set the defaultSource based on an in seed.
func SetSeed(s int64) {
	lock.Lock()
	defer lock.Unlock()
	defaultSource = rand.New(rand.NewSource(s))
}

// Set the defaultSource based on a provide *rand.Rand Source.
func SetSource(s *rand.Rand) {
	lock.Lock()
	defer lock.Unlock()
	defaultSource = s
}

// Set the defaultSource randomly.
func SetRandomSource() {
	lock.Lock()
	defer lock.Unlock()
	defaultSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

/////////////////
// DEFINITIONS //
/////////////////

// Set the defaultSource based on an in seed.
func (d *Definition) SetSeed(s int64) *Definition {
	d.source = rand.New(rand.NewSource(s))
	for _, child := range d.Children {
		child.SetSource(d.source)
	}
	return d
}

// Set the defaultSource based on a provide *rand.Rand Source.
func (d *Definition) SetSource(s *rand.Rand) *Definition {
	d.source = s
	for _, child := range d.Children {
		child.SetSource(d.source)
	}
	return d
}

// Set the defaultSource randomly.
func (d *Definition) SetRandomSource() *Definition {
	d.source = rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, child := range d.Children {
		child.SetSource(d.source)
	}
	return d
}

// Set the defaultSource randomly.
func (d *Definition) SetDefaultSource() *Definition {
	lock.Lock()
	defer lock.Unlock()
	d.source = defaultSource
	for _, child := range d.Children {
		child.SetSource(d.source)
	}
	return d
}
