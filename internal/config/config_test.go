package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGigName(t *testing.T) {
	type config func()
	type args struct {
		args []string
	}
	tests := []struct {
		name      string
		args      args
		config    config
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "with arg",
			args:      args{args: []string{"arg-value"}},
			config:    func() {},
			want:      "arg-value",
			assertion: assert.NoError,
		},
		{
			name:      "arg overwrites config",
			args:      args{args: []string{"arg-value"}},
			config:    func() { viper.Set("gig.name", "config-value") },
			want:      "arg-value",
			assertion: assert.NoError,
		},
		{
			name:      "use config",
			args:      args{args: []string{}},
			config:    func() { viper.Set("gig.name", "config-value") },
			want:      "config-value",
			assertion: assert.NoError,
		},
		{
			name:      "error",
			args:      args{args: []string{}},
			config:    func() {},
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.config()
			act, err := GigName(tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, act)
		})
	}
}
