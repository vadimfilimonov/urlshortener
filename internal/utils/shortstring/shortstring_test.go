package shortstring

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	if got := Generate(); len(got) != 6 {
		t.Errorf("Generate() = %v, want %v", len(got), 6)
	}
}
