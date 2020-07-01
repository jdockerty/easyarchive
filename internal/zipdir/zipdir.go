package zipdir

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func pathSepartor() string {
	return string(filepath.Separator)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func zipUp(filename string, files []string, archivePath string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)

	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file, archivePath); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, archivePath string) error {
	fileString := fmt.Sprintf("%s%s%s", archivePath, pathSepartor(), filename)
	fileToZip, err := os.Open(fileString)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func timeStamp() string {
	now := time.Now()
	currentDate := fmt.Sprintf("%02d-%02d-%d", now.Day(), now.Month(), now.Year())
	return currentDate
}

// ZipFiles will compress and zip the files contained within the set archive path from config.json.
func ZipFiles(files []string, archivePath string) {
	output := fmt.Sprintf("%s.zip", timeStamp())

	if err := zipUp(output, files, archivePath); err != nil {
		panic(err)
	}
	fmt.Println("Zipped File successfully:", output)
}
