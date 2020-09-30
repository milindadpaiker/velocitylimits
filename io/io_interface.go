package io

import (
	"context"
	"sync"
)

//Ingester Every source will have implement this interface
type Ingester interface {
	Read(context.Context, chan<- string, *sync.WaitGroup)
}

//Sink Every sink will have implement this interface
type Sink interface {
	Write(context.Context, <-chan string, *sync.WaitGroup)
}
