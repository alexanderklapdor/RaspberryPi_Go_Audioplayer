package audioalternative

import (
	"fmt"
	"strconv"
)

//global var definition
var status string = "default"
var needPause bool = false
var needStop bool = false
var volume int = 50

func setStatusStop() {
	status = "stop"
}

// StopAudio Function
func StopAudio() {
	needStop = true
} // end of StopAudio

//PauseAudio Function
func PauseAudio() {
	needPause = true
} // end of PauseAudio

// ResumeAudio Functi
func ResumeAudio() {
	needPause = false
} // end of ResumeAudio

func GetStatus() string {
	return status
}

func GetVolume() (string, string) {
	fmt.Println("GetVolume")
	return strconv.Itoa(volume), strconv.Itoa(volume)
}
func StartPulseaudio() {
	fmt.Println("StartPulseAudio")
}
func SetVolumeDown(value string) {
	fmt.Println("SetVolumeDown")
	SetVolume(strconv.Itoa(volume - 10))
}
func SetVolumeUp(value string) {
	fmt.Println("SetVolumeUp")
	SetVolume(strconv.Itoa(volume + 10))
}
func SetVolume(volumeValue string) {
	fmt.Println("SetVolume")
	volume, _ = strconv.Atoi(volumeValue)
}
func PlayAudio(fileName string) {
	fmt.Println("PlayAudio")
}
func CallNextSong() {
	fmt.Println("CallNextSong")
}
