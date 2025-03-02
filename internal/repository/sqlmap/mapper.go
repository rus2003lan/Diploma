package sqlmap

import (
	"errors"
	"fmt"

	"diploma-project/internal/model"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func mapCephErrToModelErr(err error) error {
	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		code := awsErr.Code()
		if code == s3.ErrCodeNoSuchBucket ||
			code == s3.ErrCodeNoSuchKey ||
			code == s3.ErrCodeNoSuchUpload ||
			code == "NotFound" ||
			code == "InvalidRange" {
			return fmt.Errorf("%s: %w", code, model.ErrNotFound)
		}

		return fmt.Errorf("internal ceph: %w", err)
	}

	return fmt.Errorf("error is not from aws: %w", err)
}
