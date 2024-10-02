package aws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *s3Provider) Upload(ctx context.Context, fileName string, cloudFolder string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	ext := filepath.Ext(file.Name())
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	fileKey := fmt.Sprintf("/%s/%s", cloudFolder, newFileName)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.s3Bucket),
		Key:           aws.String(fileKey),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err = s.service.PutObject(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com%s", s.cfg.s3Bucket, fileKey), nil
}

func (s *s3Provider) UploadFileData(ctx context.Context, fileData []byte, fileName string) (string, error) {
	fileBytes := bytes.NewReader(fileData)
	fileType := http.DetectContentType(fileData)

	fileKey := fmt.Sprintf("/%s", fileName)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(s.cfg.s3Bucket),
		Key:           aws.String(fileKey),
		Body:          fileBytes,
		ContentLength: aws.Int64(int64(len(fileData))),
		ContentType:   aws.String(fileType),
	}

	_, err := s.service.PutObject(ctx, params)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com%s", s.cfg.s3Bucket, fileKey), nil
}
