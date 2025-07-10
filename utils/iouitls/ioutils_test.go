package iouitls_test

import (
	"bytes"
	"pb_launcher/utils/iouitls"
	"testing"
)

func TestWriterInterceptor_Write(t *testing.T) {
	var intercepted []byte
	target := &bytes.Buffer{}
	interceptor := iouitls.NewWriterInterceptor(target, func(p []byte) {
		intercepted = append(intercepted, p...)
	})

	input := "hello world"
	n, err := interceptor.Write([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(input) {
		t.Fatalf("expected %d bytes written, got %d", len(input), n)
	}
	if target.String() != input {
		t.Fatalf("expected target to have %q, got %q", input, target.String())
	}
	if string(intercepted) != input {
		t.Fatalf("expected intercepted to have %q, got %q", input, string(intercepted))
	}
}
