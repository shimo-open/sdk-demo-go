package awos

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gotomicro/ego/core/econf"
)

// AwosService provides object storage service using AWS S3 compatible API
type AwosService struct {
	svc *s3.S3
}

var bucket string

// Init initializes a new AwosService instance with configured S3 client
func Init() *AwosService {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			econf.GetString("awos.accessKeyId"),
			econf.GetString("awos.secretAccessKey"),
			""),
		Region:           aws.String("minio"),
		Endpoint:         aws.String(econf.GetString("awos.endpoint")),
		S3ForcePathStyle: aws.Bool(econf.GetBool("awos.s3ForcePathStyle")),
		DisableSSL:       aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	svc := s3.New(sess)
	bucket = econf.GetString("awos.bucket")

	return &AwosService{svc: svc}
}

// Save uploads data to the object storage with the specified key
func (a *AwosService) Save(key string, data []byte) error {
	_, err := a.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}

// GetUploadURL creates an upload URL valid for the given number of seconds
func (a *AwosService) GetUploadURL(key string, expireSeconds int64) (string, error) {
	req, _ := a.svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return req.Presign(time.Duration(expireSeconds) * time.Second)
}

// Get retrieves data from the object storage by key
func (a *AwosService) Get(key string) ([]byte, error) {
	result, err := a.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(result.Body)
	return buf.Bytes(), nil
}

// GetDownloadURL creates a download URL valid for the given number of seconds
func (a *AwosService) GetDownloadURL(key string, filename string, expireSeconds int64) (string, error) {
	req, _ := a.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf(`attachment; filename="%s"`, filename)),
	})

	return req.Presign(time.Duration(expireSeconds) * time.Second)
}

// Remove deletes an object from the storage by key
func (a *AwosService) Remove(key string) error {
	_, err := a.svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
