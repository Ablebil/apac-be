package usecase

import (
	authRepository "apac/internal/app/auth/repository"
	userRepository "apac/internal/app/user/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/entity"
	"apac/internal/domain/env"
	"apac/internal/infra/email"
	"apac/internal/infra/jwt"
	"apac/internal/infra/oauth"
	"apac/internal/infra/redis"
	res "apac/internal/infra/response"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	Register(payload *dto.RegisterRequest) *res.Err
	VerifyOTP(payload *dto.VerifyOTPRequest) (string, string, *res.Err)
	Login(payload *dto.LoginRequest) (string, string, *res.Err)
	RefreshToken(payload *dto.RefreshToken) (string, string, *res.Err)
	Logout(payload *dto.LogoutRequest) *res.Err
	GoogleLogin() (string, *res.Err)
	GoogleCallback(payload *dto.GoogleCallbackRequest) (string, string, *res.Err)
	ChoosePreference(payload *dto.ChoosePreferenceResponse) *res.Err
}

type AuthUsecase struct {
	authRepository authRepository.AuthRepositoryItf
	userRepository userRepository.UserRepositoryItf
	jwt            jwt.JWTItf
	db             *gorm.DB
	redis          redis.RedisItf
	email          email.EmailItf
	env            *env.Env
	oauth          oauth.OAuthItf
}

func NewAuthUsecase(
	env *env.Env,
	db *gorm.DB,
	redis redis.RedisItf,
	authRepository authRepository.AuthRepositoryItf,
	userRepository userRepository.UserRepositoryItf,
	jwt jwt.JWTItf,
	email email.EmailItf,
	oauth oauth.OAuthItf,
) AuthUsecaseItf {
	return &AuthUsecase{
		authRepository: authRepository,
		userRepository: userRepository,
		jwt:            jwt,
		redis:          redis,
		db:             db,
		email:          email,
		env:            env,
		oauth:          oauth,
	}
}

func (uc *AuthUsecase) Register(payload *dto.RegisterRequest) *res.Err {
	user, err := uc.authRepository.FindByEmail(payload.Email)

	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return res.ErrInternalServer("Failed to hash password")
	}
	hashedPassword := string(hashed)

	otp := fmt.Sprintf("%06d", 100000+rand.New(rand.NewSource(time.Now().UnixNano())).Intn(900000))
	expiration := 5 * time.Minute

	if user != nil {
		if user.GoogleID != nil && user.Password == nil {
			if err := uc.authRepository.Update(payload.Email, &entity.User{
				Password: &hashedPassword,
			}); err != nil {
				return res.ErrInternalServer("Failed to update user")
			}

			_ = uc.redis.SetOTP(payload.Email, otp, expiration)

			if err := uc.email.SendOTPEmail(payload.Email, otp); err != nil {
				return res.ErrInternalServer("Failed to send OTP email")
			}

			return nil
		}

		return res.ErrConflict("Email already registered")
	} else {
		user = &entity.User{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: &hashedPassword,
		}

		if err := uc.authRepository.Create(user); err != nil {
			return res.ErrInternalServer("Failed to create user")
		}

		_ = uc.redis.SetOTP(payload.Email, otp, expiration)

		if err := uc.email.SendOTPEmail(payload.Email, otp); err != nil {
			return res.ErrInternalServer("Failed to send OTP email")
		}
	}

	return nil
}

func (uc *AuthUsecase) VerifyOTP(payload *dto.VerifyOTPRequest) (string, string, *res.Err) {
	user, err := uc.authRepository.FindByEmail(payload.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return "", "", res.ErrNotFound("User not found")
	}

	storedOtp, err := uc.redis.GetOTP(payload.Email)
	if err != nil || storedOtp != payload.OTP {
		return "", "", res.ErrBadRequest("Invalid or expired OTP")
	}

	_ = uc.redis.DeleteOTP(payload.Email)

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	if err := uc.authRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	user.Verified = true

	if err := uc.authRepository.Update(user.Email, user); err != nil {
		return "", "", res.ErrInternalServer("Failed to update user")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) Login(payload *dto.LoginRequest) (string, string, *res.Err) {
	user, err := uc.authRepository.FindByEmail(payload.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil || bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.Password)) != nil {
		return "", "", res.ErrUnauthorized("Incorrect email or password")
	}

	if !user.Verified {
		return "", "", res.ErrForbidden("Account not verified")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, payload.RememberMe)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	refreshTokens, err := uc.authRepository.GetUserRefreshTokens(user.ID)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to get refresh tokens")
	}

	if len(refreshTokens) >= 2 {
		if err := uc.authRepository.RemoveRefreshToken(refreshTokens[0].Token); err != nil {
			return "", "", res.ErrInternalServer("Failed to remove refresh token")
		}
	}

	if err := uc.authRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshToken(payload *dto.RefreshToken) (string, string, *res.Err) {
	user, err := uc.authRepository.FindByRefreshToken(payload.RefreshToken)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return "", "", res.ErrForbidden("Invalid refresh token")
	}

	if _, err := uc.jwt.VerifyRefreshToken(payload.RefreshToken); err != nil {
		return "", "", res.ErrForbidden("Expired refresh token")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	if err := uc.authRepository.RemoveRefreshToken(payload.RefreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to remove refresh token")
	}

	if err := uc.authRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) Logout(payload *dto.LogoutRequest) *res.Err {
	user, err := uc.authRepository.FindByRefreshToken(payload.RefreshToken)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrForbidden("Invalid refresh token")
	}

	if err := uc.authRepository.RemoveRefreshToken(payload.RefreshToken); err != nil {
		return res.ErrInternalServer("Failed to remove refresh token")
	}

	return nil
}

