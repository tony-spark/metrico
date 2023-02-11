package config

import "testing"

func TestParseConfigFileParameter(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			"-short",
			[]string{"-c", "test.json"},
			"test.json",
			false,
		},
		{
			"-short=",
			[]string{"-c=test.json"},
			"test.json",
			false,
		},
		{
			"-long",
			[]string{"-config", "test.json"},
			"test.json",
			false,
		},
		{
			"-long=",
			[]string{"-config=test.json"},
			"test.json",
			false,
		},
		{
			"--long",
			[]string{"--config", "test.json"},
			"test.json",
			false,
		},
		{
			"missed1",
			[]string{"-v", "-c"},
			"",
			true,
		},
		{
			"missed2",
			[]string{"--config="},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfigFileParameter(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfigFileParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseConfigFileParameter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
