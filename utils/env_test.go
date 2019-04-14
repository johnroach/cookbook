package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	assert.Equal(t, GetEnv("FOO", "0.0.1"), "1")
	assert.Equal(t, GetEnv("FOO2", "0.0.1"), "0.0.1")
}