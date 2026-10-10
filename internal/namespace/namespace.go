package namespace

import (
	"fmt"
	pathpkg "path"
	"strings"
)

func UserPath(userID int64, rel string) string {
	base := fmt.Sprintf("__users/%d", userID)
	rel = strings.ReplaceAll(rel, "\\", "/")
	clean := strings.TrimPrefix(pathpkg.Clean("/"+rel), "/")
	if clean == "." || clean == "" {
		return base
	}
	return base + "/" + clean
}

func PublicPath(userID int64, rel string) string {
	base := fmt.Sprintf("__users/%d", userID)
	rel = strings.Trim(strings.ReplaceAll(rel, "\\", "/"), "/")
	if rel == base {
		return ""
	}
	return strings.TrimPrefix(rel, base+"/")
}
