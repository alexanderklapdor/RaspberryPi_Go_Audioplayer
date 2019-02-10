package util

import "github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"

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
		logger.Log.Error(err)
		panic(err)
	}
}
