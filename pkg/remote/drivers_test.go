package remote

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DriverSuite struct {
	suite.Suite
}

func (s *DriverSuite) TestNewDriverOptions() {
	expectedDrivers := []string{
		"github",
		"gitlab",
	}

	for _, t := range expectedDrivers {
		d, err := NewDriver(t, &DriverOpts{
			Auth: &Auth{
				Token: "abcd123",
				Type:  "token",
			},
		})
		s.Nil(err, "Driver should be create without error")
		s.NotNil(d, "Driver should not be nil")
	}
}

func (s *DriverSuite) TestNewDriverNotExists() {
	d, err := NewDriver("blah", &DriverOpts{})
	s.Nil(d, "Driver does not exist")
	s.NotNil(err, "Error should be returned")
}

func (s *DriverSuite) TestRegisterDriver() {
	d, err := NewDriver("test", &DriverOpts{})
	s.Nil(d, "Driver does not exist")
	s.NotNil(err, "Error should be returned")

	RegisterDriver("test", gitlabCreator)
	d, err = NewDriver("test", &DriverOpts{})
	s.NotNil(d, "Driver should exist")
	gl := d.(*GitLab)
	s.NotNil(gl, "Should be a gitlab driver")
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}
