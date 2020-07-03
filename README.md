# Easy Archive
A simple archive program to upload files into AWS S3 Glacier from a designated backups folder.

## Install

Firstly, ensure that your AWS credentials are set. If you're running on an EC2 instance, then ensure it has the appropriate IAM role with permissions to create a bucket and write to S3.

```
git clone https://github.com/jdockerty/EasyArchive.git
cd EasyArchive
go mod download
go build -v easyarchive.go
```