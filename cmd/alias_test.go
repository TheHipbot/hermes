package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testConfigPath string
	testTargetFile string
)

func init() {
	testConfigPath = "/test/.hermes/"
	testTargetFile = ".hermes_test"
	viper.Set("config_path", testConfigPath)
	viper.Set("target_file", testTargetFile)
}

func TestGenerateAlias(t *testing.T) {
	assert := assert.New(t)

	alias, err := generateAlias()
	assert.Nil(err, "generateAlias should not return an error")
	assert.Contains(alias, fmt.Sprintf("%s%s", testConfigPath, testTargetFile), "generateAlias should have the correct target path")
}
