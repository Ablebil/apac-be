package supabase

import (
	"apac/internal/domain/env"
	"apac/internal/infra/response"
	"fmt"
	"io"
	"mime/multipart"

	supabase "github.com/nedpals/supabase-go"
)

type SupabaseItf interface {
	UploadFile(bucket string, filePath string, contentType string, file multipart.File) (string, error)
	UploadFileFromIOReader(bucket string, filePath string, contentType string, file io.Reader) (string, error)
	DeleteFile(bucket string, filePath string) error
}

type Supabase struct {
	Client supabase.Client
}

func NewSupabase(env *env.Env) SupabaseItf {
	storageClient := supabase.CreateClient(fmt.Sprintf("%s", env.SupabaseUrl), env.SupabaseAnonKey)

	return &Supabase{
		Client: *storageClient,
	}
}

func (s *Supabase) UploadFile(bucket, filePath, contentType string, file multipart.File) (string, error) {
	wrapper := safeWrapper(func() error {
		s.Client.Storage.From(bucket).Upload(filePath, file, &supabase.FileUploadOptions{ContentType: contentType})
		return nil
	})

	if wrapper != nil {
		return "", wrapper
	}

	publicUrl := s.Client.Storage.From(bucket).GetPublicUrl(filePath).SignedUrl
	return publicUrl, nil
}

func (s *Supabase) UploadFileFromIOReader(bucket, filePath, contentType string, file io.Reader) (string, error) {
	wrapper := safeWrapper((func() error {
		s.Client.Storage.From(bucket).Upload(filePath, file, &supabase.FileUploadOptions{ContentType: contentType})
		return nil
	}))

	if wrapper != nil {
		return "", wrapper
	}

	publicUrl := s.Client.Storage.From(bucket).GetPublicUrl(filePath).SignedUrl
	return publicUrl, nil
}

func (s *Supabase) DeleteFile(bucket, filePath string) error {
	wrapper := safeWrapper(func() error {
		s.Client.Storage.From(bucket).Remove([]string{filePath})
		return nil
	})

	return wrapper
}

func safeWrapper(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = response.ErrInternalServer(fmt.Sprintf("Unexpected panic: %v", r))
		}
	}()
	return f()
}
