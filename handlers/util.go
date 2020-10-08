package handlers

import (
	"io"
	"sync"
)

type CustomReader interface {
	ReadInjectedWithChannel(buf []byte, signal func(int64)) (int, error)
	Read(p []byte) (n int, err error)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src CustomReader, signal func(int64), byteChannel chan Datapack, id int) (written int64, err error) {
	size := 32 * 1024
	buf := make([]byte, size)
	for {
		nr, er := src.ReadInjectedWithChannel(buf, signal)
		if nr > 0 {
			byteChannel <- Datapack{nr, id}
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				byteChannel <- Datapack{nw, id}
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

type Signal struct {
	Mu    *sync.RWMutex
	Value *bool
}

func (s Signal) ReadSignal(code func(sp *bool) bool) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return code(s.Value)
}

func (s Signal) WriteSignal(code func(sp *bool)) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	code(s.Value)
}

type Value1 struct {
	LatestModified int64
	CanModify      bool
}

type Signal1 struct {
	Mu    *sync.RWMutex
	Value *Value1
}

func (s Signal1) ReadSignal1(code func(sp *Value1) bool) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return code(s.Value)
}

func (s Signal1) WriteSignal1(code func(sp *Value1)) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	code(s.Value)
}
