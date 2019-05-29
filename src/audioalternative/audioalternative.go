package audioalternative

import (
	"fmt"
	"strconv"
)

// global var definition
var status string = "default"
var needPause bool = false
var needStop bool = false
var volume int = 50

// SetStatusStop Function
func setStatusStop() {
	status = "stop"
}

// StopAudio Function
func StopAudio() {
	needStop = true
}

// PauseAudio Function
func PauseAudio() {
	needPause = true
}

// ResumeAudio Function
func ResumeAudio() {
	needPause = false
}

// GetStatus Function
func GetStatus() string {
	return status
}

// GetVolume Function
func GetVolume() (string, string) {
	fmt.Println("GetVolume")
	return strconv.Itoa(volume), strconv.Itoa(volume)
}

// StartPulseaudio Function
func StartPulseaudio() {
	fmt.Println("StartPulseAudio")
}

// SetVolumeDown Function
func SetVolumeDown(value string) {
	fmt.Println("SetVolumeDown")
	SetVolume(strconv.Itoa(volume - 10))
}

// SetVolumeUp Function
func SetVolumeUp(value string) {
	fmt.Println("SetVolumeUp")
	SetVolume(strconv.Itoa(volume + 10))
}

// SetVolume Function
func SetVolume(volumeValue string) {
	fmt.Println("SetVolume")
	volume, _ = strconv.Atoi(volumeValue)
}

// PlayAudio Function
func PlayAudio(fileName string) {
	fmt.Println("PlayAudio")
}

// CallNextSong Function
func CallNextSong() {
	fmt.Println("CallNextSong")
}
