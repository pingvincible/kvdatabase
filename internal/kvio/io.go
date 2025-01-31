package kvio

import (
	"bufio"
	"strings"
)

type ReadWriter struct {
	readWriter *bufio.ReadWriter
}

func NewReadWriter(reader *bufio.Reader, writer *bufio.Writer) *ReadWriter {
	return &ReadWriter{
		readWriter: bufio.NewReadWriter(reader, writer),
	}
}

func (r *ReadWriter) ReadLine() (string, error) {
	text, err := r.readWriter.ReadString('\n')
	if err != nil {
		return "", err
	}

	return text, nil
}

func (r *ReadWriter) Write(text string) error {
	text = strings.TrimSpace(text)

	err := r.write(text)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReadWriter) WriteLine(text string) error {
	text = strings.TrimSpace(text)

	err := r.write(text + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (r *ReadWriter) write(text string) error {
	_, err := r.readWriter.WriteString(text)
	if err != nil {
		return err
	}

	err = r.readWriter.Flush()
	if err != nil {
		return err
	}

	return nil
}
