package updater

import (
	"path/filepath"
	"testing"
)

func TestNewerVersionComparison(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"0.2.1", "0.2.0", true},
		{"1.0.0", "0.9.9", true},
		{"v2.0.0", "1.9.9", true},
		{"0.2.0", "0.2.0", false},
		{"0.1.9", "0.2.0", false},
	}
	for _, tc := range cases {
		if got := newer(tc.latest, tc.current); got != tc.want {
			t.Fatalf("newer(%q, %q) = %v, want %v", tc.latest, tc.current, got, tc.want)
		}
	}
}

func TestSafeJoinBlocksArchiveEscape(t *testing.T) {
	root := t.TempDir()
	if _, err := safeJoin(root, "bin/sego"); err != nil {
		t.Fatalf("expected safe relative path: %v", err)
	}
	if _, err := safeJoin(root, filepath.Join("..", "escape")); err == nil {
		t.Fatal("expected parent-directory archive escape to be rejected")
	}
}
