// main function
package main

import "flag"
import "fmt"
import "bufio"
import "log"
import "os"
import "strings"
import "path/filepath"
import "path"
import "io/ioutil"
import "./logger"

//import "strconv"

func main() {

	// Set up Logger
	logger.SetUpLogger()
	fmt.Println("############################################")
	fmt.Println("#            Music Player                  #")
	fmt.Println("############################################")

	// available Arguments (arguments are pointer!)
	input := flag.String("i", "", "input music file/folder")
	volume := flag.Int("v", 50, "music volume in percent (default 50)")
	shuffle := flag.Bool("s", false, "shuffle (default false)")
	fadeIn := flag.Int("fi", 0, "fadein in milliseconds (default 0)")
	fadeOut := flag.Int("fo", 0, "fadeout in milliseconds (default 0)")

	flag.Parse()

	fmt.Println("Input:    ", *input)
	fmt.Println("Volume:   ", *volume)
	fmt.Println("Shuffle:  ", *shuffle)
	fmt.Println("Fade in:  ", *fadeIn)
	fmt.Println("Fade out: ", *fadeOut)
	fmt.Println("Tail:     ", flag.Args())

	// check supported formats
	supportedFormats := getSupportedFormats()
	fmt.Println("Supported formats:")
	for _, format := range supportedFormats {
		fmt.Print(" ", format)
	}
	fmt.Println()

	// check if file exists
	fi, err := os.Stat(*input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("input given")
	switch mode := fi.Mode(); {
	case mode.IsDir():
		fmt.Println("Found directory")
		fileList := getFilesInFolder(*input, supportedFormats, 2)
		fmt.Println("Supported Files: ")
		for _, fileElement := range fileList {
			fmt.Println(fileElement)
		}
	case mode.IsRegular():
		fmt.Println("Found file")
		var extension = filepath.Ext(*input)
		if stringInArray(extension, supportedFormats) {
			fmt.Println("Extension supported")
		} else {
			fmt.Println("Extension not supported")
		}
	}

	fmt.Println("")
	fmt.Println("********************************************")
	fmt.Println("*             Shutdown                     *")
	fmt.Println("********************************************")
}

func getFilesInFolder(folder string, supportedExtensions []string, depth int) []string {
	fmt.Println("get files in ", folder)
	fileList := make([]string, 0)
	if depth > 0 {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			filename := joinPath(folder, file.Name())

			fi, err := os.Stat(filename)
			if err != nil {
				fmt.Println(err)
			}

			switch mode := fi.Mode(); {
			case mode.IsDir():
				newFolder := filename + "/"
				newFiles := getFilesInFolder(newFolder, supportedExtensions, depth-1)
				for _, newFile := range newFiles {
					fileList = append(fileList, newFile)
				}
			case mode.IsRegular():
				var extension = filepath.Ext(filename)
				if stringInArray(extension, supportedExtensions) {
					fileList = append(fileList, filename)
				}
			}
		}
	} else {
		//fmt.Println("Max depth reached")
	}
	return fileList
}

func joinPath(source, target string) string {
	if path.IsAbs(target) {
		return target
	}
	return path.Join(path.Dir(source), target)
}

func getSupportedFormats() []string {
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

func stringInArray(str string, list []string) bool {
	// check if string is in string-list
	for _, element := range list {
		if element == str {
			return true
		}
	}
	return false
}
