package oauth

import (
	"apac/internal/domain/dto"
	"apac/internal/domain/env"
	"context"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type OAuthItf interface {
	GenerateAuthLink(state string) (string, error)
	ExchangeToken(code string) (*oauth2.Token, error)
	GetProfile(token *oauth2.Token) (dto.GoogleProfileResponse, error)
}

type OAuth struct {
	googleOAuthConfig *oauth2.Config
}

func NewOAuth(env *env.Env) OAuthItf {
	return &OAuth{
		googleOAuthConfig: &oauth2.Config{
			ClientID:     env.GoogleClientID,
			ClientSecret: env.GoogleClientSecret,
			RedirectURL:  env.GoogleRedirectUrl,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (o *OAuth) GenerateAuthLink(state string) (string, error) {
	return o.googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (o *OAuth) ExchangeToken(code string) (*oauth2.Token, error) {
	token, err := o.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return &oauth2.Token{}, err
	}

	return token, nil
}

func (o *OAuth) GetProfile(token *oauth2.Token) (dto.GoogleProfileResponse, error) {
	userData := dto.GoogleProfileResponse{}
	client := o.googleOAuthConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return userData, err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return userData, err
	}

	userData.Username = strings.Split(userInfo.Email, "@")[0]
	userData.Email = userInfo.Email
	userData.Name = userInfo.Name
	userData.IsVerified = *userInfo.VerifiedEmail

	return userData, nil
}
