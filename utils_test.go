package btobet

import "testing"

func Test_BtoMobile(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{
			name:  "happy 1",
			phone: "0708113456",
			want:  "0708113456",
		},
		{
			name:  "happy 2",
			phone: "254708113456",
			want:  "0708113456",
		},
		{
			name:  "sad 1",
			phone: "070811345",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BtoMobile(tt.phone)
			if got != tt.want {
				t.Errorf("BtoMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}
