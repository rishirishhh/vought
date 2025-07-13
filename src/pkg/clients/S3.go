package clients

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"io"
	_ "io"
	"strings"
	_ "strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
)

type IS3Client interface {
	ListObjects(ctx context.Context) ([]string, error)
	GetObjects(ctx context.Context, key string) (io.Reader, error)
	PutObjectInput(ctx context.Context, f io.Reader, path string) error
	CreateBucketIfDoesNotExists(ctx context.Context, bucketName string) error
	RemoveObject(ctx context.Context, path string) error
}

var _ IS3Client = s3Client{}

type s3Client struct {
	awsS3Client *s3.Client
	bucket      string
}

func NewS3Client(host, region, bucket, accessKey, pwdKey string) (IS3Client, error) {
	// Create custom credential provider
	creds := credentials.NewStaticCredentialsProvider(accessKey, pwdKey, "")

	// Load default config with region and creds
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}

	// Use custom endpoint
	if host != "" {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               host,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})
	}

	// Initialise S3 client
	client := &s3Client{
		awsS3Client: s3.NewFromConfig(cfg),
		bucket:      bucket,
	}

	if err := client.CreateBucketIfDoesNotExists(context.Background(), bucket); err != nil {
		return nil, err
	}

	return client, nil
}

func (s s3Client) ListObjects(ctx context.Context) ([]string, error) {
	delimiter := "/"
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		MaxKeys:   aws.Int32(10),
		Delimiter: &delimiter,
	}

	res, err := s.awsS3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	objectsName := make([]string, len(res.CommonPrefixes))

	for index, obj := range res.CommonPrefixes {
		if obj.Prefix == nil {
			continue
		}

		objectsName[index] = strings.TrimSuffix(*obj.Prefix, "/")
	}

	return objectsName, nil
}

func (s s3Client) GetObjects(ctx context.Context, key string) (io.Reader, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	response, err := s.awsS3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (s s3Client) PutObjectInput(ctx context.Context, fileReader io.Reader, path string) error {
	uploader := manager.NewUploader(s.awsS3Client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   fileReader,
	})

	return err
}

func (s s3Client) CreateBucketIfDoesNotExists(ctx context.Context, bucketName string) error {
	bucketListOutput, err := s.awsS3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	if bucketListOutput == nil {
		return errors.New("unable to get bucket list")
	}

	for _, b := range bucketListOutput.Buckets {
		if b.Name == nil {
			continue
		}

		if strings.Contains(*b.Name, bucketName) {
			// bucket already exists, nothing to do
			return nil
		}
	}

	log.Debugf("Bucket named '%s' not found, Creating it.", bucketName)

	inputCreate := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	_, err = s.awsS3Client.CreateBucket(ctx, inputCreate)
	return err
}

func (s s3Client) RemoveObject(ctx context.Context, path string) error {
	results, err := s.awsS3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		return err
	}

	for _, content := range results.Contents {
		// If we can go deeper in the filesystem, do it to remove all objects in
		// directory

		if *content.Key != path {
			err := s.RemoveObject(ctx, *content.Key)
			if err != nil {
				return err
			}
		}

		_, err = s.awsS3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    content.Key,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
