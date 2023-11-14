package parser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/born2ngopi/dolpin/generator"
	"github.com/pterm/pterm"
	"golang.org/x/mod/modfile"
)

// GenerateTest is used to auto generate test for golang code.
func GenerateTest(dir, funcName, fileDir, mockLib, mockDir, output string) error {
	logger := pterm.DefaultLogger.
		WithLevel(pterm.LogLevelTrace)

	_getDir, _ := pterm.DefaultSpinner.Start("Getting directory")
	modulePath, projectDir, err := getDir(dir)
	if err != nil {
		return err
	}
	_getDir.InfoPrinter = &pterm.PrefixPrinter{
		MessageStyle: &pterm.Style{pterm.FgLightBlue},
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgBlack, pterm.BgLightBlue},
			Text:  " USING ",
		},
	}
	_getDir.Info(fmt.Sprintf("Modulepath: %s | Projectdir : %s", modulePath, projectDir))

	if funcName != "" {
		singgleSpinner, _ := pterm.DefaultSpinner.Start("Generate Singgle Unit Test ....")

		logger.Info("reading file to prompt")
		// genereate singgle unit test
		_prompts, err := readFileToPrompt(filepath.Join(projectDir, fileDir), funcName, modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		if len(_prompts) == 0 {
			singgleSpinner.Fail(fmt.Sprintf("function %s not found", funcName))
			return fmt.Errorf("function %s not found", funcName)
		}

		logger.Info("generate code completion....")
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

		singgleSpinner.Success("Success generate singgle unit test")
		return nil

	}

	multiSpinner, _ := pterm.DefaultSpinner.Start("Generate Multi Unit Test ....")
	// walk through the directory
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		// check if is not file and not .go extention
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		multiSpinner.UpdateText("Reading file to prompt....")
		// parse the file
		_prompts, err := readFileToPrompt(path, "", modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		multiSpinner.UpdateText("Generate code completion....")
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
		}
		multiSpinner.UpdateText("Success create test")
		logger.Info(fmt.Sprintf("Success create test for %s", path))
		return nil
	})
	if err != nil {
		multiSpinner.Fail("Failed to generate test")
		return err
	}

	multiSpinner.Success("Success generate test")
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
