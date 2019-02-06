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

func TestGenerateAliasDefault(t *testing.T) {
	assert := assert.New(t)

	alias, err := generateAlias()
	assert.Nil(err, "generateAlias should not return an error")
	assert.Contains(alias, fmt.Sprintf("function %s()", "hermes"))
	assert.Contains(alias, `$HERMES_BIN $@
	local EXIT_STATUS=$?`, "alias should capture the exit code from the hermes binaray")
	assert.Contains(alias, "return $EXIT_STATUS", "alias should return exit status from binary")
	assert.Contains(alias, fmt.Sprintf("%s%s", testConfigPath, testTargetFile), "generateAlias should have the correct target path")
}

func TestGenerateAliasWithName(t *testing.T) {
	assert := assert.New(t)

	testAliasName := "testAlias"
	viper.Set("alias_name", testAliasName)
	alias, err := generateAlias()
	assert.Nil(err, "generateAlias should not return an error")
	assert.Contains(alias, fmt.Sprintf("function %s()", testAliasName))
}
