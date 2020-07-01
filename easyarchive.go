package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/jdockerty/easyarchive/internal/md5calc"
	// "github.com/jdockerty/easyarchive/internal/glacierupload"
	"github.com/jdockerty/easyarchive/internal/zipdir"
)

const (
	runningOS  = runtime.GOOS
	configFile = "config.json"
)

type hashVal struct {
	Filename  string `json:"Filename"`
	HashValue string `json:"Value"`
}

type configArchive struct {
	ArchiveLocation string    `json:"Archive Location"`
	Hashes          []hashVal `json:"Hash Values"`
}

var (
	currentConfig configArchive
)

func createConfig() {
	configFile, _ := os.Create(configFile)
	defer configFile.Close()
}

func setArchivePath(fp string) {
	cleanPath := filepath.Clean(fp)
	currentConfig.ArchiveLocation = cleanPath

	writeArchivePathToConfig(cleanPath)
	fmt.Println("Archive path set to", currentConfig.ArchiveLocation)
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

func readConfigFile() configArchive {
	configJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		createConfig()
	}
	var data configArchive

	err = json.Unmarshal(configJSON, &data)
	if err != nil {
		return data
	}

	return data
}

func calcHashes(conf configArchive) []hashVal {
	m, err := md5calc.MD5All(conf.ArchiveLocation)
	if err != nil {
		fmt.Println("Error running MD5All on archive path.")
		panic(err)
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	var tempNewHashes []hashVal
	for _, path := range paths {

		// fmt.Printf("%x  %s\n", m[path], path)
		basePath := filepath.Base(path)
		hashString := fmt.Sprintf("%x", m[path])
		hashHolder := hashVal{Filename: basePath, HashValue: hashString}
		tempNewHashes = append(tempNewHashes, hashHolder)
	}
	return tempNewHashes
}

func writeHashes(currentConf configArchive, new []hashVal) {
	currentConf.Hashes = new

	dataStream, _ := json.MarshalIndent(currentConf, "", "\t")
	ioutil.WriteFile(configFile, dataStream, 0644)
}

func isEqualHash(old, new []hashVal) bool {
	for i, v := range new {
		if v.HashValue == old[i].HashValue {
			continue
		} else {
			return false
		}
	}
	return true
}

func hashesChanged(oldHashes, newHashes []hashVal) bool {
	// fmt.Println("old:", oldHashes)
	// fmt.Println("new:", newHashes)

	if len(newHashes) > len(oldHashes) {

		return true

	} else if len(newHashes) < len(oldHashes) {

		return true

	} else if isEqualHash(oldHashes, newHashes) == false {

		return true

	}

	return false
}

func writeArchivePathToConfig(path string) {
	currentConfig.ArchiveLocation = path
	currentConfig.Hashes = nil

	output, _ := json.MarshalIndent(currentConfig, "", "\t")
	ioutil.WriteFile(configFile, output, 0644)
}

func getFilenames(val []hashVal) []string {
	var files []string

	for _, v := range val {
		files = append(files, v.Filename)
	}
	return files
}

func main() {

	currentConfig = readConfigFile()

	fmt.Println("start conf", currentConfig)

	if archivepath := currentConfig.ArchiveLocation; len(archivepath) > 0 {
		newH := calcHashes(currentConfig)
		if hashesChanged(currentConfig.Hashes, newH) {
			fmt.Println("Change detected, updating config.json...")
			writeHashes(currentConfig, newH)

			fmt.Println("Archiving files...")
			filenames := getFilenames(currentConfig.Hashes)
			zipdir.ZipFiles(filenames, currentConfig.ArchiveLocation)

		} else {
			fmt.Println("No action required.")
		}
	} else {
		// Set the archive path and re-run the main function after writing to config.json

		fmt.Println("The archive file path is not set.")

		if runningOS == "windows" {

			fmt.Printf("Enter the file path of the folder you wish to use, e.g. 'C:\\Users\\Jack\\Desktop\\MyBackups'")

			filepath := readUserInput()
			filepath = strings.Replace(filepath, "\r\n", "", -1)

			setArchivePath(filepath)
			// writeToConfigFile() CAN POSSIBILITY REMOVE THIS FUNCTION
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
