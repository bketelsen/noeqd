package noeqd

import (
	"fmt"
	"testing"
)

func TestServe(t *testing.T) {
	g, err := NewGenerator(1, 1)
	id, err := g.Get()
	if err != nil {
		t.Error(err)
	}
	if id == 0 {
		t.Errorf("Expected an id, got %d", id)
	}
	fmt.Println(id)

}
