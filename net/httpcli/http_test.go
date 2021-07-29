package httpcli

import (
	"testing"
)

func TestDefaultHTTPClient(t *testing.T) {
	if DefaultHTTPClient() == nil {
		t.Fatal("DefaultHTTPClient() return nil")
	}
}

func TestNewDefaultHTTPClient(t *testing.T) {
	if NewDefaultHTTPClient() == nil {
		t.Fatal("NewDefaultHTTPClient() return nil")
	}
}
