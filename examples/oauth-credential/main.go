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
		// default endpoint: https://gitlab.com
		//Endpoint: gitlab.CloudEndpoint,
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
