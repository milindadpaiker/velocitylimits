package io

import (
	"bufio"
	"context"
	"log"
	"os"
	"sync"

	"github.com/milind/velocitylimit/config"
)

//File type implmenents ingester and Sink interface
type File struct {
	file *os.File
}

//NewInputFile returns File implmentation for Ingester interface
func NewInputFile(fName string) (*File, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	return &File{file: f}, nil
}

//NewOutputFile returns File implmentation for Sink interface
func NewOutputFile(fName string) (*File, error) {
	var flag int
	if config.AppendMode {
		flag = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	} else {
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	f, err := os.OpenFile(fName, flag, 0755)
	if err != nil {
		return nil, err
	}
	return &File{file: f}, nil
}

//Write implements Write() of Sink interface
func (f *File) Write(ctx context.Context, ch <-chan string, wg *sync.WaitGroup) {

	bf := bufio.NewWriter(f.file)
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
func (f *File) Read(ctx context.Context, ch chan<- string, wg *sync.WaitGroup) {
	//defer wg.Done()
	scanner := bufio.NewScanner(f.file)

	for scanner.Scan() {
		str := scanner.Text()
		ch <- str
	}
	close(ch)
}
