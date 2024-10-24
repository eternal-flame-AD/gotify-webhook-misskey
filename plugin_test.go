package main

import (
	"testing"

	"github.com/gotify/plugin-api"
	"github.com/stretchr/testify/assert"
)

func TestAPICompatibility(t *testing.T) {
	assert.Implements(t, (*plugin.Plugin)(nil), new(MisskeyHookPlugin))
	// Add other interfaces you intend to implement here
	assert.Implements(t, (*plugin.Webhooker)(nil), new(MisskeyHookPlugin))
	assert.Implements(t, (*plugin.Configurer)(nil), new(MisskeyHookPlugin))
	assert.Implements(t, (*plugin.Displayer)(nil), new(MisskeyHookPlugin))
}
