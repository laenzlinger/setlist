package repertoire

import (
	"testing"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/laenzlinger/setlist/internal/song"
	_ "github.com/laenzlinger/setlist/internal/testinginit"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		band config.Band
	}
	tests := []struct {
		name      string
		args      args
		want      Repertoire
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "good",
			args: args{band: config.Band{Source: "Band"}},
			want: Repertoire{
				columns: []string{"Title", "Year", "Description", "Composer", "Arranger", "Duration"},
				songs: []song.Song{
					{Title: "On the Alamo"},
					{Title: "Frankie and Johnnie"},
					{Title: "Nowhere to go"},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.band)
			tt.assertion(t, err)
			got.markdown = nil
			got.source = nil
			got.header = nil
			for i, song := range got.songs {
				song.TableRow = nil
				got.songs[i] = song
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
