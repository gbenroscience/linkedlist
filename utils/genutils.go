package utils

import (
	"math/rand"
	"sync"
	"time"
)

// RandomLife ...
type RandomLife struct {
	SeededRand *rand.Rand
}

// Letters of the alphabet in upper and lower case
const (
	ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DIGITS   = "0123456789"
)

// NewRnd ...
func NewRnd() RandomLife {
	return RandomLife{
		SeededRand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

//NextInt - Generates a number between 0 and max, max. excluded
func (rnd *RandomLife) NextInt(max int) int {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := 0 + rnd.SeededRand.Intn(max)
	mu.Unlock()

	return i

}

//NextInt - Generates a number between 0 and max, max. excluded
func (rnd *RandomLife) NextBool() bool {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := 0 + rnd.SeededRand.Intn(2)
	mu.Unlock()

	return i == 1

}

//NextFloat - Generates a number between 0 and 1
func (rnd *RandomLife) NextFloat() float64 {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := 0 + rnd.SeededRand.Float64()
	mu.Unlock()

	return i

}

//GetArrEntryRndInt -
func (rnd *RandomLife) GetArrEntryRndInt(arr []int) int {

	var mu sync.Mutex

	// lock/unlock when accessing the rand from a goroutine
	mu.Lock()
	i := rnd.SeededRand.Intn(len(arr))

	name := arr[i]
	mu.Unlock()

	return name
}

// GenerateRndFloat ...Supply min and max
func (rnd *RandomLife) GenerateRndFloat(min float32, max float32) float32 {
	return min + rnd.SeededRand.Float32()*(max-min)
}

// GenerateRndArray ...Generates an array of random integers. The maximum number producible is maxElem
// If the signedElems parameter is true, the minimum number producible is -maxElems.
// When signedElems is true, the elements of the array are a random mix of positive and negative integers
func (rnd *RandomLife) GenerateRndArray(arraySize int, maxElem int, signedElems bool) []int {

	arr := make([]int, arraySize)

	for i := 0; i < arraySize; i++ {
		val := rnd.NextInt(maxElem)
		if signedElems {
			if !rnd.NextBool() {
				val *= -1
			}
		}
		arr[i] = val
	}
	return arr
}
