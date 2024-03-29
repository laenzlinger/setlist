package repertoire

import (
	"testing"

	_ "github.com/laenzlinger/setlist/internal/testinginit"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		band string
	}
	tests := []struct {
		name      string
		args      args
		want      Repertoire
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "good",
			args: args{band: "Band"},
			want: Repertoire{
				columns: []string{"Title", "Year", "Copyright", "Description"},
				songs: []Song{
					{Title: "On the Alamo"},
					{Title: "Frankie and Johnnie"},
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
			for i, song := range got.songs {
				song.TableRow = nil
				got.songs[i] = song
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
