package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	h := NewHasher("secret")
	password := "hello world"
	hashedPassword, err := h.Password(password)
	require.NoError(t, err)
	require.NotEqual(t, password, hashedPassword)
	require.True(t, h.CheckPassword(password, hashedPassword))
}
