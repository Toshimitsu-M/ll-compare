package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	anthropic_option "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go"
	openai_option "github.com/openai/openai-go/option"
)

type QueryRequest struct {
	Message string `json:"message"`
}

type QueryResponse struct {
	Claude string `json:"claude"`
	GPT    string `json:"gpt"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/api/query", handleQuery)

	port := "8080"
	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Message == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var (
		claudeResp string
		gptResp    string
		claudeErr  string
		gptErr     string
		wg         sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		text, err := callClaude(req.Message)
		if err != nil {
			claudeErr = err.Error()
			return
		}
		claudeResp = text
	}()

	go func() {
		defer wg.Done()
		text, err := callGPT(req.Message)
		if err != nil {
			gptErr = err.Error()
			return
		}
		gptResp = text
	}()

	wg.Wait()

	errMsg := ""
	if claudeErr != "" {
		errMsg += "Claude error: " + claudeErr + " "
	}
	if gptErr != "" {
		errMsg += "GPT error: " + gptErr
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(QueryResponse{
		Claude: claudeResp,
		GPT:    gptResp,
		Error:  errMsg,
	})
}

func callClaude(message string) (string, error) {
	client := anthropic.NewClient(
		anthropic_option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
	)

	msg, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeOpus4_7,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
		},
	})
	if err != nil {
		return "", err
	}

	for _, block := range msg.Content {
		if t, ok := block.AsAny().(anthropic.TextBlock); ok {
			return t.Text, nil
		}
	}
	return "", nil
}

func callGPT(message string) (string, error) {
	client := openai.NewClient(
		openai_option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4o,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", nil
}
