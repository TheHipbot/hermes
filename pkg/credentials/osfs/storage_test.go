package osfs

import (
	"testing"

	"github.com/TheHipbot/hermes/pkg/credentials"

	"github.com/stretchr/testify/suite"
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

var (
	testFS billy.Filesystem
)

type FSStorageSuite struct {
	suite.Suite
}

func (suite *FSStorageSuite) SetupTest() {
	testFS = memfs.New()
}

func (suite *FSStorageSuite) TestNewFSStorer() {
	file, err := testFS.Create("test_credentials.yml")
	suite.Nil(err, "Should create test file in memory fs")

	storage := NewFSStorer(file)
	suite.NotNil(storage)
	suite.Equal(file, storage.storer)
}

func (suite *FSStorageSuite) TestPutAndGet() {
	file, err := testFS.Create("test_credentials.yml")
	suite.Nil(err, "Should create test file in memory fs")

	storage := NewFSStorer(file)
	suite.NotNil(storage)
	suite.Equal(file, storage.storer)

	testCred := credentials.Credential{
		Type:  "token",
		Token: "123abc",
	}
	storage.Put("test.com", testCred)

	storage2 := NewFSStorer(file)
	suite.NotNil(storage)
	suite.Equal(file, storage.storer)

	cred, err := storage2.Get("test.com")
	suite.Nil(err, "Error should be nil")
	suite.Equal(testCred, cred, "Test credential should be retrieved from the same file")
}

func TestConfigFSSuite(t *testing.T) {
	suite.Run(t, new(FSStorageSuite))
}
