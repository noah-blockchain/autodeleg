package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpPost(url string, payload map[string]string) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
