package cmd

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestSwaggerCommandRunsGenerator(t *testing.T) {
	called := false
	command := newSwaggerCommand(func(context.Context, io.Writer, io.Writer) error {
		called = true
		return nil
	})
	if err := command.Execute(); err != nil {
		t.Fatalf("execute Swagger command: %v", err)
	}
	if !called {
		t.Fatal("Swagger generator was not called")
	}
}

func TestSwaggerCommandPropagatesGeneratorError(t *testing.T) {
	wanted := errors.New("generator failed")
	command := newSwaggerCommand(func(context.Context, io.Writer, io.Writer) error { return wanted })
	if err := command.Execute(); !errors.Is(err, wanted) {
		t.Fatalf("execute error = %v, want %v", err, wanted)
	}
}
