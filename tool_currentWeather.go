package main

import (
	"encoding/json"
)

// WeatherTool is a tool for getting the current weather
type WeatherTool struct{}

func (wt *WeatherTool) Execute(input json.RawMessage) (json.RawMessage, error) {
	var wr WeatherRequest
	if err := json.Unmarshal(input, &wr); err != nil {
		return nil, err
	}
	output := ExecuteGetCurrentWeather(wr)
	return json.Marshal(output)
}

type Location struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Unit struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}


// WeatherRequest is a struct for weather requests
type WeatherProperties struct {
	Location Location `json:"location"`
	Unit     Unit     `json:"unit"`
}
type WeatherRequest struct {
	Location string `json:"location"`
	Unit     string `json:"unit"`
}

// WeatherResponse is a struct for weather responses
type WeatherResponse struct {
	Location    string   `json:"location"`
	Temperature string   `json:"temperature"`
	Unit        string   `json:"unit"`
	Forecast    []string `json:"forecast"`
}

func paramsDefExecuteGetCurrentWeather() ParamsDefinition {
	return ParamsDefinition{
		Type: "object",
		Properties: &WeatherProperties{
			Location: Location{
				Type:        "string",
				Description: "The city and state, e.g. San Francisco, CA",
			},
			Unit: Unit{
				Type: "string",
				Enum: []string{"celsius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
	}
}

func ExecuteGetCurrentWeather(wr WeatherRequest) WeatherResponse {
	println(wr.Location)
	if len(wr.Unit) == 0 {
		wr.Unit = "celsius"
	}
	temperature := "72"
	forecast := []string{"sunny", "windy"}
	return WeatherResponse{
		Location:    wr.Location,
		Temperature: temperature,
		Unit:        wr.Unit,
		Forecast:    forecast,
	}
}

func (wt *WeatherTool) LoadDefinition() ParamsDefinition {
	return paramsDefExecuteGetCurrentWeather()
}