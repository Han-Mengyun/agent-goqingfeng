package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainReturnsWhenAPIKeyIsMissing(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "")
	input, err := os.CreateTemp(t.TempDir(), "input")
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	if _, err := input.WriteString("你好\n"); err != nil {
		t.Fatal(err)
	}
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	previousInput, previousOutput := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = input, writer
	defer func() {
		os.Stdin, os.Stdout = previousInput, previousOutput
		writer.Close()
	}()

	main()
	writer.Close()
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(output), "DEEPSEEK_API_KEY is required") {
		t.Fatalf("expected configuration error, got %q", output)
	}
	if strings.Contains(string(output), "开始对话") {
		t.Fatalf("chat started without an API key: %q", output)
	}
}
