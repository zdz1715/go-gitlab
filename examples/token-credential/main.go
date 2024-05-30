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
		// default endpoint: https://gitlab.com
		//Endpoint: gitlab.CloudEndpoint,
		AccessToken: "token",
	}

	client, err := gitlab.NewClient(credential, &gitlab.Options{
		ClientOpts: []ghttp.ClientOption{
			ghttp.WithDebug(ghttp.DefaultDebug),
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
