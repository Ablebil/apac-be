package usecase

import (
	"apac/internal/app/user/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/env"
	"apac/internal/infra/helper"
	res "apac/internal/infra/response"
	"apac/internal/infra/supabase"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	GetProfile(userId uuid.UUID) (*dto.GetProfileResponse, *res.Err)
	EditProfile(userId uuid.UUID, payload *dto.EditProfileRequest) *res.Err
	AddPreference(userId uuid.UUID, payload *dto.AddPreferenceRequest) *res.Err
	RemovePreference(userId uuid.UUID, preferenceName string) *res.Err
}

type UserUsecase struct {
	userRepository repository.UserRepositoryItf
	supabase       supabase.SupabaseItf
	helper         helper.HelperItf
	env            *env.Env
}

func NewUserUsecase(env *env.Env, userRepository repository.UserRepositoryItf, supabase supabase.SupabaseItf, helper helper.HelperItf) UserUsecaseItf {
	return &UserUsecase{
		userRepository: userRepository,
		supabase:       supabase,
		helper:         helper,
		env:            env,
	}
}

func (uc *UserUsecase) GetProfile(userId uuid.UUID) (*dto.GetProfileResponse, *res.Err) {
	user, err := uc.userRepository.FindById(userId)
	if err != nil {
		return nil, res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return nil, res.ErrNotFound("User not found")
	}

	resp := user.ParseDTOGet()
	return &resp, nil
}

func (uc *UserUsecase) EditProfile(userId uuid.UUID, payload *dto.EditProfileRequest) *res.Err {
	user, err := uc.userRepository.FindById(userId)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if payload.Name != "" {
		user.Name = payload.Name
	}

	if payload.CurrentPassword != "" && payload.NewPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.CurrentPassword)); err != nil {
			return res.ErrForbidden("Incorrect old password")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return res.ErrInternalServer("Failed to hash password")
		}

		hashedPassword := string(hashed)
		user.Password = &hashedPassword
	}

	if payload.Photo != nil {
		if err := uc.helper.ValidateImage(payload.Photo); err != nil {
			return err
		}

		src, err := payload.Photo.Open()
		if err != nil {
			return res.ErrInternalServer("Error opening file")
		}

		defer src.Close()

		bucket := uc.env.SupabaseBucket

		if user.PhotoURL != uc.env.DefaultProfilePic {
			currentPhotoURL := user.PhotoURL
			index := strings.Index(currentPhotoURL, bucket)
			currentPhotoPath := currentPhotoURL[index+len(bucket+"/"):]

			if err := uc.supabase.DeleteFile(bucket, currentPhotoPath); err != nil {
				return res.ErrInternalServer("Failed to delete file")
			}
		}

		path := "profiles/" + user.ID.String() + "/" + payload.Photo.Filename
		contentType := payload.Photo.Header.Get("Content-Type")

		publicURL, err := uc.supabase.UploadFile(bucket, path, contentType, src)
		if err != nil {
			return res.ErrInternalServer("Failed to upload file")
		}

		user.PhotoURL = publicURL
	}

	if err := uc.userRepository.UpdateUser(userId, user); err != nil {
		return res.ErrInternalServer("Failed to update user")
	}

	return nil
}

func (uc *UserUsecase) AddPreference(userId uuid.UUID, payload *dto.AddPreferenceRequest) *res.Err {
	user, err := uc.userRepository.FindById(userId)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if payload.Preferences != nil {
		for _, pref := range payload.Preferences {
			if err := uc.userRepository.AddPreference(userId, pref); err != nil {
				return res.ErrInternalServer("Failed to add preference")
			}
		}
	}

	return nil
}

func (uc *UserUsecase) RemovePreference(userId uuid.UUID, preferenceName string) *res.Err {
	user, err := uc.userRepository.FindById(userId)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if err := uc.userRepository.RemovePreference(userId, preferenceName); err != nil {
		return res.ErrInternalServer("Failed to remove preference")
	}

	return nil
}
