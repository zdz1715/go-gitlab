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
