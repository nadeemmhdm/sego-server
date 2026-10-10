package httpapi

import "testing"

func TestPreviewTypeAllowlist(t *testing.T) {
	cases := []struct {
		name     string
		wantOK   bool
		wantKind string
	}{
		{"movie.mp4", true, "video"},
		{"clip.webm", true, "video"},
		{"song.mp3", true, "audio"},
		{"photo.png", true, "image"},
		{"manual.pdf", true, "pdf"},
		{"notes.md", true, "text"},
		{"index.html", false, ""},
		{"vector.svg", false, ""},
		{"program.exe", false, ""},
	}
	for _, tc := range cases {
		spec, ok := previewType(tc.name)
		if ok != tc.wantOK {
			t.Fatalf("previewType(%q) allowed=%v, want %v", tc.name, ok, tc.wantOK)
		}
		if ok && spec.Kind != tc.wantKind {
			t.Fatalf("previewType(%q) kind=%q, want %q", tc.name, spec.Kind, tc.wantKind)
		}
	}
}
