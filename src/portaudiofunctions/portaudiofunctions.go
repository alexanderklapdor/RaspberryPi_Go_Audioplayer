package portaudiofunctions

// Imports
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

// Global var definition
var status string = "default"
var needPause bool = false
var needStop bool = false
var stream *portaudio.Stream

// PlayAudio Function
func PlayAudio(fileName string) {
	// Set needStop to false (default)
	needStop = false

	// Run after func finished
	defer CallNextSong()
	defer setStatusStop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// Create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	util.Check(err, "Server")

	// Decode file
	util.Check(decoder.Open(fileName), "Server")
	defer decoder.Close()

	// Get audio format information
	rate, channels, _ := decoder.GetFormat()

	// Make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	// Initialize Portaudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Set Stream
	out := make([]int16, 8192)
	stream, err = portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	util.Check(err, "Server")

	// Start Stream & Set Status
	defer stream.Close()
	util.Check(stream.Start(), "Server")
	status = "play"
	defer stream.Stop()

	//Loop through the bytes
	for {
		// Check if pause flag is set
		if needPause != false {
			// Stop stream
			stream.Stop()
			// Set status to pause
			status = "pause"
			// Wait while status is pause
			for needPause == true {
				if needStop != false {
					// Exit function
					return
				}
			}
			// Start stream
			stream.Start()
			// Set status to play
			status = "play"
		}
		// Check if stop flag is set
		if needStop != false {
			// Exit function
			return
		}
		// Read byte
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		util.Check(err, "Server")

		// Write byte to stream
		util.Check(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out), "Server")
		util.Check(stream.Write(), "Server")
		select {
		case <-sig:
			return
		default:
		}
	} // End of for
} // End of PlayAudio

// Set status to Stop
func setStatusStop() {
	status = "stop"
} // End of setStatusStop

// CallNextSong function
func CallNextSong() {
	if needStop != true {
		sender.Send([]byte("{\"Command\":\"next\",\"Data\":{}}"))
	}

} // End of CallNextSong

// StopAudio Function
func StopAudio() {
	needStop = true
} // End of StopAudio

//PauseAudio Function
func PauseAudio() {
	needPause = true
} // End of PauseAudio

// ResumeAudio Function
func ResumeAudio() {
	needPause = false
} // End of ResumeAudio

// SetVolume Function
func SetVolume(volumeValue string) {
	cmd := exec.Command("amixer", "set", "Master", volumeValue+"%")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolume failed with :" + err.Error())
	}
} // End of SetVolume

// SetVolumeUp Function
func SetVolumeUp(value string) {
	cmd := exec.Command("amixer", "set", "Master", value+"%+")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolumeUp failed with :" + err.Error())
	}
} // End of SetVolumeUp

// SetVolumeDown Function
func SetVolumeDown(value string) {
	cmd := exec.Command("amixer", "set", "Master", value+"%-")
	err := cmd.Run()
	if err != nil {
		logger.Error("SetVolumeDown failed with :" + err.Error())
	}
} //End of SetVolumeDown

// StartPulseaudio Function
func StartPulseaudio() {
	cmd := exec.Command("pulseaudio", "-D")
	err := cmd.Run()
	if err != nil {
		logger.Error("StartPulseaudio failed with :" + err.Error())
	}
} // End of StartPulseaudio

// StopPulseaudio Function
func StopPulseaudio() {
	cmd := exec.Command("pulseaudio", "--kill")
	err := cmd.Run()
	if err != nil {
		logger.Error("StopPulseaudio failed with :" + err.Error())
	}
} // End of StopPulseaudio

// GetVolume Function
func GetVolume() (string, string) {
	// Var declaration
	var leftArray, rightArray []string
	var left, right string
	// Run Command
	cmd := exec.Command("amixer", "get", "Master")
	cmdOutput, err := cmd.Output()
	// Check for errors
	if err != nil {
		logger.Error("GetVolume failed with: " + err.Error())
	}
	// Regex Query - Get volume levels
	regPerc, _ := regexp.Compile("[[]([0-9]+%)[]]")
	regNumb, _ := regexp.Compile("[0-9]+")
	for _, line := range strings.Split(string(cmdOutput), "\n") {
		if strings.Contains(line, "Left") && strings.Contains(line, "[on]") {
			leftArray = regPerc.FindAllString(string(cmdOutput), 1)
		} // End of if
		if strings.Contains(line, "Right") && strings.Contains(line, "[on]") {
			rightArray = regPerc.FindAllString(string(cmdOutput), 1)
		} // End of if
	} // End of for
	if len(leftArray) != 0 {
		left = leftArray[0]
		left = regNumb.FindAllString(left, 1)[0]
	} else {
		left = "unknown"
	} //End of else
	if len(rightArray) != 0 {
		right = rightArray[0]
		right = regNumb.FindAllString(right, 1)[0]
	} else {
		right = "unknown"
	}
	return left, right
} // End of GetVolume

// GetStatus function
func GetStatus() string {
	return status
}
