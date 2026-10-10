package auth

import "testing"

func TestPasswordRoundTrip(t *testing.T) {
	h, e := HashPassword("correct horse battery staple")
	if e != nil {
		t.Fatal(e)
	}
	if !VerifyPassword(h, "correct horse battery staple") {
		t.Fatal("expected valid password")
	}
	if VerifyPassword(h, "wrong password") {
		t.Fatal("accepted wrong password")
	}
}
func TestOpaqueTokensDiffer(t *testing.T) {
	a, ha, e := NewOpaqueToken()
	if e != nil {
		t.Fatal(e)
	}
	b, hb, e := NewOpaqueToken()
	if e != nil {
		t.Fatal(e)
	}
	if a == b || ha == hb {
		t.Fatal("tokens must be unique")
	}
	if HashToken(a) != ha {
		t.Fatal("hash mismatch")
	}
}
