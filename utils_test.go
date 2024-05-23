package btobet

import "testing"

func Test_parseMobile(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		want    string
		wantErr bool
	}{
		{
			name:    "happy 1",
			phone:   "0708113456",
			want:    "0708113456",
			wantErr: false,
		},
		{
			name:    "happy 2",
			phone:   "254708113456",
			want:    "0708113456",
			wantErr: false,
		},
		{
			name:    "sad 1",
			phone:   "070811345",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMobile(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMobile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}
