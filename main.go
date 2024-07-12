package main

import (
	"log"
	"os"

	"github.com/tjgurwara99/ruby-lsp/handlers"
	"github.com/tjgurwara99/ruby-lsp/rpc"
)

func main() {
	logger := getLogger("/Users/taj/personal/ruby-lsp/log.txt")
	mux := rpc.NewMux(os.Stdin, os.Stdout, logger)
	handler := handlers.New(logger)
	mux.HandleMethod("initialize", handler.Initialize)
	mux.HandleMethod("textDocument/completion", handler.TextCompletion)
	mux.HandleMethod("textDocument/definition", handler.GoToDef)
	mux.HandleNotification("textDocument/didOpen", handler.DidOpenHandler)
	mux.HandleNotification("textDocument/didChange", handler.DidChangeHandler)
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