func (uc *AuthUsecase) GoogleLogin() (string, *res.Err) {
	stateLength := uc.env.StateLength
	bytes := make([]byte, stateLength)
	if _, err := crand.Read(bytes); err != nil {
		return "", res.ErrInternalServer("Failed to generate state")
	}

	state := base64.RawURLEncoding.EncodeToString(bytes)

	if len(state) > stateLength {
		state = state[:stateLength]
	}

	if err := uc.redis.Set(state, []byte(state), uc.env.StateExpiry); err != nil {
		return "", res.ErrInternalServer("Failed to save oauth state")
	}

	url, err := uc.oauth.GenerateLink(state)
	if err != nil {
		return "", res.ErrInternalServer("Failed to generate oauth link")
	}

	return url, nil
}

func (uc *AuthUsecase) GoogleCallback(payload *dto.GoogleCallbackRequest) (string, string, *res.Err) {
	if payload.Error != "" {
		return "", "", res.ErrInternalServer("Google callback returns with error: " + payload.Error)
	}

	state, err := uc.redis.Get(payload.State)
	if err != nil {
		return "", "", res.ErrUnauthorized("OAuth state not found")
	}

	if string(state) != payload.State {
		return "", "", res.ErrUnauthorized("OAuth state not found")
	}

	if err := uc.redis.Delete(payload.State); err != nil {
		return "", "", res.ErrInternalServer()
	}

	token, err := uc.oauth.ExchangeToken(payload.Code)
	if err != nil {
		return "", "", res.ErrInternalServer(err.Error())
	}

	profile, err := uc.oauth.GetProfile(token)
	if err != nil {
		return "", "", res.ErrInternalServer(err.Error())
	}

	user, err := uc.authRepository.FindByEmail(profile.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to find user")
	}

	if user != nil {
		if user.GoogleID == nil {
			if err := uc.authRepository.Update(user.Email, &entity.User{
				GoogleID: &profile.ID,
			}); err != nil {
				return "", "", res.ErrInternalServer("Failed to update user")
			}
		}
	} else {
		user = &entity.User{
			Email:    profile.Email,
			Name:     profile.Name,
			GoogleID: &profile.ID,
			Verified: profile.Verified,
		}

		if err := uc.authRepository.Create(user); err != nil {
			return "", "", res.ErrInternalServer("Failed to create user")
		}
	}

	if !user.Verified {
		return "", "", res.ErrForbidden("Account not verified")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken(user.ID, false)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate refresh token")
	}

	refreshTokens, err := uc.authRepository.GetUserRefreshTokens(user.ID)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to get refresh tokens")
	}

	if len(refreshTokens) >= 2 {
		if err := uc.authRepository.RemoveRefreshToken(refreshTokens[0].Token); err != nil {
			return "", "", res.ErrInternalServer("Failed to remove refresh token")
		}
	}

	if err := uc.authRepository.AddRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", res.ErrInternalServer("Failed to add refresh token")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID, user.Name, user.Email)
	if err != nil {
		return "", "", res.ErrInternalServer("Failed to generate access token")
	}

	return accessToken, refreshToken, nil
}

func (uc *AuthUsecase) ChoosePreference(payload *dto.ChoosePreferenceResponse) *res.Err {
	user, err := uc.authRepository.FindByEmail(payload.Email)
	if err != nil {
		return res.ErrInternalServer("Failed to find user")
	}

	if user == nil {
		return res.ErrNotFound("User not found")
	}

	if payload.Preferences != nil {
		for _, pref := range payload.Preferences {
			if err := uc.userRepository.AddPreference(user.ID, pref); err != nil {
				return res.ErrInternalServer("Failed to add preference")
			}
		}
	}

	return nil
}
