package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Client interface {
	PutPrivateObjectFromFile(ctx context.Context, localFile, destName string) error
	PutObjectFromFile(ctx context.Context, localFile, destName string) error

	PubPrivateObject(ctx context.Context, reader io.Reader, destName string) error
	PubObject(ctx context.Context, reader io.Reader, destName string) error

	GetObjectURL(destName string) string
	DownloadFileWithURL(signedURL, localpath string) error
}

type AliyunOSS struct {
	endpoint  string
	bucket    string
	rawClient *oss.Client
}

func NewOSSClient(endpoint, accessKeyId, accessKeySecret, bucket string) (Client, error) {
	return newAliyunOSS(endpoint, accessKeyId, accessKeySecret, bucket)
}

func newAliyunOSS(endpoint, accessKeyId, accessKeySecret, bucket string) (*AliyunOSS, error) {
	if bucket == "" {
		return nil, errors.New("Invalid bucket")
	}
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return &AliyunOSS{
		rawClient: client,
		bucket:    bucket,
		endpoint:  endpoint,
	}, nil
}

// OssProgressListener is the progress listener
type OssProgressListener struct {
}

// ProgressChanged handles progress event
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	// case oss.TransferDataEvent:

	// fmt.Printf("Transfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.\n",
	// 	event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("Transfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("Transfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func (p *AliyunOSS) PubPrivateObject(ctx context.Context, reader io.Reader, destName string) error {
	bucket, err := p.rawClient.Bucket(p.bucket)
	if err != nil {
		return err
	}

	options := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}

	return bucket.PutObject(destName, reader, options...)
}

func (p *AliyunOSS) PubObject(ctx context.Context, reader io.Reader, destName string) error {
	bucket, err := p.rawClient.Bucket(p.bucket)
	if err != nil {
		return err
	}
	return bucket.PutObject(destName, reader)
}

func (p *AliyunOSS) PutPrivateObjectFromFile(ctx context.Context, localFile, destName string) error {

	bucket, err := p.rawClient.Bucket(p.bucket)
	if err != nil {
		return err
	}

	options := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}

	return bucket.PutObjectFromFile(destName, localFile, options...)
}

func (p *AliyunOSS) PutObjectFromFile(ctx context.Context, localFile, destName string) error {
	bucket, err := p.rawClient.Bucket(p.bucket)
	if err != nil {
		return err
	}
	return bucket.PutObjectFromFile(destName, localFile, oss.Progress(&OssProgressListener{}))
}

func (p *AliyunOSS) DownloadFileWithURL(signedURL, localpath string) error {
	bucket, err := p.rawClient.Bucket(p.bucket)
	if err != nil {
		return err
	}
	return bucket.GetObjectToFileWithURL(signedURL, localpath)
}

func (p *AliyunOSS) GetObjectURL(destName string) string {
	domain := p.endpoint
	if strings.HasPrefix(p.endpoint, "http") {
		domain = strings.Replace(domain, "http://", "", -1)
		domain = strings.Replace(domain, "https://", "", -1)
	}
	return fmt.Sprintf("https://%s.%s/%s", p.bucket, domain, destName)
}
