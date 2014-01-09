package ext

import (
	"testing"
)

func TestFileRotator(t *testing.T) {
	rotator := NewFileRotator("/var/log", "testfile", "log", 10, 10)
	_, err := rotator.Write([]byte("writing to file"))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
