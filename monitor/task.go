package monitor

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	validator "gopkg.in/go-playground/validator.v9"
)

type task struct {
	RefID          string `json:"ref_id"` //uniquely identifies each task
	URL            string `json:"url" validate:"required,url"`
	CheckingPeriod int64  `json:"checking_period" validate:"min=1"`  //frequencey at which check is done in secs
	WhatToCheck    string `json:"what_to_check" validate:"required"` //the text to be checked
	UserAgent      string `json:"user_agent"`                        //user agent is sent in the header
	httpClient     *http.Client
	stop           chan bool
}

func (t *task) prepare(defCheckingPeriod int64, defUserAgent string, timeOut time.Duration) error {
	if t.CheckingPeriod == 0 {
		t.CheckingPeriod = defCheckingPeriod
	}
	if len(t.UserAgent) == 0 {
		t.UserAgent = defUserAgent
	}

	err := validator.New().Struct(t)
	if err != nil {
		log.Printf("error validating task. Task %+v. Error: %v", t, err)
		return err
	}
	t.httpClient = &http.Client{Timeout: timeOut}
	t.stop = make(chan bool, 1)
	return nil
}

func (t *task) run(results chan *result) {
	results <- t.request()
	ticker := time.NewTicker(time.Duration(t.CheckingPeriod) * time.Second)
	for {
		select {
		case <-ticker.C:
			results <- t.request()
		case <-t.stop:
			return
		}
	}
}

func (t *task) request() *result {
	//Only GET method is used now as the scope is just to check the website content
	req, err := http.NewRequest("GET", t.URL, nil)
	if err != nil {
		return &result{
			RefID:    t.RefID,
			Comments: "Unable to make request",
			DateTime: time.Now().UnixNano(),
			URL:      t.URL,
		}
	}

	header := make(http.Header)
	header.Set("User-Agent", t.UserAgent)
	req.Header = header
	start := time.Now().UnixNano()
	resp, err := t.httpClient.Do(req)
	duration := time.Now().UnixNano() - start
	if err != nil {
		return &result{
			RefID:        t.RefID,
			Comments:     fmt.Sprintf("Error processing request: %s", err.Error()),
			DateTime:     time.Now().UnixNano(),
			URL:          t.URL,
			RespDuration: duration,
		}
	}
	defer resp.Body.Close()

	found := t.isMatched(resp)
	return &result{
		RefID:        t.RefID,
		URL:          t.URL,
		StatusCode:   resp.StatusCode,
		RespDuration: duration,
		Found:        found,
		DateTime:     time.Now().UnixNano(),
		Comments:     fmt.Sprintf("Searching: %s", t.WhatToCheck),
	}
}

func (t *task) isMatched(resp *http.Response) bool {
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		switch {
		case tokenType == html.ErrorToken: //reached end of document
			return false
		case tokenType == html.TextToken:
			token := tokenizer.Token()
			if found := strings.Contains(token.Data, t.WhatToCheck); found {
				//only first matched instance checked
				return true
			}
		}
	}
}
