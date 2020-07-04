# Easy Archive
An archive program to upload files into AWS S3 Glacier from a designated backups folder.

This will calculate the hashes of specific files within the folder and re-zip when changes occur, uploading a new zip archive.

It is assumed that archives are not run incredibly frequently, such as every 30 minutes, so archives are labelled with the current date in the format DD-MM-YYYY.

## Install

Firstly, ensure that your AWS credentials are set as this program uses the Go SDK. If running on an EC2 instance, make sure it has the appropriate IAM role with permissions to create a bucket and write to S3.
```
git clone https://github.com/jdockerty/easyarchive.git
cd easyarchive
go mod download
go build -v easyarchive.go
```

Once built, you can use the program as required for the specific OS. `./easyarchive` on Linux and `easyarchive.exe` on Windows. Placing the executable within your PATH variable will enable easy access from the command line. You can also execute the binary via cron for scheduled backups to be run at set intervals.
