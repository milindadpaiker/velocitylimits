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
