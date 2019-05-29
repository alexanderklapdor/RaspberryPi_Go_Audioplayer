package screener

// Imports
import (
	"fmt"
	"strconv"

	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
)

// Start Screen
func StartScreen() {
	fmt.Println("#######################################")
	fmt.Println("#      Start Music Player Client      #")
	fmt.Println("#######################################")
	fmt.Println("")

}

// End Screen
func EndScreen() {
	fmt.Println("")
	fmt.Println("#######################################")
	fmt.Println("#       End Music Player Client       #")
	fmt.Println("#######################################")
}

// Print Files
func PrintFiles(fileList []string, printFiles bool) {
	logger.Info("Found " + strconv.Itoa(len(fileList)) + " supported files")
	if printFiles {
		logger.Info("*********Filelist (Supported Files)*********")
		for _, fileElement := range fileList {
			logger.Info(fileElement)
		}
	}

}
