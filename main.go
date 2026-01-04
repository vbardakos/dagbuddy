package main

import (
	"log"
	"os"

	"github.com/vbardakos/dagbuddy/lsp"
	"github.com/vbardakos/dagbuddy/rpc"
)

func main() {
	log := newLogger("lsp.log")

	if err := lsp.RunLspServer(log, rpc.WithOnError(LogError(log))); err != nil {
		panic(err)
	}

	log.Println("server closes gloriously!!!")
}

func LogError(log *log.Logger) func(err error) {
	return func(err error) {
		log.Printf("error: %s\n", err)
	}
}

func newLogger(fname string) *log.Logger {
	out, err := os.OpenFile(fname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)

	if err != nil {
		panic("Bro the file is wrong!")
	}

	return log.New(out, "", log.Ldate|log.Ltime|log.Lshortfile)
}
