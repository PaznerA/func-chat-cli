package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Execute(input json.RawMessage) (json.RawMessage, error)
	LoadDefinition() ParamsDefinition
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
		functionResponse := Execute(functionToCall, json.RawMessage(args))

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

func Execute(funcName string, args json.RawMessage) string {
	resp := "{}"
	functionMap := map[string]Tool{
		"ExecuteGetCurrentWeather": &WeatherTool{},
		"FileGetContents":          &FileGetContentsTool{},
		"Wget":                     &WgetTool{},
	}

	if tool, ok := functionMap[funcName]; ok {
		result, err := tool.Execute(args)
		if err != nil {
			panic(err)
		}
		resp = string(result)
	}
	return resp
}

func functionDefinitionLoader() []openai.FunctionDefinition {
	return []openai.FunctionDefinition{
		{
			Name:        "ExecuteGetCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters:  paramsDefExecuteGetCurrentWeather(),
		},
		{
			Name:        "ExecuteFileGetContents",
			Description: "Get the contents of a file",
			Parameters:  paramsDefExecuteFileGetContents(),
		},

		{
			Name:        "ExecuteWget",
			Description: "Get the contents of an URL address",
			Parameters:  paramsDefExecuteWget(),
		},
	}
}