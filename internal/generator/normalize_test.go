package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeCloudProvider(t *testing.T) {
	tC := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{raw: "", want: "none"},
		{raw: "none", want: "none"},
		{raw: "aws", want: "aws"},
		{raw: "gcp", want: "", wantErr: true},
		{raw: "azure", want: "", wantErr: true},
		{raw: "bitbucket", want: "", wantErr: true},
	}

	for _, tc := range tC {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := normalizeCloudProvider(tc.raw)
			if tc.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "unsupported cloud provider")
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
