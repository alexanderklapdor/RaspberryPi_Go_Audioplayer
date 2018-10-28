// main function
package main

import "flag"
import "fmt"

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


	fmt.Println("********************************************")
	fmt.Println("*             Shutdown                     *")
	fmt.Println("********************************************")
}
