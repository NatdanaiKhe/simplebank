package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(8)
	wrongPassword := RandomString(8)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	require.NotEmpty(t, hashedPassword)
	require.NoError(t, CheckPassword(hashedPassword, password))
	require.Error(t, CheckPassword(hashedPassword, wrongPassword))

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword, hashedPassword2)
}
