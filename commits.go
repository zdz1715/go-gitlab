package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CommitsService
// GitLab API docs: https://docs.gitlab.com/ee/api/commits.html
type CommitsService service

type Commit struct {
	ID             string            `json:"id"`
	ShortID        string            `json:"short_id"`
	Title          string            `json:"title"`
	AuthorName     string            `json:"author_name"`
	AuthorEmail    string            `json:"author_email"`
	AuthoredDate   time.Time         `json:"authored_date"`
	CommitterName  string            `json:"committer_name"`
	CommitterEmail string            `json:"committer_email"`
	CommittedDate  time.Time         `json:"committed_date"`
	CreatedAt      time.Time         `json:"created_at"`
	Message        string            `json:"message"`
	ParentIDs      []string          `json:"parent_ids"`
	Stats          CommitStats       `json:"stats"`
	Status         BuildStateValue   `json:"status"`
	LastPipeline   PipelineInfo      `json:"last_pipeline"`
	ProjectID      int               `json:"project_id"`
	Trailers       map[string]string `json:"trailers"`
	WebURL         string            `json:"web_url"`
}

type CommitStats struct {
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	Total     int `json:"total"`
}

// ListCommitsOptions represents the available ListCommits() options.
//
// GitLab API docs: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits
type ListCommitsOptions struct {
	*ListOptions
	RefName     *string    `json:"ref_name,omitempty" query:"ref_name"`
	Since       *time.Time `json:"since,omitempty" query:"since"`
	Until       *time.Time `json:"until,omitempty" query:"until"`
	Path        *string    `json:"path,omitempty" query:"path"`
	Author      *string    `json:"author,omitempty" query:"author"`
	All         *bool      `json:"all,omitempty" query:"all"`
	WithStats   *bool      `json:"with_stats,omitempty" query:"with_stats"`
	FirstParent *bool      `json:"first_parent,omitempty" query:"first_parent"`
	Trailers    *bool      `json:"trailers,omitempty" query:"trailers"`
}

// ListCommits gets a list of repository commits in a project.
//
// GitLab API docs: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits
func (s *CommitsService) ListCommits(ctx context.Context, projectId string, opts *ListCommitsOptions) ([]*Commit, error) {
	apiEndpoint := fmt.Sprintf("/api/v4/projects/%s/repository/commits", projectId)
	var v []*Commit
	if err := s.client.InvokeByCredential(ctx, http.MethodGet, apiEndpoint, opts, &v); err != nil {
		return nil, err
	}
	return v, nil
}
