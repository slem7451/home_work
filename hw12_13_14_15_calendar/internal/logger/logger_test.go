package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	require.NotNil(t, New("error"))
	require.Panics(t, func() { New("1") })
}
