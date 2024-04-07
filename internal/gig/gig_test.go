package gig

import (
	"testing"

	"github.com/laenzlinger/setlist/internal/config"
	_ "github.com/laenzlinger/setlist/internal/testinginit"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		band config.Band
		gig  string
	}
	tests := []struct {
		name      string
		args      args
		want      Gig
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "good",
			args: args{band: config.Band{Name: "MyBand", Source: "Band"}, gig: "Grand Ole Opry"},
			want: Gig{
				Name: "MyBand @ Grand Ole Opry",
				Sections: []Section{
					{
						Header:     "# Set 1\n\nSay Hello",
						SongTitles: []string{"Frankie and Johnnie", "On the Alamo"},
					},
					{
						Header:     "\n\n# Encore",
						SongTitles: []string{"Nowhere To Go"},
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.band, tt.args.gig)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
