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

func TestSliceContains(t *testing.T) {

	slice1 := []string{"One", "Two", "Three"}

	if !SliceContains(slice1, "One") {
		t.Fatalf("001")
	}

	if !SliceContains(slice1, "Two") {
		t.Fatalf("001")
	}

	if !SliceContains(slice1, "Three") {
		t.Fatalf("001")
	}

	if SliceContains(slice1, "Four") {
		t.Fatalf("001")
	}

	if SliceContains(slice1, "") {
		t.Fatalf("001")
	}
}
