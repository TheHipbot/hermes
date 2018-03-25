package fs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	testConfigPath string
	testTargetFile string
)

func init() {
	testConfigPath = "/test/.hermes/"
	testTargetFile = ".hermes_test"
	viper.Set("config_path", testConfigPath)
	viper.Set("target_name", testTargetFile)
}

func TestSetupCreateDir(t *testing.T) {
	configFs = memfs.New()
	Setup()

	if _, err := configFs.Lstat(testConfigPath); err != nil {
		t.Errorf("Setup should create config dir\n%s", err)
	}
}

func TestSetTarget(t *testing.T) {
	configFs = memfs.New()
	target := "/repo_dir/github.com/TheHipbot/hermes/"
	bs := make([]byte, 40)

	// set up to create config_dir in memfs
	Setup()

	SetTarget(target)
	file, err := configFs.Open(fmt.Sprintf("%s%s", testConfigPath, testTargetFile))
	if err != nil {
		t.Errorf("SetTarget should create a target file\n%s", err)
	}
	file.Read(bs)
	if !strings.Contains(string(bs), target) {
		t.Errorf("Target file content incorrect\n%s", strings.TrimSpace(string(bs)))
	}
}
