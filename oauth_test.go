package gitlab

import "os"

var testPasswordCredential = &PasswordCredential{
	Username: os.Getenv("TEST_GITLAB_USERNAME"),
	Password: os.Getenv("TEST_GITLAB_PASSWORD"),
}
