package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (s *service) resizeHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Expecting POST request"))
			return
		}

		// we dispatch the `?async` request argument, if not present it defaults to false (process synchronously)
		async, ok := r.URL.Query()["async"]
		processAsync := false

		if !ok {
			log.Println("Async flag wasn't specified")
		} else if async[0] == "true" {
			processAsync = true
		}


		request := resizeRequest{}
		err := json.NewDecoder(io.LimitReader(r.Body, 8*1024)).Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Failed to parse request"))
			return
		}

		// Notice the new `processAsync` argument
		results, err := s.processResizes(request, processAsync)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to process request"))
			return
		}

		data, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal response"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("content-type", "application/json")
		w.Write(data)
	})
}

func (s *service) getImageHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("fetching ", r.URL.String())
		data, ok := s.cache.Get(r.URL.String())
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "image/jpeg")


		// read the cache entry and decide which value to write in the browser based on `ongoingProcessing`
		chEntry := data.(cacheEntry)

		if chEntry.ongoingProcessing {
			fmt.Println("I am stuck here")
			//_, ok := <-chEntry.chMsg


			select {
			case <-chEntry.chMsg:
				fmt.Println("The channel is (%v)", ok)
				// force the cache to get update
			default:
			}

			data, _ := s.cache.Get(r.URL.String())
			chEntry := data.(cacheEntry)

			w.Write(chEntry.terminalValue)

		} else {
			fmt.Println("Regular branch")
			w.Write(chEntry.terminalValue)
		}
	})
}
