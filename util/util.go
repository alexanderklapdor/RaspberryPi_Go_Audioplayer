package util

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
