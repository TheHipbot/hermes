package remote

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DriverSuite struct {
	suite.Suite
}

func (s *DriverSuite) TestCreateDriverOptions() {
	expectedDrivers := []string{
		"github",
	}

	for _, t := range expectedDrivers {
		d, err := CreateDriver(t)
		s.Nil(err, "Driver should be create without error")
		s.NotNil(d, "Driver should not be nil")
	}
}

func (s *DriverSuite) TestCreateDriverNotExists() {
	d, err := CreateDriver("blah")
	s.Nil(d, "Driver does not exist")
	s.NotNil(err, "Error should be returned")
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(DriverSuite))
}
