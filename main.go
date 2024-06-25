package main

import (
	"bufio"
	"log"
	"os"

	"github.com/tjgurwara99/ruby-lsp/rpc"
)

func main() {
	logger := getLogger("/Users/taj/personal/ruby-lsp/log.txt")
	defer func() {
		if r := recover(); r != nil {
			logger.Println("Recovered in f", r)
		}
	}()
	logger.Println("started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitFunc)
	logger.Println("split set")
	for scanner.Scan() {
		logger.Println("inside scan")
		msg := scanner.Text()
		err := scanner.Err()
		if err != nil {
			logger.Println(err)
		}
		handleMessage(msg, logger)
	}
	logger.Println("stopped")
}

func handleMessage(msg string, logger *log.Logger) {
	logger.Printf("%s\n", msg)
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to initiate logger")
	}
	return log.New(logfile, "[ruby-lsp] ", log.Ldate|log.Ltime|log.Lshortfile)
}
