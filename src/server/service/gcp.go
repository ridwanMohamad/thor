package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"thor/src/server/config"

	"cloud.google.com/go/storage"
)

type gcpClient struct {
	ProjectID  string
	BucketName string
}

const (
	gcsBaseDomain = "storage.googleapis.com"
)

var (
	gcpClientOnce sync.Once
	GcpClient     *gcpClient
)

func InitializeGcp(c config.Service) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", c.Gcp.CredentialPath)
	gcpClientOnce.Do(func() {
		GcpClient = &gcpClient{
			ProjectID:  c.Gcp.ProjectID,
			BucketName: c.Gcp.BucketName,
		}
	})
}

func (g *gcpClient) UploadBase64(ctx context.Context, filename string, file io.Reader) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(GcpClient.BucketName)
	w := bucket.Object(filename).NewWriter(ctx)
	w.CacheControl = "no-store"
	defer w.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return err
}

func (g *gcpClient) Upload(ctx context.Context, filename string, file multipart.File) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(GcpClient.BucketName)
	w := bucket.Object(filename).NewWriter(ctx)
	defer w.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}

// URL Generate HTTPs url to fetch google clould file
func (g *gcpClient) URL(filename string) string {
	return fmt.Sprintf("https://%s/%s/%s", gcsBaseDomain, g.BucketName, filename)
}
