package audiofunctions

//Imports
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

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
	fmt.Println("Playing audiofile.  Press Ctrl-C to stop.")

	defer CallNextSong()
	defer setStatusStop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	util.Check(err)

	util.Check(decoder.Open(fileName))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	stream, err = portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	util.Check(err)
	// todo: here call next song

	defer stream.Close()
	util.Check(stream.Start())
	status = "play"
	defer stream.Stop()
	for {
		if needPause != false {
			stream.Stop()
			status = "pause"
			for needPause == true {
				if needStop != false {
					return
				}
			}
			stream.Start()
			status = "play"
		}
		if needStop != false {
			needStop = false
			return
		}
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		util.Check(err)

		util.Check(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
		util.Check(stream.Write())
		select {
		case <-sig:
			return
		default:
		}
	} // end of for
} // end of PlayAudio

func setStatusStop() {
	status = "stop"
}

func CallNextSong() {
	sender.Send([]byte("{\"Command\":\"next\",\"Data\":{}}"))
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
		logger.Log.Error("SetVolume failed with " + err.Error() + "\n")
	}
} // end of SetVolume

// SetVolumeUp Function
func SetVolumeUp() {
	cmd := exec.Command("amixer", "set", "Master", "2%+")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("SetVolumeUp failed with " + err.Error() + "\n")
	}
} // end of SetVolumeUp

// SetVolumeDown Function
func SetVolumeDown() {
	cmd := exec.Command("amixer", "set", "Master", "2%-")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("SetVolumeDown failed with " + err.Error() + "\n")
	}
} //end of SetVolumeDown

// StartPulseaudio Function
func StartPulseaudio() {
	cmd := exec.Command("pulseaudio", "-D")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("StartPulseaudio failed with " + err.Error() + "\n")
	}
} // end of StartPulseaudio

func GetStatus() string {
	return status
}
