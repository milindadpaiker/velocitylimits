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
	stringPtr := flag.String("backend", "memory", "datastore for the validation service. Choice: memory or sqlite")

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

	input, err = io.NewInputFile(*inFilePtr)
	if err != nil {
		log.Panicf("Failed to load ingestion file %s. Error: %v", *inFilePtr, err)
	}
	// input, err = io.NewInputTerminal(*inFilePtr)
	// if err != nil {
	// 	log.Panicf("Failed to load ingestion file %s. Error: %v", *inFilePtr, err)
	// }
	output, err = io.NewOutputFile(*outFilePtr)
	if err != nil {
		log.Panicf("Failed to load output file %s. Error: %v", *outFilePtr, err)
	}
	// output, err = io.NewOutputTerminal(*outFilePtr)
	// if err != nil {
	// 	log.Panicf("Failed to load output file %s. Error: %v", *outFilePtr, err)
	// }
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
	defer file.Close()
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
