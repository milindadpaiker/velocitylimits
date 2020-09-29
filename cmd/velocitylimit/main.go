package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/milind/velocitylimit/config"

	"github.com/milind/velocitylimit/io"
	validate "github.com/milind/velocitylimit/validate"
)

func main() {
	configPtr := flag.String("config", "config.json", "configuration file for the application")
	logPtr := flag.String("log", "velocitylimit.log", "log file for the application")
	boolPtr := flag.Bool("append", false, "only applicable when sink is a file. Appends results to given outfile")
	inFilePtr := flag.String("infile", "input.txt", "ingestion of funds through text file")
	outFilePtr := flag.String("outfile", "output.txt", "output file to send fund status")
	stringPtr := flag.String("backend", "memory", "datastore for the validation service. Options: memory or sqlite")
	stdinPtr := flag.Bool("stdin", false, "When set to true ingestion of funds will be through terminal. Do not set any other ingestion mode if this true.")
	stdoutPtr := flag.Bool("stdout", false, "When set to true output will be sent to the terminal. Do not set any other sink mode if this true")

	var input io.Ingester
	var output io.Sink
	var err error

	flag.Parse()
	config.AppendMode = *boolPtr

	err = loadLogFile(*logPtr)
	if err != nil {
		//Log
		panic(err)
	}
	log.Println("Application started")
	err = loadConfig(*configPtr)
	if err != nil {
		log.Panicf("Failed to load configuration from file %s. Error: %v", *configPtr, err)
	}

	//set ingestion. Use factory pattern in future.
	if *stdinPtr {
		input, err = io.NewInputTerminal()
		if err != nil {
			log.Panicf("Failed to load ingestion file %s. Error: %v", *inFilePtr, err)
		}
	} else {
		input, err = io.NewInputFile(*inFilePtr)
		if err != nil {
			log.Panicf("Failed to load ingestion file %s. Error: %v", *inFilePtr, err)
		}
	}

	if *stdoutPtr {
		output, err = io.NewOutputTerminal()
		if err != nil {
			log.Panicf("Failed to load output file %s. Error: %v", *outFilePtr, err)
		}
	} else {
		output, err = io.NewOutputFile(*outFilePtr)
		if err != nil {
			log.Panicf("Failed to load output file %s. Error: %v", *outFilePtr, err)
		}
	}

	validator, err := validate.NewValidator(*stringPtr)
	if err != nil {
		log.Panicf("Failed to load velocity limits validation module. Error: %v", err)
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
		log.Println("Received Shutdown signal")
		cancel()
	}()

	io.Process(ctx, input, output, validator)
	log.Println("Application winding down")
}

func loadConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&config.Configuration)
}

func loadLogFile(logPath string) error {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}
