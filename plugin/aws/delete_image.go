package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func converting(r []*types.ObjectIdentifier) []types.ObjectIdentifier {
	convertedArr := []types.ObjectIdentifier{}
	for _, item := range r {
		convertedArr = append(convertedArr, *item)
	}

	return convertedArr
}

func (s *s3Provider) DeleteImages(ctx context.Context, fileKeys []string) error {
	del := &types.Delete{
		Objects: converting(toOIDs(fileKeys)),
		Quiet:   aws.Bool(false),
	}

	doi := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Delete: del,
	}

	res, err := s.service.DeleteObjects(ctx, doi)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}

func toOIDs(keys []string) []*types.ObjectIdentifier {
	ret := make([]*types.ObjectIdentifier, len(keys))
	for i := 0; i < len(ret); i++ {
		oid := &types.ObjectIdentifier{
			Key: &(keys[i]),
		}
		ret[i] = oid
	}
	return ret
}

func (s *s3Provider) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.s3Bucket),
		Key:    aws.String(key),
	}

	res, err := s.service.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	s.logger.Infoln(res)

	return nil
}
