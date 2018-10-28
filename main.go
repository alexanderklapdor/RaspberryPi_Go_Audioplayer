// main function
package main

import "flag"
import "fmt"
import "bufio"
import "os"
import "strings"
import "path/filepath"
import "path"
import "io/ioutil"
import "github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"

import "strconv"

func main() {

	// Set up Logger
	logger.SetUpLogger()
	fmt.Println("############################################")
	fmt.Println("#            Music Player                  #")
	fmt.Println("############################################")

	// available Arguments (arguments are pointer!)
	input := flag.String("i", "", "input music file/folder")
	volume := flag.Int("v", 50, "music volume in percent (default 50)")
	depth := flag.Int("d", 2, "audio file searching depth (default/recommended 2)")
	shuffle := flag.Bool("s", false, "shuffle (default false)")
	fadeIn := flag.Int("fi", 0, "fadein in milliseconds (default 0)")
	fadeOut := flag.Int("fo", 0, "fadeout in milliseconds (default 0)")

	logger.Log.Notice("Start Parsing cli parameters")
	flag.Parse()

	logger.Log.Info("Input:    " + *input)
	logger.Log.Info("Volume:   " + strconv.Itoa(*volume))
	logger.Log.Info("Depth:    " + strconv.Itoa(*depth))
	fmt.Println("Shuffle:  " , *shuffle)
	logger.Log.Info("Fade in:  " + strconv.Itoa(*fadeIn))
	logger.Log.Info("Fade out: " + strconv.Itoa(*fadeOut))
	//logger.Log.Info("Tail:     " + flag.Args())

	// check supported formats
	logger.Log.Notice("Parsing supported formats")
	supportedFormats := getSupportedFormats()
	formatString := "  "
	for _, format := range supportedFormats {
		formatString = formatString + format + ", "
	}
	formatString = formatString[:len(formatString)-2] // todo: Exception handling
	logger.Log.Info("Supported formats: " + formatString)

	// check if given file/folder exists
	logger.Log.Notice("Check if folder/(file exists")
	fi, err := os.Stat(*input)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		// directory given
		logger.Log.Info("Directory found")
		logger.Log.Notice("Getting files inside of the folder")
		fileList := getFilesInFolder(*input, supportedFormats, *depth)
		logger.Log.Notice("Supported Files: ")
		for _, fileElement := range fileList {
			logger.Log.Info(fileElement)
		}
	case mode.IsRegular():
		// file given
		logger.Log.Info("File found")
		var extension = filepath.Ext(*input)
		if stringInArray(extension, supportedFormats) {
			logger.Log.Notice("Extension supported")
		} else {
			logger.Log.Warning("Extension not supported")
		}
	}

	fmt.Println("")
	fmt.Println("********************************************")
	fmt.Println("*             Shutdown                     *")
	fmt.Println("********************************************")
}

func getFilesInFolder(folder string, supportedExtensions []string, depth int) []string {
	// fmt.Println("get files in ", folder)
	fileList := make([]string, 0)
	if depth > 0 {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			logger.Log.Error(err)
		}
		for _, file := range files {
			filename := joinPath(folder, file.Name())

			fi, err := os.Stat(filename)
			if err != nil {
				logger.Log.Error(err)
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
	supportedFormats := make([]string, 0)

	// Opening file
	file, err := os.Open("supportedFormats.cfg")
	if err != nil {
		logger.Log.Error(err)
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
		logger.Log.Error(err)
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
