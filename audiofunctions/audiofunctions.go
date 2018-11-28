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
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
)

// PlayAudio Function
func PlayAudio(fileName string) {
	fmt.Println("Playing audiofile.  Press Ctrl-C to stop.")

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
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	util.Check(err)
	defer stream.Close()

	util.Check(stream.Start())
	defer stream.Stop()
	for {
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
	}
}

// StopAudio Function
func StopAudio() {

}

// EndAudio Function
func EndAudio(fileList []string, printFiles bool) {

}

// SetVolume Function
func SetVolume(volumeValue string) {
	cmd := exec.Command("pactl", "set-sink-volume", "0", volumeValue+"%")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("cmd.Run() failed with %s\n", err)
	}
}

// SetVolume Function
func SetVolumeUp() {
	cmd := exec.Command("pactl", "set-sink-volume", "0", "+1%")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("cmd.Run() failed with %s\n", err)
	}
}

// SetVolume Function
func SetVolumeDown() {
	cmd := exec.Command("pactl", "set-sink-volume", "0", "-1%")
	err := cmd.Run()
	if err != nil {
		logger.Log.Error("cmd.Run() failed with %s\n", err)
	}
}

// StartPulseaudio
func StartPulseaudio() {
	cmd := exec.Command("pulseaudio", "-D")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log.Error("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
