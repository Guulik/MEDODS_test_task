package service

import (
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_generateRefreshToken(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
		{
			name: "2",
		},
		{
			name: "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				cfg:           nil,
				log:           nil,
				tokenModifier: nil,
				tokenProvider: nil,
			}
			refreshToken, err := s.generateRefreshToken()
			require.NoError(t, err)

			decodedToken, err := base64.StdEncoding.DecodeString(refreshToken)
			require.NoError(t, err)

			encoded := base64.StdEncoding.EncodeToString(decodedToken)
			require.Equal(t, refreshToken, encoded)
		})
	}
}
