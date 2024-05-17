package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// FileGetContentsTool is a tool for getting the contents of a file
type FileGetContentsTool struct{}

func (fgct *FileGetContentsTool) Execute(input json.RawMessage) (json.RawMessage, error) {
	var fgc FileGetContentsRequest
	if err := json.Unmarshal(input, &fgc); err != nil {
		return nil, err
	}
	output := ExecuteFileGetContents(fgc)
	return json.Marshal(output)
}

func (fgct *FileGetContentsTool) LoadDefinition() ParamsDefinition {
	return paramsDefExecuteFileGetContents()
}

func ExecuteFileGetContents(fgc FileGetContentsRequest) FileGetContentsResponse {
	content, err := os.ReadFile(fgc.Filename)
	if err != nil {
		return FileGetContentsResponse{Content: fmt.Sprintf("Error reading file: %v", err)}
	}
	return FileGetContentsResponse{Content: string(content)}
}
func paramsDefExecuteFileGetContents() ParamsDefinition {
	return ParamsDefinition{
		Type: "object",
		Properties: &FileProperties{
			Filename: Filename{
				Type:        "string",
				Description: "The path to the file to read",
			},
		},
		Required: []string{"filename"},
	}
}

type ParamsDefinition struct {
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
	Required   []string    `json:"required"`
}



type FileProperties struct {
	Filename Filename `json:"filename"`
}



type Filename struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}



// FileGetContentsRequest is a struct for file contents requests
type FileGetContentsRequest struct {
	Filename string `json:"filename"`
}

// FileGetContentsResponse is a struct for file contents responses
type FileGetContentsResponse struct {
	Content string `json:"content"`
}
