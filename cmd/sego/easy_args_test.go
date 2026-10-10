package main

import (
	"reflect"
	"testing"
)

func TestNormalizeEasyArgs(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{"no arguments opens dashboard", []string{"sego"}, []string{"sego", "open"}},
		{"update alias", []string{"sego", "update"}, []string{"sego", "update", "auto"}},
		{"remote alias", []string{"sego", "remote"}, []string{"sego", "tunnel", "setup"}},
		{"unlock alias", []string{"sego", "unlock"}, []string{"sego", "emergency", "unlock"}},
		{"low mode alias", []string{"sego", "low"}, []string{"sego", "mode", "low-resource"}},
		{"normal mode alias", []string{"sego", "normal"}, []string{"sego", "mode", "standard"}},
		{"check alias", []string{"sego", "check"}, []string{"sego", "status"}},
		{"explicit command unchanged", []string{"sego", "tunnel", "status"}, []string{"sego", "tunnel", "status"}},
		{"flagged command unchanged", []string{"sego", "--config", "custom.json", "doctor"}, []string{"sego", "--config", "custom.json", "doctor"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeEasyArgs(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("normalizeEasyArgs(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}
