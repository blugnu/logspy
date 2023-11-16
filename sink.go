package logspy

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

var sink = &bytes.Buffer{}

// Sink returns an io.Writer reference to the buffer that logspy
// uses to capture log output.
func Sink() io.Writer {
	return sink
}

// Strings returns a slice containing all non-empty entries in the log,
// where each 'entry' is delimited in the log by '\n' separators.
//
// Entries are trimmed of any leading and trailing whitespace.
func Strings() (result []string) {
	for _, e := range strings.Split(Content(), "\n") {
		if len(strings.TrimSpace(e)) > 0 {
			result = append(result, e)
		}
	}
	return result
}

// JsonObjects returns the log as a slice of maps, one map per log entry.
func JsonObjects() (result []map[string]interface{}, err error) {
	r := strings.NewReader(sink.String())
	d := json.NewDecoder(r)
	for {
		var o map[string]interface{}
		if err := d.Decode(&o); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, nil
}

// Reset resets the sink, clearing any captured logs..
func Reset() {
	sink.Reset()
}

// Content returns a string containing the entire contents of the
// log captured since the most recent reset.
func Content() string {
	return sink.String()
}
