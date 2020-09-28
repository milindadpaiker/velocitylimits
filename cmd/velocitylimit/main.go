package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/milind/velocitylimit/config"

	"github.com/milind/velocitylimit/io"
)

func main() {
	boolPtr := flag.Bool("recover", false, "recover last run")
	flag.Parse()
	config.RecoverMode = *boolPtr
	fmt.Println("recover:", *boolPtr)

	var input io.Ingester
	var output io.Sink
	var err error
	input, err = io.NewInputFile("input.txt")
	if err != nil {
		//log
		panic(err)
	}
	output, err = io.NewOutputFile("output.txt")
	if err != nil {
		//log
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		//Listen for Ctrl + C, so programm can shut gracefully
		sig := make(chan os.Signal, 1)
		signal.Notify(sig,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		<-sig
		fmt.Println("Received Shutdown")
		cancel()
	}()
	io.Process(ctx, input, output)

	//should wait for graceful shutdown
}
