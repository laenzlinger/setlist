package gig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSectionHeaderText(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "heading",
			header: "# Set 1\n",
			want:   "Set 1",
		},
		{
			name:   "empty",
			header: "",
			want:   "",
		},
		{
			name:   "no heading",
			header: "just text\n",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{Header: []byte(tt.header)}
			assert.Equal(t, tt.want, s.HeaderText())
		})
	}
}
