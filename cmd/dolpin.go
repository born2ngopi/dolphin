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
			// dir, _ := cmd.Flags().GetString("dir")
			// mockDir, _ := cmd.Flags().GetString("mockdir")

			if err := parser.GenerateTest(".", "", "", "", ""); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("generate called")
		},
	}

	generateCmd.Flags().StringP("dir", "d", ".", "Specify the directory")
	generateCmd.Flags().StringP("mockdir", "m", "./mocs", "Specify the mock directory")

	rootCmd.AddCommand(generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
