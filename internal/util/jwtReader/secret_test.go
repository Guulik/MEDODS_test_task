package jwtReader

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadJWTSecret(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want []byte
	}{
		{
			name: "local",
			env:  "local",
			want: []byte("default"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, LoadJWTSecret(tt.env))
		})
	}
}
