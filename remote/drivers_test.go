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
	}

	for _, t := range expectedDrivers {
		d, err := NewDriver(t)
		s.Nil(err, "Driver should be create without error")
		s.NotNil(d, "Driver should not be nil")
	}
}

func (s *DriverSuite) TestNewDriverNotExists() {
	d, err := NewDriver("blah")
	s.Nil(d, "Driver does not exist")
	s.NotNil(err, "Error should be returned")
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}
