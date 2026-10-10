package config

import "testing"

func TestApplyModeLowResource(t *testing.T) {
	c := Default(t.TempDir())
	c.LowResourceMode = false
	c.UploadConcurrency = 5
	if err := ApplyMode(&c, ModeLowResource); err != nil {
		t.Fatal(err)
	}
	if !c.LowResourceMode || c.UploadConcurrency != 1 {
		t.Fatalf("unexpected low-resource config: %+v", c)
	}
}

func TestApplyModeStandard(t *testing.T) {
	c := Default(t.TempDir())
	if err := ApplyMode(&c, ModeStandard); err != nil {
		t.Fatal(err)
	}
	if c.LowResourceMode || c.UploadConcurrency < 3 {
		t.Fatalf("unexpected standard config: %+v", c)
	}
}

func TestApplyModeRejectsUnknown(t *testing.T) {
	c := Default(t.TempDir())
	if err := ApplyMode(&c, "turbo"); err == nil {
		t.Fatal("expected unknown mode to fail")
	}
}
