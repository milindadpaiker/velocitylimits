package io

import (
	"context"
	"sync"
)

//Ingester ...
type Ingester interface {
	Read(context.Context, chan<- string, *sync.WaitGroup)
}

//Sink ...
type Sink interface {
	Write(context.Context, <-chan string, *sync.WaitGroup)
}

//Process processes orchestrates flow. Read from input call validator and then sink to output
func Process(ctx context.Context, input Ingester, output Sink) {

	in := make(chan string)
	out := make(chan string)
	var wg sync.WaitGroup
	ct := context.Background()
	ioctx, iocancel := context.WithCancel(ct)
	wg.Add(1)
	go input.Read(ioctx, in, &wg)
	go output.Write(ioctx, out, &wg)
	defer func() {
		iocancel()
		wg.Wait()
	}()
	for {
		select {
		case data := <-in:
			if data == "" {
				close(out)
				iocancel()
				return
			}

			//validate
			out <- data
		case <-ctx.Done():
			//Ctrl + C hit, shut gracefully
			iocancel()
			return
		}
	}

}
