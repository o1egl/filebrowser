package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := "hello world"
	hashedPassword, err := Password(password)
	require.NoError(t, err)
	require.NotEqual(t, password, hashedPassword)
	require.True(t, CheckPassword(password, hashedPassword))
}
