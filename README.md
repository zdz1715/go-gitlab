# go-gitlab
GitLab Go SDK

## Contents
- [Installation](#Installation)
- [Quick start](#quick-start)
- [ToDo](#todo)

## Installation
```shell
go get -u github.com/zdz1715/go-gitlab@latest
```


## Quick start
### OAuth授权码模式
```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/zdz1715/go-gitlab"

	"github.com/zdz1715/ghttp"
)

func main() {
	// OAuth授权码模式
	// docs: https://docs.gitlab.com/ee/api/oauth2.html#authorization-code-flow
	clientID := "YourClientID"
	clientSecret := "YourClientSecret"
	redirectURI := "http://127.0.0.1"
	credential := &gitlab.OAuthCredential{
		Endpoint:     gitlab.CloudEndpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}

	client, err := gitlab.NewClient(credential, &gitlab.Options{
		ClientOpts: []ghttp.ClientOption{
			ghttp.WithDebug(true),
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 先生成页面,获取code
	authURL := client.OAuth.GenerateAuthorizeURL(clientID, redirectURI, "api")
	fmt.Printf("click url: %s", authURL)

	// 监听输出，手动输入获取的code
	codeChan := make(chan string, 1)
	go func() {
		buf := bufio.NewScanner(os.Stdin)
		fmt.Print("\ninput code: ")
		for buf.Scan() {
			codeChan <- buf.Text()
		}
	}()

	select {
	case code := <-codeChan:
		_ = os.Stdin.Close()
		// 通过code先手动获取一次token，获取之后在token有效期内，请求别的接口会自动带上token
		fmt.Printf("auth by code: %+v\n", code)
		tk, err := client.OAuth.GetAccessToken(context.Background(), &gitlab.GetAccessTokenOptions{
			Code: code,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("auth token: %+v\n", tk)
		// 获取版本
		ver, err := client.Version.GetVersion(context.Background())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("version: %+v\n", ver)
		// 若是想刷新token
		tk, err = client.OAuth.GetAccessToken(context.Background(), &gitlab.GetAccessTokenOptions{
			RefreshToken: tk.RefreshToken,
		})
		fmt.Printf("RefreshToken: %+v\n", tk)
	}

}

```

### OAuth密码模式
```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zdz1715/go-gitlab"

	"github.com/zdz1715/ghttp"
)

func main() {
	// OAuth密码模式
	// docs: https://docs.gitlab.com/ee/api/oauth2.html#resource-owner-password-credentials-flow
	credential := &gitlab.PasswordCredential{
		Endpoint: gitlab.CloudEndpoint,
		Username: "YourUsername",
		Password: "YourPassword",
	}

	client, err := gitlab.NewClient(credential, &gitlab.Options{
		ClientOpts: []ghttp.ClientOption{
			ghttp.WithDebug(true),
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 无需手动获取token，执行下面方法会自动获取token，在有效期内不会重复请求获取token，当然你也可以手动获取token存起来
	// 获取版本
	ver, err := client.Version.GetVersion(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tk, _ := client.OAuth.GetAccessToken(context.Background())
	fmt.Printf("version: %+v\ntoken: %s\n", ver, tk.AccessToken)
}

```
### 直接设置token
```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zdz1715/ghttp"
	"github.com/zdz1715/go-gitlab"
)

func main() {
	// 直接设置token
	credential := &gitlab.TokenCredential{
		Endpoint:    gitlab.CloudEndpoint,
		AccessToken: "token",
	}

	client, err := gitlab.NewClient(credential, &gitlab.Options{
		ClientOpts: []ghttp.ClientOption{
			ghttp.WithDebug(true),
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 获取版本
	ver, err := client.Version.GetVersion(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("version: %+v\n", ver)

}

```

## ToDo
> [!NOTE]
> 现在提供的方法不多，会逐渐完善，也欢迎来贡献代码，只需要编写参数结构体、响应结构体就可以很快添加一个方法，参考下方代码。
```go
type Version struct {
    Version  string `json:"version"`
    Revision string `json:"revision"`
}

func (vs *VersionService) GetVersion(ctx context.Context) (*Version, error) {
    const apiEndpoint = "/api/v4/version"
    var v Version
    if err := vs.client.InvokeByCredential(ctx, http.MethodGet, apiEndpoint, nil, &v); err != nil {
        return nil, err
    }
    return &v, nil
}
```