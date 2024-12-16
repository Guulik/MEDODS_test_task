package ip

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestGetIPAddress(t *testing.T) {

	tests := []struct {
		name       string
		headers    map[string]string
		remoteAddr string
		want       string
	}{
		{
			name:       "X-forwarded",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.1"},
			remoteAddr: "192.0.1.1:22",
			want:       "203.0.113.1",
		},
		{
			name:       "no X-Forwarded",
			headers:    map[string]string{},
			remoteAddr: "192.0.1.1:961",
			want:       "192.0.1.1",
		},
		{
			name:       "no X-Forwarded another port",
			headers:    map[string]string{},
			remoteAddr: "192.0.1.1:555",
			want:       "192.0.1.1",
		},
		{
			name:       "Empty RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "",
			want:       "",
		},
		{
			name:       "invalid remote addr",
			headers:    map[string]string{},
			remoteAddr: "ooidgi",
			want:       "",
		},
		{
			name:       "invalid X-Forwarded header",
			headers:    map[string]string{"X-Forwarded-For": "volvo"},
			remoteAddr: "",
			want:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header:     http.Header{},
				RemoteAddr: tt.remoteAddr,
			}
			for header, value := range tt.headers {
				req.Header.Set(header, value)
			}

			ip := GetIPAddress(req)

			require.Equal(t, tt.want, ip)
		})
	}
}
