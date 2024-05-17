# OpenAI GPT CLI chat with function calling in GO
This project is a command-line interface (CLI) chat application that leverages the OpenAI API to generate responses and execute custom functions based on user input. The application supports custom tools that can be executed based on the conversation context, including fetching weather information, reading file contents, and making HTTP requests.

## Features

- Interactive CLI for chatting with the OpenAI API.
- Supports custom functions that can be called during the conversation.
- Includes tools for fetching weather data, reading file contents, and making HTTP requests.
- Allows saving and debugging conversation history.

## Prerequisites

- Go 1.16 or later
- OpenAI API key

## Installation

1. Clone the repository:

```sh
git clone https://github.com/PaznerA/func-chat-cli
cd func-chat-cli
```
2. Install dependencies:
```sh

go mod tidy
```
Create a file named apikey.txt in the project root directory and add your OpenAI API key to it.

### Usage
Run the CLI application:

```sh

go run main.go
```
### Commands

`exit, :q, ZZ` - Exit the application.

`debug` - Print the conversation history in JSON format for debugging.

`file` - Save the conversation history to a JSON file with a timestamped filename.

### Example Tools:

#### Weather Tool - fetches the current weather for a given location.

Request:

```json
{
  "location": "San Francisco, CA",
  "unit": "celsius"
}
```
Response:

```json
{
  "location": "San Francisco, CA",
  "temperature": "72",
  "unit": "celsius",
  "forecast": ["sunny", "windy"]
}
```
#### FileGetContents Tool

Reads the contents of a specified file.

Request:

```json
{
  "filename": "example.txt"
}
```
Response:

```json
{
  "content": "File content here..."
}
```
Wget Tool
Makes an HTTP request to a specified URL.

Request:

```json
{
  "url": "https://api.example.com/data",
  "method": "GET"
}
```
Response:

```json
{
  "status_code": 200,
  "body": "Response body here..."
}
```
### Extending the Application

To add a new tool, implement the Tool interface and add the tool to the functionMap in main.go. 

See the WeatherTool, FileGetContentsTool, and WgetTool implementations for reference.

Example Tool Implementation
```go
type ExampleTool struct{}

func (et *ExampleTool) Execute(input json.RawMessage) (json.RawMessage, error) {
    var req ExampleRequest
    if err := json.Unmarshal(input, &req); err != nil {
        return nil, err
    }
    // Process the request and generate a response
    res := ExampleResponse{ /* ... */ }
    return json.Marshal(res)
}

func (et *ExampleTool) LoadDefinition() ParamsDefinition {
    return ParamsDefinition{
        Type: "object",
        Properties: &ExampleProperties{
            /* ... */
        },
        Required: []string{ /* ... */ },
    }
}
```


## Feel free to send a PR
