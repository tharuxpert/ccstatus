package cmd

import "testing"

func TestDisplayVersionPrefixesReleaseVersion(t *testing.T) {
	oldVersion := version
	t.Cleanup(func() {
		version = oldVersion
	})

	version = "1.4.1"

	if got := DisplayVersion(); got != "v1.4.1" {
		t.Fatalf("expected v1.4.1, got %q", got)
	}
}

func TestDisplayVersionDoesNotDoublePrefixReleaseVersion(t *testing.T) {
	oldVersion := version
	t.Cleanup(func() {
		version = oldVersion
	})

	version = "v1.4.1"

	if got := DisplayVersion(); got != "v1.4.1" {
		t.Fatalf("expected v1.4.1, got %q", got)
	}
}

func TestDisplayVersionDoesNotPrefixDev(t *testing.T) {
	oldVersion := version
	t.Cleanup(func() {
		version = oldVersion
	})

	version = "dev"

	if got := DisplayVersion(); got == "vdev" {
		t.Fatal("expected dev version not to be prefixed with v")
	}
}
