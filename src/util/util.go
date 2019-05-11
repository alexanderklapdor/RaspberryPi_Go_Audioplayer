package util

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/serverfunctions/connectionfunctions"
	id3 "github.com/mikkyang/id3-go"
)

// check if string is element of the array
func StringInArray(str string, list []string) bool {
	// check if string is in string-list
	for _, element := range list {
		if element == str {
			return true
		}
	}
	return false
}

// error check
func Check(err error, source string) {
	if err != nil {
		if source == "Server" {
			connectionfunctions.Close()
		}
		logger.Error(err.Error())
		panic(err)
	}
}

// Shuffel String-Array
func Shuffle(array []string) []string {
	// Create new random variable based on the current time
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// swap elemnts random for each string position
	for n := len(array); n > 0; n-- {
		randIndex := r.Intn(n)
		array[n-1], array[randIndex] = array[randIndex], array[n-1]
	}
	//Return array
	return array
}

// JoinPath function
func JoinPath(source, target string) string {
	if path.IsAbs(target) {
		return target
	} // end of if
	return path.Join(path.Dir(source), target)
} // end of JoinPath

// GetFilesInFolder function
func GetFilesInFolder(folder string, supportedExtensions []string, depth int) []string {
	fileList := make([]string, 0)
	// Check if depth is > 0
	if depth > 0 {
		// Read directory
		files, err := ioutil.ReadDir(folder)
		Check(err, "Server")
		// For each file
		for _, file := range files {
			filename := JoinPath(folder, file.Name())
			fi, err := os.Stat(filename)
			Check(err, "Server")
			// Check if dir or file
			switch mode := fi.Mode(); {

			// Directory
			case mode.IsDir():
				newFolder := filename + "/"
				// Go into folder
				newFiles := GetFilesInFolder(newFolder, supportedExtensions, depth-1)
				// append files to fileList
				for _, newFile := range newFiles {
					fileList = append(fileList, newFile)
				} // end of for

			//File
			case mode.IsRegular():
				var extension = filepath.Ext(filename)
				// append files to fileList
				if StringInArray(extension, supportedExtensions) {
					fileList = append(fileList, filename)
				} // end of if
			} // end of switch
		} // end of for
	} else {
		logger.Info("Max depth reached")
	}
	return fileList
} // end of getFilesInFolder

//Get Minute and Secons from Seconds
func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	str := fmt.Sprintf("%dmin %dsec", minutes, seconds)
	return str
}

// Get Go Path
func GetGoExPath() string {
	gopath := os.Getenv("GOROOT")
	if gopath == "" {
		gopath = build.Default.GOROOT
	}
	return (gopath + "/bin/go")
}

// PrintMp3Infos function
func PrintMp3Infos(filePath string) string {
	//Check if Path exists
	if _, err := os.Stat(filePath); err == nil {
		//open file for id3 tags
		mp3File, err := id3.Open(filePath)
		Check(err, "Server")
		//close file at the end
		defer mp3File.Close()
		//get Tag Infos
		title := mp3File.Title()
		artist := mp3File.Artist()
		album := mp3File.Album()
		//get Audio length
		blength, lengtherr := exec.Command("mp3info", "-p", "%S", filePath).Output()
		Check(lengtherr, "Server")
		//check if one information is empty
		if title == "" || artist == "" || album == "" || string(blength[:]) == "" {
			return filePath
		} else {
			//print Infos
			length, err := strconv.Atoi(string(blength[:]))
			Check(err, "Server")
			return ("Title: " + title + "\t\tArtist: " + artist + "\t\tAlbum: " + album + "\t\tLength: " + SecondsToMinutes(length))
		}
	}
	return "Path not exists"
}
