package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/screener"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/sender"
	// "github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/util"
)

type Request struct {
	Command string
	Data    Data
}

// Data struct
type Data struct {
	Depth   int
	FadeIn  int
	FadeOut int
	Path    string
	Shuffle bool
	Loop    bool
        Values  []string
	Volume  int
}


func main() {
	// Set up Logger
        //todo: logging path in configuration file
	logger.SetUpLogger("logs/client.log")

	// Start Screen
	screener.StartScreen()

	// check if no argument is given
	if len(os.Args) < 2 {
		logger.Log.Error("Missing required argument")
		return
	}

	// define flags
	command := flag.String("c", "default", "command for the server")
	input := flag.String("i", "", "input music file/folder")
	volume := flag.Int("v", 50, "music volume in percent (default 50)")
	depth := flag.Int("d", 2, "audio file searching depth (default/recommended 2)")
	shuffle := flag.Bool("s", false, "shuffle (default false)")
	loop := flag.Bool("l", false, "loop (default false)")
	fadeIn := flag.Int("fi", 0, "fadein in milliseconds (default 0)")
	fadeOut := flag.Int("fo", 0, "fadeout in milliseconds (default 0)")

	// parsing flags
	logger.Log.Notice("Start Parsing cli parameters")
	flag.Parse()

        var values []string
	// if argument without flagname is given parse it as command
	if flag.NArg() > 1 {
		*command = flag.Arg(0)
                for id, arg := range flag.Args(){
                        if id != 0 {
                                values = append(values, arg)
                        } //end of if
                } // end of for
	} else {
                if flag.NArg() == 1 && *command == "default"{
                        *command = flag.Arg(0)
                } // end of if
	} // end of else

	// check received arguments
	logger.Log.Notice("Check received arguments")
	if *volume < 0 || *depth < 0 || *fadeIn < 0 || *fadeOut < 0 {
		logger.Log.Error(fmt.Errorf("no negative values allowed"))
		return
	}
	if *volume > 100 {
		logger.Log.Info("No volume above 100 allowed")
		*volume = 100
	}

	// print received argument
	logger.Log.Notice("Given arguments:")
	logger.Log.Info("Commabd   " + *command)
	logger.Log.Info("Input:    " + *input)
	logger.Log.Info("Volume:   " + strconv.Itoa(*volume))
	logger.Log.Info("Depth:    " + strconv.Itoa(*depth))
	logger.Log.Info("Shuffle:  " + strconv.FormatBool(*shuffle))
	logger.Log.Info("Loop:     " + strconv.FormatBool(*loop))
	logger.Log.Info("Fade in:  " + strconv.Itoa(*fadeIn))
	logger.Log.Info("Fade out: " + strconv.Itoa(*fadeOut))
	//logger.Log.Info("Tail:     " + flag.Args())

	// parsings songs

	// parsing to json
	logger.Log.Notice("Parsing argument to json")

	dataInfo := &Data{
		Depth:   *depth,
		FadeIn:  *fadeIn,
		FadeOut: *fadeOut,
		Shuffle: *shuffle,
		Loop:    *loop,
		Path:    *input,
                Values:  values,
		Volume:  *volume}
	requestInfo := &Request{
		Command: string(*command),
		Data:    *dataInfo}
	requestJson, _ := json.Marshal(requestInfo)
	logger.Log.Info("JSON String : " + string(requestJson))

	sender.Send(requestJson)

}
