package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/born2ngopi/dolphin/parser"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate a new unit test",
		Run: func(cmd *cobra.Command, args []string) {

			var opt parser.Option
			pwd, _ := os.Getwd()
			var configFilePath = filepath.Join(pwd, "dolphin.yml")
			_, err := os.Stat(filepath.Join(pwd, "dolphin.yml"))
			if err != nil {
				_, err = os.Stat(filepath.Join(pwd, "dolphin.yaml"))
				if err != nil {
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

					opt = parser.Option{
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
				}
				configFilePath = filepath.Join(pwd, "dolphin.yaml")
			} else {
				b, err := os.ReadFile(configFilePath)
				if err != nil {
					log.Fatalf("read %s error :%v", configFilePath, err)
				}

				if err := yaml.Unmarshal(b, &opt); err != nil {
					log.Fatalf("unmarshal %s error :%v", configFilePath, err)
				}
			}

			if err := parser.GenerateTest(opt); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "init a new dolphin.yml file",
		Run: func(cmd *cobra.Command, args []string) {
			pwd, _ := os.Getwd()
			var configFilePath = filepath.Join(pwd, "dolphin.yml")
			_, err := os.Stat(filepath.Join(pwd, "dolphin.yml"))
			if err == nil {
				fmt.Println("dolphin.yml already exists")
				os.Exit(1)
			}

			content := `
root-dir: "."
# set directory where you want to generate test
dir: "internal/app"
# set function name if you want to generate test for specific function
func: "GetUser"
# set file name if you want to generate test for specific file
file: "user.go"
# library for mocking
mock-lib: "gomock"
# folder/path location of mock file
mock-path: "internal/app/mock"
output: "internal/app"
# set model
model: ""
# llm is a language model or AI model
# posible value ["ollama", "openai", "gemini"]
llm: "ollama"
# host llm, if u use ollama on local u can use localhost:8080
# or if you use openai or gemini you can put host url on here
llm-host: "http://localhost:8080"

# token llm for authentication
# llm-token:

# if you want to debug prompt, set this to true
# debug: false`

			if err := os.WriteFile(configFilePath, []byte(content), 0644); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Dolphin",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Dolphin v1.0.0")
		},
	}
)

func init() {
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
}

func Execute() {
	rootCmd := &cobra.Command{Use: "dolphin"}

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
