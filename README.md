# Easy Archive
An archive program to upload files into AWS S3 Glacier from a designated backups folder.

This will calculate the hashes of specific files within the folder and re-zip when changes occur, uploading a new zip archive.

It is assumed that archives are not run incredibly frequently, such as every 30 minutes, so archives are labelled with the current date in the format DD-MM-YY.

## Install

Firstly, ensure that your AWS credentials are set as this program uses the Go SDK. If running on an EC2 instance, make sure it has the appropriate IAM role with permissions to create a bucket and write to S3.
```
go get github.com/jdockerty/easyarchive
```
OR
```
git clone https://github.com/jdockerty/easyarchive.git
cd EasyArchive
go mod download
go build -v easyarchive.go
```

Once built, you can use the program as required for the specific OS. `./easyarchive` on Linux and `easyarchive.exe` on Windows. Placing the executable within your PATH variable will enable easy access from the command line.