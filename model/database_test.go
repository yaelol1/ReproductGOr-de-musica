package reproductgor/model

import "testing"

func TestNewDatabase (t *testing.T) {
	expected := "Hello, World!"
	if observed := NewDatabase(); observed != expected {
		t.Fatalf("HelloWorld() = %v, want %v", observed, expected)
	}
