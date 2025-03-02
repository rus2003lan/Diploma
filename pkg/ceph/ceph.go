package ceph

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type File struct {
	ContentType   string
	Data          io.ReadCloser
	ContentLength int
}

type Config struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
}

var ErrCephRequiredField = errors.New("no required fields in config")

type Client struct {
	s3 *s3.S3
}

func New(cfg *Config) (*Client, error) {
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, ErrCephRequiredField
	}

	cred := credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	ccfg := aws.NewConfig().
		WithRegion("us-east-1").
		WithCredentials(cred).
		WithEndpoint(cfg.Endpoint).
		WithS3ForcePathStyle(true)

	sess, err := session.NewSession(ccfg)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	svc := s3.New(sess)

	cl := &Client{s3: svc}

	err = cl.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *Client) PutFile(ctx context.Context, bucket, filename, contentType string, data []byte) error {
	return c.PutRFile(
		ctx, bucket, filename, contentType, bytes.NewReader(data),
	)
}

// PutRFile sends data to ceph.
func (c *Client) PutRFile(
	ctx context.Context,
	bucket, filename, contentType string,
	data io.ReadSeeker,
) error {
	//nolint:exhaustruct
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		ContentType: aws.String(contentType),
		Body:        data,
	}

	err := input.Validate()
	if err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	_, err = c.s3.PutObjectWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("put object %s: %w", filename, err)
	}

	return nil
}

func (c *Client) FetchFile(ctx context.Context, bucket, filename string) (*File, error) {
	//nolint:exhaustruct
	obj, err := c.s3.GetObjectWithContext(
		ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get object with filename %s: %w", filename, err)
	}

	res := File{
		ContentType:   zeroValueIfNull(obj.ContentType),
		Data:          obj.Body,
		ContentLength: int(zeroValueIfNull(obj.ContentLength)),
	}

	return &res, nil
}

func (c *Client) DeleteFile(ctx context.Context, bucket, filename string) error {
	_, err := c.s3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{ //nolint:exhaustruct
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	return fmt.Errorf("delete file: %w", err)
}

func (c *Client) CreateBucket(ctx context.Context, bucketName string) error {
	if !c.bucketExists(ctx, bucketName) {
		_, err := c.s3.CreateBucketWithContext(ctx, &s3.CreateBucketInput{ //nolint:exhaustruct
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{ //nolint:exhaustruct
				LocationConstraint: aws.String("us-east-1"),
			},
		})
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
	}

	return nil
}

func (c *Client) bucketExists(ctx context.Context, bucketName string) bool {
	_, err := c.s3.HeadBucketWithContext(ctx, &s3.HeadBucketInput{ //nolint:exhaustruct
		Bucket: aws.String(bucketName),
	})

	return err == nil
}

func zeroValueIfNull[T any](p *T) T {
	if p == nil {
		var t T
		return t
	}

	return *p
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.s3.ListBucketsWithContext(ctx, nil)
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}
