package hash

import (
	"fmt"
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

func TestName(t *testing.T) {
	as := []int{1, 2, 3, 4, 5}
	fmt.Println(as[2:2])
}
