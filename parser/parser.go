package parser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/born2ngopi/dolpin/generator"
	"golang.org/x/mod/modfile"
)

// GenerateTest is used to auto generate test for golang code.
func GenerateTest(dir, funcName, fileDir, mockLib, mockDir, output string) error {
	modulePath, projectDir, err := getDir(dir)
	if err != nil {
		return err
	}

	if funcName != "" {
		// genereate singgle unit test
		_prompts, err := readFileToPrompt(filepath.Join(projectDir, fileDir), funcName, modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		if len(_prompts) == 0 {
			return fmt.Errorf("function %s not found", funcName)
		}

		for _, prompt := range _prompts {

			promptStr, err := prompt.Generate()
			if err != nil {
				return err
			}
			fmt.Println(promptStr)

			err = generateAddWriteTestFile(promptStr, output)
			if err != nil {
				return err
			}
		}

	}

	// walk through the directory
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		// check if is not file and not .go extention
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// parse the file
		_prompts, err := readFileToPrompt(path, "", modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		for _, _prompt := range _prompts {

			promptStr, err := _prompt.Generate()
			if err != nil {
				return err
			}

			outputPath := strings.Replace(path, ".go", "_test.go", 1)

			err = generateAddWriteTestFile(promptStr, outputPath)
			if err != nil {
				return err
			}

			// TODO: log success generate test
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func generateAddWriteTestFile(prompt string, outputPath string) error {
	codeCompletion, err := generator.Generate(prompt)
	if err != nil {
		return err
	}

	var (
		f *os.File
	)

	// check if the file exist
	if _, err = os.Stat(outputPath); os.IsNotExist(err) {
		// create the file
		f, err = os.Create(outputPath)
		if err != nil {
			return err
		}

	} else {

		f, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}

	defer f.Close()

	f.WriteString(codeCompletion)

	return nil
}

func getDir(dir string) (modulePath string, projectDir string, err error) {
	modulePath = filepath.Join(dir, "go.mod")
	// check if go.mod exist
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {

		// get pwd
		pwd, err := os.Getwd()
		if err != nil {
			return "", "", err
		}

		return "", pwd, nil
	}

	modFile, err := modfile.Parse(modulePath, nil, nil)
	if err != nil {
		return "", "", err
	}

	modulePath = modFile.Module.Mod.Path

	// get project dir

	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", modulePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("failed to run 'go list' command: %v", err)
	}

	projectDir = strings.TrimSpace(string(output))

	return modulePath, projectDir, nil
}
