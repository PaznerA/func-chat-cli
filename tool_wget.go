package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// WgetTool is a tool for making HTTP requests to a given URL
type WgetTool struct{}

func (wt *WgetTool) Execute(input json.RawMessage) (json.RawMessage, error) {
	var wr WgetRequest
	if err := json.Unmarshal(input, &wr); err != nil {
		return nil, err
	}
	output, err := ExecuteWget(wr)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (wt *WgetTool) LoadDefinition() ParamsDefinition {
	return paramsDefExecuteWget()
}

func ExecuteWget(wr WgetRequest) (WgetResponse, error) {
	req, err := http.NewRequest(wr.Method, wr.URL, nil)
	if err != nil {
		return WgetResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WgetResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WgetResponse{}, err
	}

	return WgetResponse{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil
}

func paramsDefExecuteWget() ParamsDefinition {
	return ParamsDefinition{
		Type: "object",
		Properties: &WgetProperties{
			URL: URL{
				Type:        "string",
				Description: "The URL to fetch",
			},
			Method: Method{
				Type: "string",
				Enum: []string{"GET", "POST", "PUT", "DELETE"},
			},
		},
		Required: []string{"url", "method"},
	}
}

type WgetProperties struct {
	URL    URL    `json:"url"`
	Method Method `json:"method"`
}

type URL struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Method struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}

// WgetRequest is a struct for wget requests
type WgetRequest struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

// WgetResponse is a struct for wget responses
type WgetResponse struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}
