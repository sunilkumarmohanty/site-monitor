package monitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestMonitor_Run(t *testing.T) {
	type fields struct {
		input    []byte
		waitTime int64
	}
	type wantResult struct {
		count int64
		found bool
	}
	tests := []struct {
		name        string
		fields      fields
		wantResults map[string]*wantResult
	}{
		{
			name: "Match Found",
			fields: fields{
				input:    []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"http://inject_from_test", "what_to_check": "response"}]}`),
				waitTime: 10,
			},
			wantResults: map[string]*wantResult{
				"1": &wantResult{
					count: 4,
					found: true,
				},
			},
		},
		{
			name: "Match Not Found",
			fields: fields{
				input:    []byte(`{"checking_period": 3, "tasks":[{"ref_id":"1", "url":"http://inject_from_test", "what_to_check": "missing"}]}`),
				waitTime: 2,
			},
			wantResults: map[string]*wantResult{
				"1": &wantResult{
					count: 1,
					found: false,
				},
			},
		},
		{
			name: "Multiple sites",
			fields: fields{
				input:    []byte(`{"checking_period": 3, "tasks":[{"ref_id":"1", "url":"http://inject_from_test", "what_to_check": "missing"},{"ref_id":"2", "url":"http://inject_from_test","checking_period":6, "what_to_check": "missing"}]}`),
				waitTime: 10,
			},
			wantResults: map[string]*wantResult{
				"1": &wantResult{
					count: 4,
					found: true,
				},
				"2": &wantResult{
					count: 2,
					found: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wantReqCount int64
			handler := func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt64(&wantReqCount, int64(1))
				w.Write([]byte("response from body"))
			}
			server := httptest.NewServer(http.HandlerFunc(handler))
			writer := &bytes.Buffer{}
			m, _ := New(tt.fields.input, writer, "", "")

			for i := range m.work.Tasks {
				m.work.Tasks[i].URL = server.URL
			}

			go func() {
				time.Sleep(time.Duration(tt.fields.waitTime) * time.Second)
				m.Stop()
			}()
			m.Run()
			gotFoundCount := make(map[string]int64)
			var gotReqCount int64
			finalLog := strings.TrimSuffix(writer.String(), "\n")
			for _, m := range strings.Split(finalLog, "\n") {
				result := &result{}
				json.Unmarshal([]byte(m), result)
				if tt.wantResults[result.RefID] != nil {
					if result.Found == tt.wantResults[result.RefID].found {
						gotFoundCount[result.RefID]++
					}
				}

				gotReqCount++
			}
			for k, v := range gotFoundCount {
				if v != tt.wantResults[k].count {
					t.Errorf("Found count does not match, want %v, got %v", tt.wantResults[k].count, v)
				}
			}
		})
	}
}

func TestMonitor_RunInvalidLocation(t *testing.T) {
	type fields struct {
		input     []byte
		userAgent string
		timeOut   string
		waitTime  int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Non existing location",
			fields: fields{
				input:    []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"http://location", "what_to_check": "response"}]}`),
				waitTime: 2,
				timeOut:  "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			m, _ := New(tt.fields.input, writer, tt.fields.userAgent, tt.fields.timeOut)
			go func() {
				time.Sleep(time.Duration(tt.fields.waitTime) * time.Second)
				m.Stop()
			}()
			m.Run()
			finalLog := strings.TrimSuffix(writer.String(), "\n")
			for _, m := range strings.Split(finalLog, "\n") {
				result := &result{}
				json.Unmarshal([]byte(m), result)
				if result.Found || !(strings.HasPrefix(result.Comments, "Error processing request")) {
					t.Errorf("Comments got: %v, want to begin with: Error processing request", result.Comments)
				}
			}

		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		input     []byte
		userAgent string
		timeOut   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Monitor
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{
				input:     []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"http://validurl", "what_to_check": "response"}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					CheckingPeriod: 3,
					Tasks: []task{{
						RefID:          "1",
						URL:            "http://validurl",
						CheckingPeriod: 3,
						WhatToCheck:    "response",
						UserAgent:      "dummy_user_Agent",
					}},
				},
			},
		},
		{
			name: "Invalid JSON input",
			args: args{
				input:     []byte(`{`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			wantErr: true,
		},
		{
			name: "Missing URL in task JSON",
			args: args{
				input:     []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1", "what_to_check": "response"}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					Tasks: []task{},
				},
			},
		},
		{
			name: "Invalid URL",
			args: args{
				input:     []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"validurl", "what_to_check": "response"}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					Tasks: []task{},
				},
			},
		},
		{
			name: "Missing what_to_check in task JSON",
			args: args{
				input:     []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"http://validurl"}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					Tasks: []task{},
				},
			},
		},
		{
			name: "Task checking_period override",
			args: args{
				input:     []byte(`{"checking_period": 3, "tasks":[{"ref_id": "1","url":"http://validurl", "what_to_check": "response","checking_period": 5}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					CheckingPeriod: 3,
					Tasks: []task{{
						RefID:          "1",
						URL:            "http://validurl",
						CheckingPeriod: 5,
						WhatToCheck:    "response",
						UserAgent:      "dummy_user_Agent",
					}},
				},
			},
		},
		{
			name: "Checking period missing at root and task level",
			args: args{
				input:     []byte(`{"tasks":[{"ref_id": "1","url":"http://validurl", "what_to_check": "response"}]}`),
				userAgent: "dummy_user_Agent",
				timeOut:   "10",
			},
			want: &Monitor{
				work: &work{
					Tasks: []task{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.input, nil, tt.args.userAgent, tt.args.timeOut)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got.work.Tasks) != len(tt.want.work.Tasks) {
					t.Errorf("New() got tasks count = %v, want tasks count %v", len(got.work.Tasks), len(tt.want.work.Tasks))
					return
				}

				for i := range tt.want.work.Tasks {
					// below 2 lines of code are to allow deep equal
					got.work.Tasks[i].httpClient = nil
					got.work.Tasks[i].stop = nil
					if !reflect.DeepEqual(got.work.Tasks[i], tt.want.work.Tasks[i]) {
						t.Errorf("New() = %v, want %v", got.work.Tasks[i], tt.want.work.Tasks[i])
					}
				}
			}
		})
	}
}
