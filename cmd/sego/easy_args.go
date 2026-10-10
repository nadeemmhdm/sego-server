package main

import "os"

// init normalizes a small set of beginner-friendly aliases before flag parsing.
// The existing explicit commands remain supported for administrators and scripts.
func init() {
	os.Args = normalizeEasyArgs(os.Args)
}

func normalizeEasyArgs(args []string) []string {
	out := append([]string(nil), args...)
	if len(out) == 1 {
		return append(out, "open")
	}
	if len(out) != 2 {
		return out
	}

	switch out[1] {
	case "update":
		return []string{out[0], "update", "auto"}
	case "remote":
		return []string{out[0], "tunnel", "setup"}
	case "unlock":
		return []string{out[0], "emergency", "unlock"}
	case "low":
		return []string{out[0], "mode", "low-resource"}
	case "normal":
		return []string{out[0], "mode", "standard"}
	case "check":
		return []string{out[0], "status"}
	}
	return out
}
