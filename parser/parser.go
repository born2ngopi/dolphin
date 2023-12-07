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
func GenerateTest(rootDir, dir, funcName, fileDir, mockLib, mockDir, output, model string) error {
	logger := pterm.DefaultLogger.
		WithLevel(pterm.LogLevelTrace)

	_getDir, _ := pterm.DefaultSpinner.Start("Getting directory")
	modulePath, projectDir, err := getDir(rootDir)
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
	//if rootDir != "." {
	dir = filepath.Join(projectDir, dir)
	fileDir = filepath.Join(projectDir, fileDir)
	mockDir = filepath.Join(modulePath, mockDir)
	output = filepath.Join(projectDir, output)

	if funcName != "" {

		err := prepareStruct(projectDir)
		if err != nil {
			return err
		}
		singgleSpinner, _ := pterm.DefaultSpinner.Start("Generate Singgle Unit Test ....")

		logger.Info("reading file to prompt")
		// genereate singgle unit test
		_prompt, packageName, err := readFileToPrompt(filepath.Join(projectDir, fileDir), funcName, modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		promptStr, err := _prompt.Generate()
		if err != nil {
			return err
		}

		err = generateAddWriteTestFile(promptStr, model, output, packageName)
		if err != nil {
			return err
		}

		singgleSpinner.Success("Success generate singgle unit test")
		return nil

	}

	multiSpinner, _ := pterm.DefaultSpinner.Start("Generate Multi Unit Test ....")

	multiSpinner.UpdateText("preparing struct .....")
	err = prepareStruct(projectDir)
	if err != nil {
		return err
	}

	multiSpinner.UpdateText("Generate code completion....")
	// walk through the directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// check if is not file and not .go extention
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		// check if the file is test file
		if strings.Contains(path, "_test.go") {
			return nil
		}
		// parse the file
		_prompt, packageName, err := readFileToPrompt(path, "", modulePath, dir, mockLib, mockDir)
		if err != nil {
			return err
		}

		if _prompt.SourceCode == "" {
			return nil
		}

		promptStr, err := _prompt.Generate()
		if err != nil {
			return err
		}

		fmt.Println(promptStr)
		return nil

		outputPath := strings.Replace(path, ".go", "_test.go", 1)

		err = generateAddWriteTestFile(promptStr, model, outputPath, packageName)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Success create test for %s", path))
		return nil
	})
	if err != nil {
		multiSpinner.Fail("Failed to generate test")
		return err
	}

	multiSpinner.UpdateText("run go imports ....")
	// run goimports
	cmd := exec.Command("goimports", "-w", projectDir)
	err = cmd.Run()
	if err != nil {
		multiSpinner.Fail("Failed to run goimports")
		return err
	}

	multiSpinner.Success("Success generate test")
	return nil
}

func generateAddWriteTestFile(prompt, model, outputPath, packageName string) error {
	codeCompletion, err := generator.Generate(prompt, model)
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

		f.WriteString("package " + packageName + "\n\n")

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
	pwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	modulePath = filepath.Join(pwd, modulePath)
	// check if go.mod exist
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {

		return "", pwd, nil
	}

	modfileBytes, err := os.ReadFile(modulePath)
	if err != nil {
		return "", "", err
	}

	modFile, err := modfile.Parse(modulePath, modfileBytes, nil)
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
