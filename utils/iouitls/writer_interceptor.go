package iouitls

import (
	"io"
)

type WriterInterceptorHandler func([]byte)

type WriterInterceptor struct {
	onWrite WriterInterceptorHandler
	target  io.Writer
}

func NewWriterInterceptor(target io.Writer, onWrite WriterInterceptorHandler) *WriterInterceptor {
	return &WriterInterceptor{
		onWrite: onWrite,
		target:  target,
	}
}

var _ io.Writer = (*WriterInterceptor)(nil)

func (wi *WriterInterceptor) Write(p []byte) (int, error) {
	wi.onWrite(p)
	return wi.target.Write(p)
}
