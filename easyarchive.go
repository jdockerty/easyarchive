package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	runningOS = runtime.GOOS
)

type configArchive struct {
	ArchiveLocation string `json:"Archive Location"`
}

var (
	currentConfig configArchive
)

func createConfig() {
	configFile, _ := os.Create("config.json")
	defer configFile.Close()
}

func writeToConfigFile() {

	fmt.Println("Writing to config.json...")

	jsonString, _ := json.MarshalIndent(currentConfig, "", "\t")
	ioutil.WriteFile("config.json", jsonString, os.ModePerm)
}

func setArchivePath(fp string) {
	cleanPath := filepath.Clean(fp)
	currentConfig.ArchiveLocation = cleanPath

	fmt.Println("Archive path set to", currentConfig.ArchiveLocation)
	writeToConfigFile()
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading user input string")
		panic(err)
	}

	return input
}

// func configFileExists(archivePath string) {
// Check file exists, if not then create it.
// }

func readConfigFile() configArchive {
	plan, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading file.")
		panic(err)
	}
	var data configArchive

	err = json.Unmarshal(plan, &data)
	if err != nil {
		return currentConfig
	}

	return data
}

func main() {

	conf := readConfigFile()

	if len(conf.ArchiveLocation) > 0 {
		fmt.Println("Path set")
	} else {
		// Set the archive path and re-run the main function after writing to config.json

		fmt.Println("The archive file path is not set.")

		if runningOS == "windows" {

			fmt.Printf("Enter the file path of the folder you wish to use, e.g. 'C:\\Users\\Jack\\Desktop\\MyBackups'")

			filepath := readUserInput()
			filepath = strings.Replace(filepath, "\r\n", "", -1)

			setArchivePath(filepath)
			main()

		} else if runningOS == "linux" {

			fmt.Println("Enter the file path of the folder you wish to use, e.g. /home/jack/mybackupfolder")

			filepath := readUserInput()
			filepath = strings.Replace(filepath, "\n", "", -1)

			setArchivePath(filepath)
			main()
		}
	}
}
