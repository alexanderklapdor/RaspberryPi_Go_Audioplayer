// main function
package main

import "flag"
import "fmt"
import "bufio"
import "log"
import "os"
import "strings"
//import "strconv"

func main() {
	fmt.Println("############################################")
	fmt.Println("#            Music Player                  #")
	fmt.Println("############################################")

	// available Arguments (arguments are pointer!) 
	inputFile 	:= flag.String("i", "", "input music file") 
	inputFolder 	:= flag.String("f", "", "input music folder")
	volume 		:= flag.Int("v", 50, 	"music volume in percent (default 50)")
	shuffle 	:= flag.Bool("s", false,"shuffle (default false)")
	fadeIn 		:= flag.Int("fi", 0, 	"fadein in milliseconds (default 0)")
	fadeOut 	:= flag.Int("fo", 0, 	"fadeout in milliseconds (default 0)")

	flag.Parse()

	fmt.Println("File:     ", *inputFile)
	fmt.Println("Folder:   ", *inputFolder)
	fmt.Println("Volume:   ", *volume)
	fmt.Println("Shuffle:  ", *shuffle)
	fmt.Println("Fade in:  ", *fadeIn)
	fmt.Println("Fade out: ", *fadeOut)
	fmt.Println("Tail:     ", flag.Args())

	// check supported formats
	supportedFormats := getSupportedFormats()
	fmt.Println("Supported formats:")
	for _, format := range supportedFormats {
		fmt.Print(" ",format)
	}


	fmt.Println("")
	fmt.Println("********************************************")
	fmt.Println("*             Shutdown                     *")
	fmt.Println("********************************************")
}


func getSupportedFormats() []string{
	// get supported audio formats of 'supportedFormats.cfg' file
	fmt.Println("Get supported audio formats")
	supportedFormats := make([]string, 0)

	// Opening file
	file, err := os.Open("supportedFormats.cfg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if !strings.ContainsAny(line, "#") {
			supportedFormats = append(supportedFormats, line)
			//fmt.Println("format", line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return supportedFormats
}


