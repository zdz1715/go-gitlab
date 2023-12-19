package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GetAccessTokenOptions struct {
	Code         string
	RefreshToken string
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	CreatedAt    int64  `json:"created_at"`
}

// OAuthService
// GitLab API Docs: https://docs.gitlab.com/ee/api/oauth2.html
type OAuthService struct {
	client     *Client
	credential Credential
	store
}

type store struct {
	val    *AccessToken
	expire time.Time
}

func (s *store) value() *AccessToken {
	return s.val
}

func (s *store) memory(at *AccessToken) {
	s.val = at
	if at != nil {
		// 提前5分钟过期, 避免网络带来的延时
		s.expire = time.Now().Add(time.Duration(at.ExpiresIn)*time.Second - 5*time.Minute)
	}
}

func (s *store) IsExpired() bool {
	if s.val == nil {
		return true
	}
	if s.val.AccessToken == "" {
		return true
	}
	return time.Now().After(s.expire)
}

func (oa *OAuthService) GenerateAuthorizeURL(clientId, redirectUri, scope string) string {
	return fmt.Sprintf("%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s",
		oa.client.cc.Endpoint(),
		clientId,
		url.QueryEscape(redirectUri),
		url.QueryEscape(scope),
	)
}

func (oa *OAuthService) GetAccessToken(ctx context.Context, opts ...*GetAccessTokenOptions) (*AccessToken, error) {
	opt := new(GetAccessTokenOptions)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}
	// 如果是刷新token，则强制刷新
	if opt.RefreshToken == "" && !oa.store.IsExpired() {
		return oa.store.value(), nil
	}

	if oa.credential == nil {
		return nil, ErrCredential
	}

	if err := oa.credential.Valid(); err != nil {
		return nil, err
	}

	body := oa.credential.Body(opt)

	if body == nil {
		return nil, nil
	}

	const apiEndpoint = "/oauth/token"
	var respBody AccessToken
	if err := oa.client.Invoke(ctx, http.MethodPost, apiEndpoint, body, &respBody); err != nil {
		return nil, err
	}

	oa.store.memory(&respBody)

	return &respBody, nil
}
