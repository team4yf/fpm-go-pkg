package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA256(t *testing.T) {
	sum := Sha256Encode("public")

	assert.Equal(t, sum, "efa1f375d76194fa51a3556a97e641e61685f914d446979da50a551a4333ffd7", "error")
}

func TestGetHostname(t *testing.T) {
	hname := GetHostname()

	t.Logf("hostname: %v", hname)

}
