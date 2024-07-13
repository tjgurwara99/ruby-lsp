package index

import (
	"log"
	"testing"
)

func TestIndex(t *testing.T) {
	i := New("./testdata")
	i.Start(log.Default())
	_ = i
}
