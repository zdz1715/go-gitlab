package gitlab

import (
	"context"
	"testing"

	"github.com/zdz1715/ghttp"
	"github.com/zdz1715/go-utils/goutils"
)

func TestProjectsService_ListProjects(t *testing.T) {
	client, err := NewClient(testPasswordCredential, &Options{
		ClientOpts: []ghttp.ClientOption{
			ghttp.WithDebug(true),
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	// 	查询全部部门
	reply, err := client.Projects.ListProjects(context.Background(), &ListProjectsOptions{
		ListOptions: NewListOptions(1, 10),
		OrderBy:     goutils.Ptr("last_activity_at"),
		Membership:  goutils.Ptr(true),
		Search:      goutils.Ptr(""),
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", reply)
	}
}
