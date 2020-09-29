package io

import (
	"bufio"
	"context"
	"log"
	"os"
	"sync"
)

//Terminal type implmenents ingester and Sink interface
type Terminal struct {
}

//NewInputTerminal returns terminal implmentation for Ingester interface
func NewInputTerminal() (*Terminal, error) {
	return &Terminal{}, nil
}

//NewOutputTerminal returns terminal implmentation for Sink interface
func NewOutputTerminal() (*Terminal, error) {
	return &Terminal{}, nil
}

//Write implements Write() of Sink interface
func (f *Terminal) Write(ctx context.Context, ch <-chan string, wg *sync.WaitGroup) {

	bf := bufio.NewWriter(os.Stdout)
	defer func() {
		_ = bf.Flush()
		wg.Done()
	}()
	for {
		select {
		case data := <-ch:
			if data == "" {
				return
			}
			//how about bf.WriteString(data+ "\n")?
			_, err := bf.WriteString(data)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = bf.Write([]byte{'\n'})
			if err != nil {
				log.Println(err)
				return
			}
			_ = bf.Flush()
		case <-ctx.Done():
			return
		}
	}
}

//Read implements Write() of Ingester interface
func (f *Terminal) Read(ctx context.Context, ch chan<- string, wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()
		ch <- str
	}
	close(ch)
}
