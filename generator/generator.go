package generator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms/local"
)

var (
	bin         = os.Getenv("LOCAL_LLM_BIN")
	model       = os.Getenv(("LOCAL_LLM_MODEL"))
	gpuLayers   = os.Getenv(("LOCAL_LLM_NUM_GPU_LAYERS"))
	threads     = os.Getenv(("LOCAL_LLM_NUM_CPU_CORES"))
	contextSize = os.Getenv(("LOCAL_LLM_CONTEXT"))
)

func Generate(prompt string) (string, error) {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory")
	}

	bin := fmt.Sprintf("%s/%s", wd, bin)
	args := fmt.Sprintf("-m %s/%s -t %s --temp 0 -eps 1e-5 -c %s -ngl %s -p",
		wd, model, threads, contextSize, gpuLayers)

	llm, err := local.New(
		local.WithBin(bin),
		local.WithArgs(args),
	)
	if err != nil {
		return "", errors.New("cannot create local LLM")
	}

	completion, err := llm.Call(context.Background(), prompt)
	if err != nil {
		log.Println("Cannot get completion:", err.Error())
		return "", errors.New("cannot get completion")
	}

	// remove the question if it appears in the response
	completion = strings.ReplaceAll(completion, prompt, "")

	if strings.Contains(completion, "```") {
		split := strings.Split(completion, "```")
		if len(split) > 1 {
			return split[1], nil
		}
	}

	return completion, nil
}
