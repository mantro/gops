package main

import "testing"

func TestRelPath(t *testing.T) {

	parentPath := "/Users/john/git/coreapp"
	subPath := "/Users/john/git/coreapp/server/assets"

	result := RelPath(parentPath, subPath)

	expected := "./server/assets"

	if result != "./server/assets" {
		t.Fatalf("RelPath(%q, %q) should be %q, but is %q", parentPath, subPath, expected, result)
	}
}
