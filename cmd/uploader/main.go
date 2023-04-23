package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/felipedias-dev/fullcycle-go-expert-upload-aws-s3/configs"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
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

	uploadControl := make(chan struct{}, 100)
	errorUploadControl := make(chan string, 10)

	go func() {
		for {
			select {
			case fileName := <-errorUploadControl:
				uploadControl <- struct{}{}
				wg.Add(1)
				go uploadFile(fileName, uploadControl, errorUploadControl)
			}
		}
	}()

	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}
		wg.Add(1)
		uploadControl <- struct{}{}
		go uploadFile(files[0].Name(), uploadControl, errorUploadControl)
	}
	wg.Wait()
}

func uploadFile(filename string, uploadControl <-chan struct{}, errorUploadControl chan<- string) {
	defer wg.Done()
	completeFilename := fmt.Sprintf("./tmp/%s", filename)
	fmt.Printf("Uploading file %s to bucket %s\n", completeFilename, s3Bucket)
	f, err := os.Open(completeFilename)
	if err != nil {
		fmt.Printf("Error opening file %s\n", completeFilename)
		errorUploadControl <- filename
		<-uploadControl
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
		errorUploadControl <- filename
		<-uploadControl
		return
	}
	fmt.Printf("File %s uploaded successfully\n", completeFilename)
	<-uploadControl
}
