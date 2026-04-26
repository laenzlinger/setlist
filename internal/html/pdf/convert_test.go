package convert_test

import (
	"testing"

	convert "github.com/laenzlinger/setlist/internal/html/pdf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPageCount(t *testing.T) {
	tests := []struct {
		name string
		pdf  string
		want int
	}{
		{
			name: "single page",
			pdf:  "%PDF-1.4\n1 0 obj<</Type /Page>>endobj\n2 0 obj<</Type /Pages /Kids[1 0 R]/Count 1>>endobj\n",
			want: 1,
		},
		{
			name: "two pages",
			pdf: "%PDF-1.4\n1 0 obj<</Type /Page>>endobj\n" +
				"2 0 obj<</Type /Page>>endobj\n" +
				"3 0 obj<</Type /Pages /Kids[1 0 R 2 0 R]/Count 2>>endobj\n",
			want: 2,
		},
		{
			name: "empty",
			pdf:  "%PDF-1.4\n",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert.PageCount([]byte(tt.pdf))
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
