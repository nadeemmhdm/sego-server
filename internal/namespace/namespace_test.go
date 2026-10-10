package namespace

import "testing"

func TestUserPathCannotEscapeNamespace(t *testing.T) {
	got := UserPath(7, "../../8/private.txt")
	if got != "__users/7/8/private.txt" {
		t.Fatalf("unexpected path: %s", got)
	}
}

func TestPublicPathHidesNamespace(t *testing.T) {
	if got := PublicPath(7, "__users/7/docs/a.txt"); got != "docs/a.txt" {
		t.Fatalf("unexpected public path: %s", got)
	}
}
