package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/jdockerty/easyarchive/internal/glacierupload"
	"github.com/jdockerty/easyarchive/internal/md5calc"
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
	BucketName      string    `json:"S3 Bucket"`
	Hashes          []hashVal `json:"Hashes"`
}

var (
	currentConfig configArchive
)

func createConfig() {

	configFile, err := os.Create(configFile)
	if err != nil {
		log.Println("Error creating config.json.")
		panic(err)
	}
	defer configFile.Close()
}

func setArchivePathAndBucket(fp, bucket string) {
	cleanPath := filepath.Clean(fp)

	currentConfig.ArchiveLocation = cleanPath
	currentConfig.BucketName = bucket
	currentConfig.Hashes = nil

	writeArchivePathAndBucketToConfig(cleanPath, bucket)
	log.Println("Archive path set to", cleanPath)
	log.Println("You may now place files into the set path and they will be archived into S3 Glacier upon running the program again.")
	time.Sleep(2)
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading user input string")
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

func calcHashes() []hashVal {
	m, err := md5calc.MD5All(currentConfig.ArchiveLocation)
	if err != nil {
		log.Println("Error running MD5All on archive path.")
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

func writeHashes(new []hashVal) {
	currentConfig.Hashes = new

	dataStream, _ := json.MarshalIndent(currentConfig, "", "\t")
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

	if len(oldHashes) != len(newHashes) {

		return true

	} else if isEqualHash(oldHashes, newHashes) == false {

		// Hashes have changed if they are not equal
		return true

	}

	// All hashes are the same, no action needed.
	return false
}

func writeArchivePathAndBucketToConfig(path, bucket string) {
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

func archiveBucketExist() bool {
	if len(currentConfig.BucketName) > 0 {
		return true
	}

	return false
}

func initialSetup(eol string) {
	filepath := readUserInput()
	filepath = strings.Replace(filepath, eol, "", -1)

	bucket := glacierupload.CreateBucket()

	setArchivePathAndBucket(filepath, bucket)

	main()
}


func main() {
	currentConfig = readConfigFile()

	// If archive path is set and an S3 bucket has been created in config.json, proceed.
	if archivepath := currentConfig.ArchiveLocation; len(archivepath) > 0 && archiveBucketExist() {

		newH := calcHashes()

		if hashesChanged(currentConfig.Hashes, newH) {

			log.Println("Change detected, writing new hashes to config.json...")
			writeHashes(newH)

			log.Println("Archiving files...")
			filenames := getFilenames(currentConfig.Hashes)
			outputZip := zipdir.ZipFiles(filenames, archivepath)
			glacierupload.UploadArchive(currentConfig.BucketName, outputZip)

		} else {
			log.Println("No changes detected.")
			}

	} else {
		// Set the archive path and re-run the main function after writing to config.json

		log.Println("The archive file path is not set.")

		if runningOS == "windows" {

			log.Println("Enter the file path of the folder you wish to use, e.g. 'C:\\Users\\Jack\\Desktop\\MyBackups'")

			initialSetup("\r\n")


		} else if runningOS == "linux" {

			log.Println("Enter the file path of the folder you wish to use, e.g. /home/jack/mybackupfolder")

			initialSetup("\n")
			
		}
	}
}
