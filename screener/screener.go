package screener

//Imports
import (
	"fmt"
	"github.com/alexanderklapdor/RaspberryPi_Go_Audioplayer/logger"
	"strconv"
)

// Start Screen
func StartScreen() {

	fmt.Println("############################################")
	fmt.Println("#               Welcome to                 #")
	fmt.Println("#              Music Player                #")
	fmt.Println("############################################")
	fmt.Println("")

}

// End Screen
func EndScreen() {

	fmt.Println("")
	fmt.Println("********************************************")
	fmt.Println("*               Shutdown                   *")
	fmt.Println("********************************************")
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
