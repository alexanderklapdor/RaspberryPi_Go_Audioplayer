package audiofunctions

//Imports
import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/sender"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
)

//global var definition
var status string = "default"
var needPause bool = false
var needStop bool = false
var stream *portaudio.Stream

// PlayAudio Function
func PlayAudio(fileName string) {
	// set needStop to false (default)
	needStop = false

	// Run after func finished
	defer CallNextSong()
	defer setStatusStop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	util.Check(err)

	// Decode file
	util.Check(decoder.Open(fileName))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	// Initialize Portaudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Set Stream
	out := make([]int16, 8192)
	stream, err = portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	util.Check(err)

	// Start Stream & Set Status
	defer stream.Close()
	util.Check(stream.Start())
	status = "play"
	defer stream.Stop()

	//Loop through the bytes
	for {
		// Check if pause flag is set
		if needPause != false {
			//stop stream
			stream.Stop()
			// set status to pause
			status = "pause"
			// wait while status is pause
			for needPause == true {
				if needStop != false {
					//exit function
					return
				}
			}
			//start stream
			stream.Start()
			//set status to play
			status = "play"
		}
		// check if stop flag is set
		if needStop != false {
			//exit function
			return
		}
		// read byte
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		util.Check(err)

		// write byte to stream
		util.Check(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
		util.Check(stream.Write())
		select {
		case <-sig:
			return
		default:
		}
	} // end of for
} // end of PlayAudio

// set status to Stop
func setStatusStop() {
	status = "stop"
} // end of setStatusStop

// CallNextSong function
func CallNextSong() {
	if needStop != true {
		sender.Send([]byte("{\"Command\":\"next\",\"Data\":{}}"))
	}

} // end of CallNextSong

// StopAudio Function
func StopAudio() {
	needStop = true
} // end of StopAudio

//PauseAudio Function
func PauseAudio() {
	needPause = true
} // end of PauseAudio

// ResumeAudio Function
func ResumeAudio() {
	needPause = false
} // end of ResumeAudio

// SetVolume Function
func SetVolume(volumeValue string) {
	cmd := exec.Command("amixer", "set", "Master", volumeValue+"%")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolume failed with :" + err.Error() + "\n")
	}
} // end of SetVolume

// SetVolumeUp Function
func SetVolumeUp(value string) {
	cmd := exec.Command("amixer", "set", "Master", value+"%+")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolumeUp failed with :" + err.Error() + "\n")
	}
} // end of SetVolumeUp

// SetVolumeDown Function
func SetVolumeDown(value string) {
	cmd := exec.Command("amixer", "set", "Master", value+"%-")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolumeDown failed with :" + err.Error() + "\n")
	}
} //end of SetVolumeDown

// StartPulseaudio Function
func StartPulseaudio() {
	cmd := exec.Command("pulseaudio", "-D")
	err := cmd.Run()
	if err != nil {
		logger.Error("StartPulseaudio failed with :" + err.Error() + "\n")
	}
} // end of StartPulseaudio

// GetVolume Function
func GetVolume() (string, string) {
	// var declaration
	var leftArray, rightArray []string
	var left, right string
	// run Command
	cmd := exec.Command("amixer", "get", "Master")
	cmdOutput, err := cmd.Output()
	// check for errors
	if err != nil {
		logger.Error("GetVolume failed with: " + err.Error() + "\n")
	}
	// regex
	regPerc, _ := regexp.Compile("[[]([0-9]+%)[]]")
	regNumb, _ := regexp.Compile("[0-9]+")
	for _, line := range strings.Split(string(cmdOutput), "\n") {
		if strings.Contains(line, "Left") && strings.Contains(line, "[on]") {
			leftArray = regPerc.FindAllString(string(cmdOutput), 1)
		} // end of if
		if strings.Contains(line, "Right") && strings.Contains(line, "[on]") {
			rightArray = regPerc.FindAllString(string(cmdOutput), 1)
		} // end of if
	} // end of for
	if len(leftArray) != 0 {
		left = leftArray[0]
		left = regNumb.FindAllString(left, 1)[0]
	} else {
		left = "unknown"
	} //end of else
	if len(rightArray) != 0 {
		right = rightArray[0]
		right = regNumb.FindAllString(right, 1)[0]
	} else {
		right = "unknown"
	}
	return left, right
} // end of GetVolume

// GetStatus function
func GetStatus() string {
	return status
}
