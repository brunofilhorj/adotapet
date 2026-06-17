package out

import "context"

type MediaStoragePort interface {
	CreateUploadURL(ctx context.Context, cmd CreateUploadURLCommand) (PresignedUpload, error)
	DeleteObject(ctx context.Context, objectKey string) error
}

type CreateUploadURLCommand struct {
	FileName    string
	ContentType string
	Purpose     string
}

type PresignedUpload struct {
	ObjectKey string
	UploadURL string
	PublicURL string
	ExpiresIn int
}
