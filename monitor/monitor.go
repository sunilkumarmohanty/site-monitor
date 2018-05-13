package monitor

import (
	"encoding/json"
	"io"
	"strconv"
	"sync"
	"time"
)

//Monitor is the main worker struct
type Monitor struct {
	results   chan *result
	work      *work
	reporter  *reporter
	userAgent string
	timeOut   time.Duration
}

type work struct {
	CheckingPeriod int64  `json:"checking_period,omitempty"` //the default frequency at which the website will be monitored
	Tasks          []task `json:"tasks,omitempty"`
}

//New instantiates a new Monitor
func New(input []byte, writer io.Writer, userAgent, timeOut string) (*Monitor, error) {
	w := &work{}
	err := json.Unmarshal(input, w)
	if err != nil {
		return nil, err
	}
	var t int64
	t, _ = strconv.ParseInt(timeOut, 10, 64)

	m := &Monitor{
		results:   make(chan *result),
		work:      w,
		userAgent: userAgent,
		timeOut:   time.Duration(t) * time.Second,
	}
	for i := range m.work.Tasks {
		err := m.work.Tasks[i].prepare(m.work.CheckingPeriod, m.userAgent, m.timeOut)
		if err != nil {
			//remove task from list
			m.work.Tasks = append(m.work.Tasks[:i], m.work.Tasks[i+1:]...)
		}
	}

	m.reporter = newReporter(writer, m.results)
	return m, nil
}

//Run starts monitoring the websites
func (m *Monitor) Run() {
	var wg sync.WaitGroup
	go func() {
		m.reporter.run()

	}()

	wg.Add(len(m.work.Tasks))
	for i := range m.work.Tasks {
		go func(task *task) {
			task.run(m.results)
			wg.Done()
		}(&m.work.Tasks[i])
	}
	wg.Wait()

	close(m.results)
	// wait till the reporter finishes reporting its report
	<-m.reporter.done
}

// Stop kills all the tasks gracefully
func (m *Monitor) Stop() {
	for i := range m.work.Tasks {
		m.work.Tasks[i].stop <- true
	}
}
