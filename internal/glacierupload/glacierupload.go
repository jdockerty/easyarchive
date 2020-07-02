package glacierupload

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"

	"fmt"
	"os"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func s3Session() (*session.Session, *s3.S3) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)
	if err != nil {
		exitErrorf("Error creating session\n%v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	return sess, svc
}

// CreateBucket will create an S3 bucket that is ready for use, combining a prefix with a generated UUID for a DNS-compliant name.
// This returns the bucketName to write into the config.json file.
func CreateBucket() string {

	fmt.Println("An S3 bucket is not present, creating a new one for you in the form of easyarchive-[random_ID_suffix].")

	randomIDSuffix := uuid.New().String()
	bucketName := fmt.Sprintf("easyarchive-%s", randomIDSuffix)

	_, svc := s3Session()

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucketName, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucketName)
	}

	fmt.Printf("Bucket %q successfully created\n", bucketName)

	return bucketName
}

// UploadArchive will upload the provided .zip into S3 Glacier.
func UploadArchive(bucket, zipFile string) {

	sess, _ := s3Session()

	file, err := os.Open(zipFile)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(zipFile),
		Body:        file,
		ContentType: aws.String("application/zip"),
		// StorageClass: aws.String("GLACIER"),
	})

	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", zipFile, bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", zipFile, bucket)
}
