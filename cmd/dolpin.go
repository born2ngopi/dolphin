package cmd

import (
	"fmt"
	"os"

	"github.com/born2ngopi/dolpin/parser"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{Use: "dolpin"}

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

			if err := parser.GenerateTest(dir, funcName, fileDir, mockLib, mockDir, output, model); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("generate called")
		},
	}

	generateCmd.Flags().StringP("dir", "d", ".", "Specify the directory")
	generateCmd.Flags().StringP("func", "f", "", "Specify the function name")
	generateCmd.Flags().StringP("file", "F", "", "Specify the file directory")
	generateCmd.Flags().StringP("mock-lib", "m", "", "Specify the mock library")
	generateCmd.Flags().StringP("mock-path", "M", "./mocs", "Specify the mock path")
	generateCmd.Flags().StringP("output", "o", "", "Specify the output directory")
	generateCmd.Flags().String("model", "llama2", "Specify the model")

	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
