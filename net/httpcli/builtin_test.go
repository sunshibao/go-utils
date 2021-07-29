package httpcli

import (
	"testing"
)

func TestDefault(t *testing.T) {
	if Default() == nil {
		t.Fatal("Default() return nil")
	}
}
