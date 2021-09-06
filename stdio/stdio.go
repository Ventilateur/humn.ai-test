package stdio

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"humn.ai/phan/logger"
	"humn.ai/phan/models"
)

// Read reads from io.Reader (which can be os.Stdin) line by line, unmarshal lines into models.Input objects and
// push them to a buffered channel. Read terminates when EOF is reached.
func Read(r io.Reader, buf chan models.Input) {
	defer close(buf)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Bytes()
		input := models.Input{}
		if err := json.Unmarshal(raw, &input); err != nil {
			logger.Errorf("failed to unmarshal input, skipping: %s", err.Error())
			continue
		}
		select {
		case buf <- input:
			// Push input value to buffer
		default:
			logger.Fatalf("failed to push to input channel: buffer overflowed")
			return
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatalf("failed to read input stream: %s", err.Error())
	}
}

// Write reads models.Output objects from a buffered channel, marshals them into strings of compact JSON, then
// writes to an io.Writer (which can be os.Stdout). Write terminates when the channel is closed.
func Write(w io.Writer, buf chan models.Output, finished chan bool) {
	for out := range buf {
		raw, err := json.Marshal(&out)
		if err != nil {
			logger.Errorf("failed to unmarshal data, skipping: %s", err.Error())
			continue
		}
		compact := new(bytes.Buffer)
		if err := json.Compact(compact, raw); err != nil {
			logger.Errorf("failed to compact json, skipping: %s", err.Error())
			continue
		}
		_, err = fmt.Fprintln(w, compact.String())
		if err != nil {
			logger.Errorf("failed to print, skipping: %s", err.Error())
			continue
		}
	}
	logger.Infof("Output channel is closed, nothing more to be read")
	finished <- true
}
