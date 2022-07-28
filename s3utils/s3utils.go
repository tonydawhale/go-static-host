package s3utils

import (
	"log"
	"os"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var S3Client *s3.S3
var Uploader *s3manager.Uploader

func checkCredentials() {
	_, ok1 := os.LookupEnv("AWS_ACCESS_KEY_ID")
	if !ok1 {
		log.Fatal("AWS_ACCESS_KEY_ID not set")
	}
	_, ok2 := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !ok2 {
		log.Fatal("AWS_SECRET_ACCESS_KEY not set")
	}
	_, ok3 := os.LookupEnv("AWS_ENDPOINT")
	if !ok3 {
		log.Fatal("AWS_ENDPOINT not set")
	}
	_, ok4 := os.LookupEnv("BUCKET_NAME")
	if !ok4 {
		log.Fatal("BUCKET_NAME not set")
	}
}

func Init() {
	checkCredentials()
	log.Println("S3 Credentials Validated")
}

func GetS3Object(key string) (out *s3.GetObjectOutput, err error) {
	client, err := createClient()
	if err != nil {
		return nil, err
	}

	return client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key: aws.String(key),
	})
}

func UploadS3Object(name string, body io.Reader) (*s3manager.UploadOutput, error){
	client, err := createClient()
	if err != nil {
		return nil, err
	}
	Uploader = s3manager.NewUploaderWithClient(client)
	return Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:	aws.String(name),
		Body:	body,
		ACL:   aws.String("public-read"),
	})
}

func createClient() (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"), 
			os.Getenv("AWS_SECRET_ACCESS_KEY"), 
			"",
		),
		Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
		Region: aws.String("default"),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}