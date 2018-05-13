package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type result struct {
	RefID        string `json:"ref_id,omitempty"`
	URL          string `json:"url,omitempty"`
	StatusCode   int    `json:"status_code,omitempty"`
	RespDuration int64  `json:"resp_duration,omitempty"` //Response duration in Nano seconds
	Found        bool   `json:"found,omitempty"`         // Whethere the match was found or not
	Comments     string `json:"comments,omitempty"`      //Comments
	DateTime     int64  `json:"date_time,omitempty"`     // When was the task completed and sent to reporter
}

type reporter struct {
	writer  io.Writer
	results chan *result
	done    chan bool
}

func newReporter(writer io.Writer, results chan *result) *reporter {
	if writer == nil {
		writer = os.Stdout
	}
	return &reporter{
		writer:  writer,
		results: results,
		done:    make(chan bool, 1),
	}
}

func (r *reporter) run() {
	for res := range r.results {
		r.print(res)
	}
	//result channel is closed and the reporter reports its work is done
	r.done <- true
}

func (r *reporter) print(res *result) {
	msg, err := json.Marshal(res)
	if err != nil {
		log.Println("Unable to marshal message")
	}
	fmt.Fprintf(r.writer, "%+v\n", string(msg))
}
