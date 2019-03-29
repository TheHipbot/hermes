package cmd

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionHandler(t *testing.T) {
	assert := assert.New(t)

	data := &VersionData{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		BuildTime: "Wed Mar 27 21:56:42 CDT 2019",
	}

	result, err := generateVersionOutput(versionTemplate, data)
	assert.NotEqual("", result)
	assert.Nil(err)
	assert.Regexp(regexp.MustCompile("Version:\\s+v1\\.0\\.0"), result, "Version output should have a Version field with the correct version")
	assert.Regexp(regexp.MustCompile("Git Commit:\\s+abc123"), result, "Version output should have a Git Commit field with the correct commit")
	assert.Regexp(regexp.MustCompile("Built:\\s+Wed Mar 27 21:56:42 CDT 2019"), result, "Version output should have a Built field with the correct time")
}

func TestVersionHandlerWithBadTemplate(t *testing.T) {
	assert := assert.New(t)

	data := &VersionData{
		Version:   "v1.0.0",
		GitCommit: "abc123",
		BuildTime: "Wed Mar 27 21:56:42 CDT 2019",
	}

	result, err := generateVersionOutput("{{ define ... }}", data)
	assert.Equal("", result)
	assert.Equal(errParseTemplate, err)
}

func TestVersionHandlerWithBadData(t *testing.T) {
	assert := assert.New(t)

	_, err := generateVersionOutput(versionTemplate, nil)
	assert.Equal(errExecuteTemplate, err)
}
