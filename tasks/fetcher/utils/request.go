package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func initClient(codeBlock string, req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error getting response for %s, %s", codeBlock, err)
		return nil
	}
	return resp
}

func SendRequest(codeBlock string, req *http.Request, data any, fromJson bool) (any, error) {

	var resp *http.Response

	for {
		counter := 0
		resp = initClient(codeBlock, req)
		defer resp.Body.Close()
		if resp == nil {
			return nil, fmt.Errorf("Unable to get client for %s", codeBlock)
		}
		if resp.StatusCode != 429 || counter > 10 {
			break
		}
		counter += 1
		time.Sleep(10 * time.Millisecond)
	}
	if resp.StatusCode != 200 {
		return data, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading response body for %s: %s", codeBlock, err)
		return data, err
	}

	if fromJson {
		err = json.Unmarshal(body, data)
		if err != nil {
			log.Printf("Error to unmarshall body for %s: %s", codeBlock, err)
			return nil, err
		}
		return data, nil
	}
	return string(body), nil
}
