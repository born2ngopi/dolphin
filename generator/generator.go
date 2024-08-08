package generator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"google.golang.org/api/option"
)

type Option struct {
	Llm    string
	Model  string
	Token  string
	Prompt string
	Host   string
}

func generateWithOllama(opt Option) (string, error) {
	if opt.Model == "" {
		// default model is llama2
		opt.Model = "codegemma:7b"
	}

	ollamaOpts := []ollama.Option{
		ollama.WithModel(opt.Model),
	}

	if opt.Host != "" {
		ollamaOpts = append(ollamaOpts, ollama.WithServerURL(opt.Host))
	}

	llm, err := ollama.New(
		ollamaOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	completion, err := llm.Call(ctx, opt.Prompt,
		llms.WithTemperature(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	return normalizeCompletionResponse(completion, opt.Prompt), nil
}

func generateWithGemini(opt Option) (string, error) {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(opt.Token))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if opt.Model == "" {
		opt.Model = "gemini-1.5-flash"
	}
	_model := client.GenerativeModel(opt.Model)

	resp, err := _model.GenerateContent(ctx, genai.Text(opt.Prompt))
	if err != nil {
		return "", err
	}

	completion := formatResponse(resp)

	return normalizeCompletionResponse(completion, opt.Prompt), nil

}

func generateWithOpenAI(opt Option) (string, error) {
	ctx := context.Background()
	client := openai.NewClient(opt.Token)
	req := openai.ChatCompletionRequest{
		Model: opt.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a Senior Go developer who is writing a good test",
			},
			{
				Role:    "user",
				Content: opt.Prompt,
			},
		},
		MaxTokens:        1000,
		Temperature:      0.5,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Stream:           false,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	var completion string
	if len(resp.Choices) > 0 {
		for _, choice := range resp.Choices {
			if choice.Message.Role == "assistant" {
				completion = choice.Message.Content
				break
			}
		}
	}

	return normalizeCompletionResponse(completion, opt.Prompt), nil
}

func Generate(opt Option) (string, error) {

	switch opt.Llm {
	case "gemini":
		return generateWithGemini(opt)
	case "openai":
		return generateWithOpenAI(opt)
	default:
		return generateWithOllama(opt)
	}
}

func normalizeCompletionResponse(completion, prompt string) string {
	// remove the question if it appears in the response
	completion = strings.ReplaceAll(completion, prompt, "")
	if strings.Contains(completion, "```") {
		split := strings.Split(completion, "```")
		if len(split) > 1 {

			lines := strings.Split(split[1], "\n")

			// Find the index where the test function starts
			startIndex := findTestFunctionIndex(lines)

			// Remove lines before the test function
			if startIndex >= 0 {
				lines = lines[startIndex:]
			}

			return strings.Join(lines, "\n")
		}

	}

	return completion
}

func findTestFunctionIndex(lines []string) int {
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "func Test") {
			return i
		}
	}
	return -1
}

func formatResponse(resp *genai.GenerateContentResponse) string {
	var formattedContent strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formattedContent.String()
}
