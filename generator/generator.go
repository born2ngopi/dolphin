package generator

import (
	"context"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func Generate(prompt, model string) (string, error) {

	if model == "" {
		// default model is llama2
		model = "llama2"
	}

	llm, err := ollama.New(
		ollama.WithModel(model),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(1),
		// llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		// 	fmt.Print(string(chunk))
		// 	return nil
		// }),
	)
	if err != nil {
		log.Fatal(err)
	}

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

			return strings.Join(lines, "\n"), nil
		}

	}

	return completion, nil
}

func findTestFunctionIndex(lines []string) int {
	for i, line := range lines {
		if strings.Contains(line, "func Test") {
			return i
		}
	}
	return -1
}
