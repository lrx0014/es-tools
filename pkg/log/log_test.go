package log

import (
	"testing"
)

func TestClient(t *testing.T) {
	if _, err := CreateClient(); err != nil {
		t.Fatal(err)
	}
}
