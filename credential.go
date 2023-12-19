package gitlab

import (
	"errors"
	"net/http"

	"github.com/zdz1715/ghttp"
)

var (
	ErrCredential = errors.New("invalid credential")
)

type Credential interface {
	Body(opts *GetAccessTokenOptions) any
	GetEndpoint() string
	GenerateCallOptions(token *AccessToken) (*ghttp.CallOptions, error)
	Valid() error
}

// TokenCredential
// Docs: https://docs.gitlab.com/ee/api/rest/#authentication
type TokenCredential struct {
	Endpoint    string    `json:"endpoint" xml:"endpoint"`
	TokenType   TokenType `json:"type" xml:"type"`
	AccessToken string    `json:"token" xml:"token"`
}

func (t *TokenCredential) Body(opts *GetAccessTokenOptions) any {
	return nil
}

func (t *TokenCredential) GetEndpoint() string {
	return t.Endpoint
}

func (t *TokenCredential) GenerateCallOptions(token *AccessToken) (*ghttp.CallOptions, error) {
	callOpts := &ghttp.CallOptions{}
	if t.AccessToken == "" {
		return callOpts, nil
	}
	switch t.TokenType {
	case TokenTypeTypeJob:
		callOpts.BeforeHook = func(request *http.Request) error {
			request.Header.Set("JOB-TOKEN", t.AccessToken)
			return nil
		}
	case TokenTypePrivate:
		callOpts.BeforeHook = func(request *http.Request) error {
			request.Header.Set("PRIVATE-TOKEN", t.AccessToken)
			return nil
		}
	default:
		callOpts.BearerToken = t.AccessToken
	}
	return callOpts, nil
}

func (t *TokenCredential) Valid() error {
	if t.AccessToken == "" {
		return ErrCredential
	}
	return nil
}

// PasswordCredential
// note: The Resource Owner Password Credentials is disabled for users with two-factor authentication turned on.
// These users can access the API using personal access tokens instead.
// docs: https://docs.gitlab.com/ee/api/oauth2.html#resource-owner-password-credentials-flow
type PasswordCredential struct {
	Endpoint string `json:"endpoint" xml:"endpoint"`
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
}

func (p *PasswordCredential) Body(opts *GetAccessTokenOptions) any {
	return map[string]string{
		"grant_type": "password",
		"username":   p.Username,
		"password":   p.Password,
	}
}

func (p *PasswordCredential) GetEndpoint() string {
	return p.Endpoint
}

func (p *PasswordCredential) GenerateCallOptions(token *AccessToken) (*ghttp.CallOptions, error) {
	tk := ""
	if token != nil {
		tk = token.AccessToken
	}
	return &ghttp.CallOptions{
		BearerToken: tk,
	}, nil
}

func (p *PasswordCredential) Valid() error {
	if p.Username == "" || p.Password == "" {
		return ErrCredential
	}
	return nil
}

// OAuthCredential
// docs: https://docs.gitlab.com/ee/api/oauth2.html#authorization-code-flow
type OAuthCredential struct {
	Endpoint     string `json:"endpoint" xml:"endpoint"`
	ClientID     string `json:"client_id" xml:"client_id"`
	ClientSecret string `json:"client_secret" xml:"client_secret"`
	RedirectURI  string `json:"redirect_uri" xml:"redirect_uri"`
}

func (c *OAuthCredential) GetEndpoint() string {
	return c.Endpoint
}

func (c *OAuthCredential) Body(opts *GetAccessTokenOptions) any {
	body := map[string]string{
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"redirect_uri":  c.RedirectURI,
	}
	if opts.Code != "" {
		body["grant_type"] = "authorization_code"
		body["code"] = opts.Code
	}

	if opts.RefreshToken != "" {
		body["grant_type"] = "refresh_token"
		body["refresh_token"] = opts.RefreshToken
	}
	return body
}

func (c *OAuthCredential) GenerateCallOptions(token *AccessToken) (*ghttp.CallOptions, error) {
	tk := ""
	if token != nil {
		tk = token.AccessToken
	}
	return &ghttp.CallOptions{
		BearerToken: tk,
	}, nil
}

func (c *OAuthCredential) Valid() error {
	if c.ClientID == "" || c.ClientSecret == "" || c.RedirectURI == "" {
		return ErrCredential
	}
	return nil
}
