package service_test

import (
	"testing"

	"github.com/seshoo/bookFinder/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestNewServices(t *testing.T) {
	deps := service.Deps{
		DpUrlTmp: "http://example.com",
	}

	services := service.NewServices(deps)

	assert.NotNil(t, services)
	assert.NotNil(t, services.Parser)
}
