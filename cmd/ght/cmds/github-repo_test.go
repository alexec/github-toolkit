package cmds

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewGithubRepo(t *testing.T) {
	expected := GithubRepo{accessToken: os.Getenv("ACCESS_TOKEN"), host: "github.com", owner: "alexec", repo: "github-toolkit"}
	t.Run("SSH", func(t *testing.T) {
		r := NewGithubRepo("git@github.com:alexec/github-toolkit.git")
		assert.Equal(t, r, expected)
		assert.Equal(t, r.BaseURL().String(), "https://api.github.com/")

	})
	t.Run("HTTPS", func(t *testing.T) {
		r := NewGithubRepo("https://github.com/alexec/github-toolkit.git")
		assert.Equal(t, r, expected)
		assert.Equal(t, r.BaseURL().String(), "https://api.github.com/")
	})
	t.Run("Enterprise", func(t *testing.T) {
		r := NewGithubRepo("git@github.my-company.com:my-argo/my-proj")
		assert.Equal(t,  r.BaseURL().String(), "https://github.my-company.com/api/v3/")
	})
}
