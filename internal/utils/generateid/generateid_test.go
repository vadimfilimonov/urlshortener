package utils

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	ID := GenerateID()
	actual := len(ID)

	if actual != MaxSizeOfID {
		t.Errorf("Generate() = %v, want %v", actual, MaxSizeOfID)
	}
}
