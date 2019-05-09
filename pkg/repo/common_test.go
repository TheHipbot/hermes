package repo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClonerSuite struct {
	suite.Suite
}

type clonerSuiteCloner struct{}

func (*clonerSuiteCloner) Clone(path string, opts *CloneOptions) error {
	return nil
}

func (suite *ClonerSuite) TestRegisterCloner() {
	testCloner := &clonerSuiteCloner{}
	testCreator := func() (Cloner, error) {
		return testCloner, nil
	}
	RegisterCloner("test", testCreator)
	creator, ok := creators["test"]
	cloner, err := creator()
	suite.True(ok)
	suite.Equal(testCloner, cloner)
	suite.Nil(err)
}

func (suite *ClonerSuite) TestNewCloner() {
	testCloner := &clonerSuiteCloner{}
	testCreator := func() (Cloner, error) {
		return testCloner, nil
	}
	RegisterCloner("test", testCreator)
	c, err := NewCloner("test")
	suite.Nil(err)
	suite.Equal(c, testCloner)
}

func TestClonerSuite(t *testing.T) {
	suite.Run(t, new(ClonerSuite))
}
