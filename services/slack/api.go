package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func OpenDialog(input DialogOptions, headers map[string]string) (result map[string]interface{}, err error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", SlackAPIURL, ActionOpenDialog),
		bytes.NewReader(data),
	)

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respData))
	if err != nil {
		return nil, err
	}

	json.Unmarshal(respData, &result)

	return result, nil
}
