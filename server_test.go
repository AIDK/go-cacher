package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandleGetUser(t *testing.T) {

	s := NewServer()
	ts := httptest.NewServer(http.HandlerFunc(s.handleGetUser))
	wg := &sync.WaitGroup{} // we're using a WaitGroup to wait for all goroutines to finish

	for i := 0; i < 1000; i++ {
		wg.Add(1) // we're adding 1 to the WaitGroup counter for each goroutine
		// we're using a goroutine to simulate concurrent requests
		go func(i int) {
			// we're generating a id between 1 and 100
			id := i%100 + 1
			// we're creating a url like: http://localhost:8080/?id=1
			url := fmt.Sprintf("%s/?id=%d", ts.URL, id)

			// we're making a GET request to our server with
			// the above URL and getting the response
			resp, err := http.Get(url)
			if err != nil {
				t.Error(err)
			}

			// we're decoding the response body into a User struct
			// and checking if it's valid or not and printing it
			user := &User{}
			if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
				t.Error(err)
			}

			fmt.Printf("%+v\n", user)
			wg.Done() // we call Done() on the WaitGroup for each goroutine
		}(i) // we're passing i to the goroutine

		// we're sleeping for 1 millisecond to give the goroutines time to finish
		time.Sleep(time.Millisecond * 1)
	}

	wg.Wait() // we're waiting for all goroutines to finish
	fmt.Println("database hit count:", s.hit)
}
