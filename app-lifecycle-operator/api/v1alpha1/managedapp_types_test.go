package v1alpha1

import "testing"

func TestGroupVersion(t *testing.T) {
	if GroupVersion.Group != "platform.example.com" || GroupVersion.Version != "v1alpha1" {
		t.Fatalf("unexpected group/version: %s", GroupVersion.String())
	}
}
