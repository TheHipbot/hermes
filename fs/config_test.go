package fs

import (
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/src-d/go-billy.v4/memfs"
)

func TestSetupCreateDir(t *testing.T) {
	testConfigPath := "/test/.hermes/"
	configFs = memfs.New()

	viper.Set("config_path", testConfigPath)
	if _, err := configFs.ReadDir(testConfigPath); err != nil {
		t.Errorf("Setup should create config dir\n%s", err)
	}
}
