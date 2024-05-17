package mai2n

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Execute(input interface{}) (output interface{}, err error)
}

type Response interface{}

func main() {
	reader := bufio.NewReader(os.Stdin)
	messages := []openai.ChatCompletionMessage{}

	for {
		fmt.Print("User: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" || message == ":q" || message == "ZZ" {
			break
		}
		if message == "debug" {
			fmt.Println("-----------------")
			debugMessages, _ := json.MarshalIndent(messages, "", " ")
			fmt.Println(string(debugMessages))
			fmt.Println("-----------------")
			continue
		}
		if message == "file" {
			fmt.Println("-----------------")
			file, _ := json.MarshalIndent(messages, "", " ")
			fileName := "messages" + time.Now().Format("20060102150405") + ".json"
			_ = os.WriteFile(fileName, file, 0644)
			fmt.Println("saved!")
			fmt.Println("-----------------")
			continue
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})

		fmt.Print("Bot: ")
		messages = RunConversation(messages)
		fmt.Println("")
	}
}

func RunConversation(messages []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	apikey := readApiKey("apikey.txt")
	functionDefs := functionDefinitionLoader()
	// debug print content of functionDefs
	fmt.Println("-----------------")
	functionDefsDebug, _ := json.MarshalIndent(functionDefs, "", " ")
	fmt.Println(string(functionDefsDebug))
	fmt.Println("-----------------")

	client := openai.NewClient(apikey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo0613,
			Messages:  messages,
			Functions: functionDefs,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion round 1 error: %v\n", err)
		return nil
	}

	if resp.Choices[0].FinishReason == "function_call" {
		messages = append(messages, resp.Choices[0].Message)

		functionToCall := resp.Choices[0].Message.FunctionCall.Name
		args := resp.Choices[0].Message.FunctionCall.Arguments
		functionResponse := Execute(functionToCall, args)

		message3 := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleFunction,
			Content: functionResponse,
			Name:    functionToCall,
		}

		messages = append(messages, message3)
	}

	resp, err = client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo0613,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion round 2 error: %v\n", err)
		return nil
	}
	println(resp.Choices[0].Message.Content)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})
	return messages
}

func readApiKey(fname string) string {
	fileContent, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(fileContent))
}

func Execute(funcName string, args string) string {
	resp := "{}"
	functionMap := map[string]Tool{
		"ExecuteGetCurrentWeather": &WeatherTool{},
		"FileGetContents":          &FileGetContentsTool{},
	}

	if tool, ok := functionMap[funcName]; ok {
		// Get the type of the input expected by the tool
		toolType := reflect.TypeOf(tool).Elem()
		executeMethod, found := toolType.MethodByName("Execute")
		if !found {
			panic("method Execute not found")
		}
		inputType := executeMethod.Type.In(1)

		// Create a new instance of the input type
		inputValue := reflect.New(inputType).Interface()

		// Unmarshal JSON arguments into the input instance
		if err := json.Unmarshal([]byte(args), inputValue); err != nil {
			panic(err)
		}

		// Execute the tool with the input value
		result, err := tool.Execute(inputValue)
		if err != nil {
			panic(err)
		}

		// Marshal the result to JSON
		r, _ := json.Marshal(result)
		resp = string(r)
	}
	return resp
}

// WeatherTool is a tool for getting the current weather
type WeatherTool struct{}

func (wt *WeatherTool) Execute(input interface{}) (output interface{}, err error) {
	wr := input.(WeatherRequest)
	return ExecuteGetCurrentWeather(wr), nil
}

// FileGetContentsTool is a tool for getting the contents of a file
type FileGetContentsTool struct{}

func (fgct *FileGetContentsTool) Execute(input interface{}) (output interface{}, err error) {
	fgc := input.(FileGetContentsRequest)
	return FileGetContents(fgc), nil
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

func FileGetContents(fgc FileGetContentsRequest) FileGetContentsResponse {
	content, err := os.ReadFile(fgc.Filename)
	if err != nil {
		return FileGetContentsResponse{Content: fmt.Sprintf("Error reading file: %v", err)}
	}
	return FileGetContentsResponse{Content: string(content)}
}

func functionDefinitionLoader() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        "ExecuteGetCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters:  paramsDefExecuteGetCurrentWeather(),
		},
		{
			Name:        "FileGetContents",
			Description: "Get the contents of a file",
			Parameters:  paramsDefExecuteFileGetContents(),
		},
	}
}

func paramsDefExecuteGetCurrentWeather() ParamsDefinition {
	return ParamsDefinition{
		Type: "object",
		Properties: WeatherProperties{
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

func paramsDefExecuteFileGetContents() ParamsDefinition {
	return ParamsDefinition{
		Type: "object",
		Properties: FileProperties{
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

type WeatherProperties struct {
	Location Location `json:"location"`
	Unit     Unit     `json:"unit"`
}

type FileProperties struct {
	Filename Filename `json:"filename"`
}

type Location struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Unit struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}

type Filename struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// WeatherRequest is a struct for weather requests
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

// FileGetContentsRequest is a struct for file contents requests
type FileGetContentsRequest struct {
	Filename string `json:"filename"`
}

// FileGetContentsResponse is a struct for file contents responses
type FileGetContentsResponse struct {
	Content string `json:"content"`
}
