package gitlab

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zdz1715/ghttp"
)

type service struct {
	client *Client
}

type Options struct {
	ClientOpts []ghttp.ClientOption
}

type Client struct {
	cc   *ghttp.Client
	opts *Options

	common service

	OAuth *OAuthService

	Branches      *BranchesService
	Commits       *CommitsService
	MergeRequests *MergeRequestsService
	Tags          *TagsService
	Users         *UsersService
	Projects      *ProjectsService
	Version       *VersionService
	Metadata      *MetadataService
}

func NewClient(credential Credential, opts *Options) (*Client, error) {
	if opts == nil {
		opts = &Options{}
	}

	clientOptions := []ghttp.ClientOption{
		ghttp.WithEndpoint(CloudEndpoint),
	}

	if len(opts.ClientOpts) > 0 {
		clientOptions = append(clientOptions, opts.ClientOpts...)
	}

	// 覆盖错误
	clientOptions = append(clientOptions,
		ghttp.WithNot2xxError(func() ghttp.Not2xxError {
			return new(Error)
		}),
	)

	cc := ghttp.NewClient(clientOptions...)

	c := &Client{
		cc:   cc,
		opts: opts,
	}

	c.common.client = c

	c.OAuth = &OAuthService{client: c.common.client}
	c.Branches = (*BranchesService)(&c.common)
	c.Commits = (*CommitsService)(&c.common)
	c.MergeRequests = (*MergeRequestsService)(&c.common)
	c.Tags = (*TagsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Projects = (*ProjectsService)(&c.common)
	c.Metadata = (*MetadataService)(&c.common)
	c.Version = (*VersionService)(&c.common)

	if credential != nil {
		if err := c.SetCredential(credential); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) SetCredential(credential Credential) error {
	if credential == nil {
		return ErrCredential
	}

	if err := credential.Valid(); err != nil {
		return err
	}

	c.cc.SetEndpoint(credential.GetEndpoint())

	if c.OAuth != nil {
		c.OAuth.credential = credential
	}

	return nil
}

func (c *Client) InvokeByCredential(ctx context.Context, method, path string, args any, reply any) error {
	accessToken, err := c.OAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	callOpts, err := c.OAuth.credential.GenerateCallOptions(accessToken)
	if err != nil {
		return err
	}

	return c.Invoke(ctx, method, path, args, reply, callOpts)
}

func (c *Client) Invoke(ctx context.Context, method, path string, args any, reply any, opts ...*ghttp.CallOptions) error {
	callOpts := new(ghttp.CallOptions)

	if len(opts) > 0 && opts[0] != nil {
		callOpts = opts[0]
	}

	if method == http.MethodGet && args != nil {
		callOpts.Query = args
		args = nil
	}

	_, err := c.cc.Invoke(ctx, method, path, args, reply, callOpts)
	return err
}

// Error data-validation-and-error-reporting + OAuth error
// GitLab API docs: https://docs.gitlab.com/ee/api/rest/#data-validation-and-error-reporting
type Error struct {
	Message any `json:"message"`

	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *Error) String() string {
	if e.ErrorDescription != "" {
		return e.ErrorDescription
	}
	if e.Error != "" {
		return e.Error
	}
	if e.Message != nil {
		switch msg := e.Message.(type) {
		case string:
			return msg
		default:
			b, _ := json.Marshal(e.Message)
			return string(b)
		}
	}
	return ""
}

func (e *Error) Reset() {
	e.Message = nil
	e.Error = ""
	e.ErrorDescription = ""
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
// GitLab API docs: https://docs.gitlab.com/ee/api/rest/index.html#pagination
type ListOptions struct {
	// GitLab API docs: https://docs.gitlab.com/ee/api/rest/index.html#offset-based-pagination
	// For paginated result sets, page of results to retrieve.
	Page int `query:"page,omitempty" json:"page,omitempty"`
	// For paginated result sets, the number of results to include per page.
	PerPage int `query:"per_page,omitempty" json:"per_page,omitempty"` // default: 20 max: 100

	// GitLab API docs: https://docs.gitlab.com/ee/api/rest/index.html#keyset-based-pagination
	Pagination string `query:"pagination,omitempty" json:"pagination,omitempty"`
	OrderBy    string `query:"order_by,omitempty" json:"order_by,omitempty"`
	Sort       Sort   `query:"sort,omitempty" json:"sort,omitempty"`
}

func NewListOptions(page int, perPage ...int) *ListOptions {
	if page <= 0 {
		page = 1
	}
	l := &ListOptions{
		Page:    page,
		PerPage: 20,
	}
	if len(perPage) > 0 && perPage[0] > 0 {
		l.PerPage = perPage[0]
	}
	return l
}

func NewKetSetListOptions(orderBy string, sort Sort, perPage ...int) *ListOptions {
	l := &ListOptions{
		Pagination: "keyset",
		OrderBy:    orderBy,
		Sort:       sort,
	}
	if len(perPage) > 0 && perPage[0] > 0 {
		l.PerPage = perPage[0]
	}
	return l
}
