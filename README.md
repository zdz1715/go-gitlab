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
- [OAuth授权码模式](./examples/oauth-credential/main.go)
- [OAuth密码模式](./examples/password-credential/main.go)
- [Token](./examples/token-credential/main.go)
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