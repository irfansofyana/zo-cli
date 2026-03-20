package api

import (
	"bufio"
	"io"
	"strings"
)

// ReadSSE reads Server-Sent Events from r, calling handler for each data line.
func ReadSSE(r io.Reader, handler func(data string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return nil
		}
		if err := handler(data); err != nil {
			return err
		}
	}
	return scanner.Err()
}
