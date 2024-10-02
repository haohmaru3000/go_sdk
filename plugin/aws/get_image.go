package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *s3Provider) GetImageWithExpireLink(ctx context.Context, imageKey string, duration time.Duration) (string, error) {
	req, err := s.presignService.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Key:    aws.String(imageKey),
	}, func(o *s3.PresignOptions) {
		o.Expires = duration
	})

	return req.URL, err
}
