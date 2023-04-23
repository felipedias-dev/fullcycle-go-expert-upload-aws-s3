package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/felipedias-dev/fullcycle-go-expert-upload-aws-s3/configs"
)

var (
	s3Client *s3.S3
	s3Bucket string
)

func init() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(config.AwsRegion),
			Credentials: credentials.NewStaticCredentials(
				config.AwsKey,
				config.AwsSecret,
				"",
			),
		},
	)
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(sess)
	s3Bucket = config.S3Bucket
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}
		uploadFile(files[0].Name())
	}
}

func uploadFile(filename string) {
	completeFilename := fmt.Sprintf("./tmp/%s", filename)
	fmt.Printf("Uploading file %s to bucker %s\n", completeFilename, s3Bucket)
	f, err := os.Open(completeFilename)
	if err != nil {
		fmt.Printf("Error opening file %s\n", completeFilename)
		return
	}
	defer f.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		fmt.Printf("Error uploading file %s\n", completeFilename)
		return
	}
	fmt.Printf("File %s uploaded successfully\n", completeFilename)
}
