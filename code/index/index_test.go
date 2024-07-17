package index

import (
	"fmt"
	"log"
	"testing"
)

func TestIndex(t *testing.T) {
	i := New("./testdata")
	i.Start(log.Default())
	fmt.Println()
}
