package aws

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/haohmaru3000/go_sdk/logger"
	"github.com/haohmaru3000/go_sdk/sdkcm"
)

var (
	ErrS3ApiKeyMissing       = sdkcm.CustomError("ErrS3ApiKeyMissing", "AWS S3 API key is missing")
	ErrS3ApiSecretKeyMissing = sdkcm.CustomError("ErrS3ApiSecretKeyMissing", "AWS S3 API secret key is missing")
	ErrS3RegionMissing       = sdkcm.CustomError("ErrS3RegionMissing", "AWS S3 region is missing")
	ErrS3BucketMissing       = sdkcm.CustomError("ErrS3ApiKeyMissing", "AWS S3 bucket is missing")
)

type s3Provider struct {
	name   string
	prefix string
	logger logger.Logger

	cfg s3Config

	config         *aws.Config
	service        *s3.Client
	presignService *s3.PresignClient
}

type s3Config struct {
	s3ApiKey    string
	s3ApiSecret string
	s3Region    string
	s3Bucket    string
}

func NewS3Provider(prefix ...string) *s3Provider {
	pre := "aws-s3"

	if len(prefix) > 0 {
		pre = prefix[0]
	}

	return &s3Provider{
		name:   "aws-s3",
		prefix: pre,
	}
}

func (s *s3Provider) Get() interface{} {
	return s
}

func (s *s3Provider) Name() string {
	return s.name
}

func (s *s3Provider) InitFlags() {
	flag.StringVar(&s.cfg.s3ApiKey, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-key"), "", "S3 API key")
	flag.StringVar(&s.cfg.s3ApiSecret, fmt.Sprintf("%s-%s", s.GetPrefix(), "api-secret"), "", "S3 API secret key")
	flag.StringVar(&s.cfg.s3Region, fmt.Sprintf("%s-%s", s.GetPrefix(), "region"), "", "S3 region")
	flag.StringVar(&s.cfg.s3Bucket, fmt.Sprintf("%s-%s", s.GetPrefix(), "bucket"), "", "S3 bucket")
}

func (s *s3Provider) Configure() error {
	ctx := context.TODO()

	s.logger = logger.GetCurrent().GetLogger(s.Name())

	if err := s.cfg.check(); err != nil {
		s.logger.Errorln(err)
		return err
	}

	credential := credentials.NewStaticCredentialsProvider(s.cfg.s3ApiKey, s.cfg.s3ApiSecret, "")
	_, err := credential.Retrieve(ctx)
	if err != nil {
		s.logger.Errorln(err)
		return err
	}

	config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.cfg.s3Region),
		config.WithCredentialsProvider(credential),
	)
	if err != nil {
		s.logger.Errorln(err)
		return err
	}

	service := s3.NewFromConfig(config)
	presignService := s3.NewPresignClient(s3.NewFromConfig(config))

	s.config = &config
	s.service = service
	s.presignService = presignService

	return nil
}

func (s *s3Provider) GetPrefix() string {
	return s.prefix
}

func (s *s3Provider) Run() error {
	return s.Configure()
}

func (s *s3Provider) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (cfg *s3Config) check() error {
	if len(cfg.s3ApiKey) < 1 {
		return ErrS3ApiKeyMissing
	}
	if len(cfg.s3ApiSecret) < 1 {
		return ErrS3ApiSecretKeyMissing
	}
	if len(cfg.s3Bucket) < 1 {
		return ErrS3BucketMissing
	}
	if len(cfg.s3Region) < 1 {
		return ErrS3RegionMissing
	}
	return nil
}
