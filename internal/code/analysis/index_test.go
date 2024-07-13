package analysis_test

import (
	"log"
	"os"
	"testing"

	"github.com/tjgurwara99/ruby-lsp/internal/code/analysis"
)

func TestIndexStart(t *testing.T) {
	logger := log.New(os.Stdin, "INSIDE TEST", log.Default().Flags())
	index := analysis.New("/Users/taj/github/github", logger)
	index.Start()
	t.Log(index)
}
