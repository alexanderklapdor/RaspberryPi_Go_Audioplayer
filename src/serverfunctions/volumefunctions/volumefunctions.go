package volumefunctions

import (
	"strconv"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/portaudiofunctions"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/structs"
)

// SetVolume Function
func SetVolume(data structs.Data) string {
	var volume string
	if len(data.Values) != 0 {
		volume = data.Values[0]
	} else {
		volume = strconv.Itoa(data.Volume)
	} // end of else
	portaudiofunctions.SetVolume(volume)
	logger.Info("Executing: Set volume to " + volume)
	return "Set volume to " + volume
} // end of setVolume

// IncreaseVolume Function
func IncreaseVolume() string {
	logger.Info("Executing: Increase volume")
	portaudiofunctions.SetVolumeUp("10")
	return "Increased volume by 10 \n" + PrintVolume()
} // end of increaseVolume

// DecreaseVolume function
func DecreaseVolume() string {
	logger.Info("Executing: Decrease volume")
	portaudiofunctions.SetVolumeDown("10")
	return "Decreased volume by 10 \n" + PrintVolume()
} // end of decreaseVolume

// getVolume function
func GetVolume() (string, string) {
	left, right := portaudiofunctions.GetVolume()
	return left, right
} // end of getVolume()

// printVolume function
func PrintVolume() string {
	left, right := GetVolume()
	return "Current Volume:  Left(" + left + ")  Right(" + right + ")"
}
