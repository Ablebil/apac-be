package usecase

import (
	"apac/internal/app/auth/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/entity"
	"apac/internal/domain/env"
	"apac/internal/infra/email"
	"apac/internal/infra/jwt"
	res "apac/internal/infra/response"

	"golang.org/x/crypto/bcrypt"

	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(*dto.RegisterRequest) *res.Err
	VerifyOTP(*dto.VerifyOTPRequest) (string, string, *res.Err)
	ChoosePreference(*dto.ChoosePreference) *res.Err
	Login(*dto.LoginRequest) (string, string, *res.Err)
	RefreshToken(*dto.RefreshToken) (string, string, *res.Err)
	Logout(*dto.LogoutRequest) *res.Err
}

type AuthUsecase struct {
	repo  repository.AuthRepositoryItf
	jwt   jwt.JWTItf
	db    *gorm.DB
	email email.EmailItf
	env   *env.Env
}

func NewAuthUsecase(env *env.Env, db *gorm.DB, authRepository repository.AuthRepositoryItf, jwt jwt.JWTItf, email email.EmailItf) AuthUsecaseItf {
	return &AuthUsecase{
		repo:  authRepository,
		jwt:   jwt,
		db:    db,
		email: email,
		env:   env,
	}
}

func (uc *AuthUsecase) Register(payload *dto.RegisterRequest) *res.Err {
	user, err := uc.repo.FindByEmail(payload.Email)

	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return res.ErrInternalServer("Failed to hash password")
	}
	hashedPassword := string(hashed)

	otp := fmt.Sprintf("%06d", 100000+rand.New(rand.NewSource(time.Now().UnixNano())).Intn(900000))
	otpExpiresAt := time.Now().Add(5 * time.Minute)

	if user != nil {
		if user.GoogleID != nil && user.Password == nil {
			if err := uc.repo.Update(payload.Email, &entity.User{
				Password:     &hashedPassword,
				OTP:          &otp,
				OTPExpiresAt: &otpExpiresAt,
			}); err != nil {
				return res.ErrInternalServer("Failed to update user")
			}

			if err := uc.email.SendOTPEmail(payload.Email, otp); err != nil {
				return res.ErrInternalServer("Failed to send OTP email")
			}

			return nil
		}

		return res.ErrBadRequest("Email already registered")
	}

	user = &entity.User{
		Name:         payload.Name,
		Email:        payload.Email,
		Password:     &hashedPassword,
		OTP:          &otp,
		OTPExpiresAt: &otpExpiresAt,
	}

	if err := uc.repo.Create(user); err != nil {
		return res.ErrInternalServer("Failed to create user")
	}

	if err := uc.email.SendOTPEmail(payload.Email, otp); err != nil {
		return res.ErrInternalServer("Failed to send OTP email")
	}

	return nil
}

func (uc *AuthUsecase) VerifyOTP(payload *dto.VerifyOTPRequest) (string, string, *res.Err) {
	user, err := uc.repo.FindByEmail(payload.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return "", "", res.ErrNotFound("User not found")
	}

	if user.OTP == nil || *user.OTP != payload.OTP || time.Now().After(*user.OTPExpiresAt) {
		return "", "", res.ErrBadRequest("Invalid or expired OTP")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	if err := uc.repo.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	if err := uc.repo.Update(payload.Email, &entity.User{
		OTP:          nil,
		OTPExpiresAt: nil,
		Verified:     true,
	}); err != nil {
		return "", "", res.ErrInternalServer("Failed to update user")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) Login(payload *dto.LoginRequest) (string, string, *res.Err) {
	user, err := uc.repo.FindByEmail(payload.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return "", "", res.ErrNotFound("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.Password)); err != nil {
		return "", "", res.ErrBadRequest("Incorrect email or password")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, payload.RememberMe)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	refreshTokens, err := uc.repo.GetUserRefreshTokens(user.ID)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to get refresh tokens")
	}

	if len(refreshTokens) >= 2 {
		uc.repo.RemoveRefreshToken(refreshTokens[0].Token)
	}

	if err := uc.repo.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshToken(payload *dto.RefreshToken) (string, string, *res.Err) {
	user, err := uc.repo.FindByRefreshToken(payload.RefreshToken)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return "", "", res.ErrNotFound("User not found")
	}

	if _, err := uc.jwt.VerifyRefreshToken(payload.RefreshToken); err != nil {
		return "", "", res.ErrUnauthorized("Invalid refresh token")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	if err := uc.repo.RemoveRefreshToken(payload.RefreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to remove refresh token")
	}

	if err := uc.repo.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) Logout(payload *dto.LogoutRequest) *res.Err {
	user, err := uc.repo.FindByRefreshToken(payload.RefreshToken)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if err := uc.repo.RemoveRefreshToken(payload.RefreshToken); err != nil {
		return res.ErrInternalServer("Failed to remove refresh token")
	}

	return nil
}

func (uc *AuthUsecase) ChoosePreference(payload *dto.ChoosePreference) *res.Err {
	user, err := uc.repo.FindByEmail(payload.Email)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if payload.Preferences != nil {
		for _, pref := range payload.Preferences {
			if err := uc.repo.AddPreference(user.ID, pref); err != nil {
				return res.ErrInternalServer("Failed to add preference")
			}
		}
	}

	return nil
}
