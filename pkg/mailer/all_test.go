package mailer

import "testing"

// Unit tests only: go test -v -short
// Integration test only: go test -run Integration

// TestSomething is a base unit test reference.
func TestSomething(t *testing.T) {
	t.Parallel()
	t.Skip("Skiping unit test at the moment.")
}

// TestSendIntegration is a base unit test reference.
func TestSendIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skiping integration test.")
	}
	if false {
		t.Errorf("Status failed: %t", true)
	}
}
