package main

import (
	"log"
	"os"

	"github.com/tjgurwara99/ruby-lsp/handlers"
	"github.com/tjgurwara99/ruby-lsp/rpc"
)

func main() {
	mux := rpc.NewMux(os.Stdin, os.Stdout)
	mux.HandleMethod("initialize", handlers.Initialize)
	logger := getLogger("/Users/taj/personal/ruby-lsp/log.txt")
	for {
		if err := mux.Process(); err != nil {
			logger.Println(err)
			return
		}
	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to initiate logger")
	}
	return log.New(logfile, "[ruby-lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
