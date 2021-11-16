package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
)

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(s *session.Session, file graphql.Upload) (string, error) {
  // get the file size and read
  // the file content into a buffer
  size := file.Size
  buffer := make([]byte, size)
  _, err:= file.File.Read(buffer)

  if err != nil {
	  return "", err 
  }

  // create a unique file name for the file
  tempFileName := "shorts/" + bson.NewObjectId().Hex() + filepath.Ext(file.Filename)
	
  // config settings: this is where you choose the bucket,
  // filename, content-type and storage class of the file
  // you're uploading
  _, err = s3.New(s).PutObject(&s3.PutObjectInput{
     Bucket:               aws.String("audioshortbucket"),
     Key:                  aws.String(tempFileName),
     ACL:                  aws.String("public-read"),// could be private if you want it to be access by only authorized users
     Body:                 bytes.NewReader(buffer),
     ContentLength:        aws.Int64(int64(size)),
     ContentType:        aws.String(http.DetectContentType(buffer)),
     ContentDisposition:   aws.String("attachment"),
     ServerSideEncryption: aws.String("AES256"),
     StorageClass:         aws.String("INTELLIGENT_TIERING"),
  })

  if err != nil {
     return "", err
  }

  return "https://audioshortbucket.s3.amazonaws.com/"+tempFileName, nil
}

func HandleFileUpload(file graphql.Upload) (string, error) {
	// allow only 1MB of file size
	maxSize := int64(1024000) 
	
	
	if file.Size > maxSize {
	   fmt.Println("Image too large. Max Size: " + string(maxSize) + ", Uploaded File Size: " + string(file.Size))
	   return "", errors.New("Image too large. Max Size: " + string(maxSize) + ", Uploaded File Size: " + string(file.Size))
	}

	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			string(os.Getenv("AWS_ID")), // id
			string(os.Getenv("AWS_KEY")),   // secret
			""),  // token can be left blank for now
	})

	if err != nil {
		fmt.Println("Create Session Error: ",err)
		return "", errors.New("Create Session Error")
	}

	// upload file to s3
	fileName, err := UploadFileToS3(s, file)
	
	if err != nil {
		fmt.Println("Upload File Error: ",err)
		return "", errors.New("Upload File Error")
	}

	fmt.Println("Image uploaded successfully: ", fileName)

	return fileName, nil
}