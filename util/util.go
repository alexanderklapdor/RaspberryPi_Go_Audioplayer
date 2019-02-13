package util

import (
	"math/rand"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

// check if string is element of the array
func StringInArray(str string, list []string) bool {
	// check if string is in string-list
	for _, element := range list {
		if element == str {
			return true
		}
	}
	return false
}

// error check
func Check(err error) {
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}

// Shuffel String-Array
func Shuffle(array []string) []string {
	// Create new random variable based on the current time
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// swap elemnts random for each string position
	for n := len(array); n > 0; n-- {
		randIndex := r.Intn(n)
		array[n-1], array[randIndex] = array[randIndex], array[n-1]
	}
	//Return array
	return array
}
