package cmd

import (
	"fmt"
	"os"

	"github.com/born2ngopi/dolphin/parser"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{Use: "dolphin"}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a new unit test",
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("dir")
			funcName, _ := cmd.Flags().GetString("func")
			fileDir, _ := cmd.Flags().GetString("file")
			mockLib, _ := cmd.Flags().GetString("mock-lib")
			mockDir, _ := cmd.Flags().GetString("mock-path")
			output, _ := cmd.Flags().GetString("output")
			model, _ := cmd.Flags().GetString("model")
			llm, _ := cmd.Flags().GetString("llm")
			llmHost, _ := cmd.Flags().GetString("llm-host")
			llmToken, _ := cmd.Flags().GetString("llm-token")
			rootDir, _ := cmd.Flags().GetString("root-dir")
			debugPrompt, _ := cmd.Flags().GetBool("debug")

			opt := parser.Option{
				RootDir:     rootDir,
				Dir:         dir,
				FuncName:    funcName,
				FileDir:     fileDir,
				MockLib:     mockLib,
				MockDir:     mockDir,
				Output:      output,
				Model:       model,
				Llm:         llm,
				LlmHost:     llmHost,
				LlmToken:    llmToken,
				DebugPrompt: debugPrompt,
			}

			if err := parser.GenerateTest(opt); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	generateCmd.Flags().StringP("dir", "d", ".", "Specify the directory")
	generateCmd.Flags().StringP("func", "f", "", "Specify the function name")
	generateCmd.Flags().StringP("file", "F", "", "Specify the file directory")
	generateCmd.Flags().StringP("mock-lib", "m", "", "Specify the mock library")
	generateCmd.Flags().StringP("mock-path", "M", "./mocks", "Specify the mock path")
	generateCmd.Flags().StringP("output", "o", "", "Specify the output directory")
	generateCmd.Flags().String("model", "", "Specify the model")
	generateCmd.Flags().String("llm", "ollama", "Specify the llm")
	generateCmd.Flags().String("llm-host", "", "Specify the llm host")
	generateCmd.Flags().StringP("llm-token", "T", "", "Specify the llm token")
	generateCmd.Flags().StringP("root-dir", "r", ".", "Specify the root directory")
	generateCmd.Flags().BoolP("debug", "D", false, "Debug prompt")

	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
