package cmds

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewGithubRepo(t *testing.T) {
	expected := GithubRepo{accessToken: os.Getenv("ACCESS_TOKEN"), host: "github.com", owner: "alexec", repo: "github-toolkit"}
	assert.Equal(t, NewGithubRepo("git@github.com:alexec/github-toolkit.git"), expected)
	assert.Equal(t, NewGithubRepo("https://github.com/alexec/github-toolkit.git"), expected)
}
