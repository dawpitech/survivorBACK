package utils

import (
	"net/http"
	"fmt"
	"log"
	"io"
	"encoding/json"
)

func SendRequest(codeBlock string, req *http.Request, data any, fromJson bool) (any, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error getting response for %s, %s", codeBlock, err)
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200{
		log.Printf("Error in status code for %s: %d", codeBlock, resp.StatusCode)
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


