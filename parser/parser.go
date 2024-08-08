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

var multiSpinner *pterm.SpinnerPrinter

type Option struct {
	RootDir     string
	Dir         string
	FuncName    string
	FileDir     string
	MockLib     string
	MockDir     string
	Output      string
	Model       string
	Llm         string
	LlmHost     string
	LlmToken    string
	DebugPrompt bool
}

// GenerateTest is used to auto generate test for golang code.
func GenerateTest(opt Option) error {

	logger := pterm.DefaultLogger.
		WithLevel(pterm.LogLevelTrace)

	_getDir, _ := pterm.DefaultSpinner.Start("Getting directory")
	modulePath, projectDir, err := getDir(opt.RootDir)
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

	// check user have goimports tools
	checkGoImports(projectDir)

	_getDir.Info(fmt.Sprintf("Modulepath: %s | Projectdir : %s", modulePath, projectDir))
	//if rootDir != "." {
	dir := filepath.Join(projectDir, opt.Dir)
	fileDir := filepath.Join(projectDir, opt.FileDir)
	mockPath := filepath.Join(modulePath, opt.MockDir)
	mockDir := filepath.Join(projectDir, opt.MockDir)
	gitDir := filepath.Join(projectDir, "./.git")
	output := filepath.Join(projectDir, opt.Output)

	if opt.FuncName != "" {

		err := prepareStruct(projectDir)
		if err != nil {
			return err
		}
		singgleSpinner, _ := pterm.DefaultSpinner.Start("Generate Singgle Unit Test ....")

		logger.Info("reading file to prompt")
		// genereate singgle unit test

		conf := Config{
			Path:          fileDir,
			FuncName:      opt.FuncName,
			ModulePath:    modulePath,
			Dir:           dir,
			MockLib:       opt.MockLib,
			MockDir:       mockPath,
			ExistingTests: nil,
		}

		_prompt, packageName, err := readFileToPrompt(conf)
		if err != nil {
			return err
		}

		promptStr, err := _prompt.Generate()
		if err != nil {
			return err
		}

		if opt.DebugPrompt {
			fmt.Println(promptStr)
			return nil
		}

		err = generateAddWriteTestFile(promptStr, output, packageName, opt)
		if err != nil {
			return err
		}

		singgleSpinner.Success("Success generate singgle unit test")
		return nil

	}

	multiSpinner, _ = pterm.DefaultSpinner.Start("Generate Multi Unit Test ....")

	multiSpinner.UpdateText("preparing struct .....")
	err = prepareStruct(projectDir)
	if err != nil {
		return err
	}

	multiSpinner.UpdateText("Check existing unit test....")

	// get list of existing test function
	var existingTests = make(map[string]string)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, errArg error) error {
		// skip mocks and .git folder
		if path == mockDir || path == gitDir {
			return filepath.SkipDir
		}

		// only read test file
		if !strings.Contains(path, "_test.go") {
			return nil
		}

		funcNames := getListFunctionName(path)
		if len(funcNames) != 0 {
			for _, name := range funcNames {
				existingTests[name] = path
			}
		}

		return nil
	})
	if err != nil {
		multiSpinner.Fail("Failed to generate test")
		return err
	}

	multiSpinner.UpdateText("Generate code completion....")

	multiSpinner.Stop()
	// walk through the directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, errArg error) error {
		// skip mocks folder
		if path == mockDir || path == gitDir {
			return filepath.SkipDir
		}

		// check if is not file and not .go extention
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		// check if the file is test file
		if strings.Contains(path, "_test.go") {
			return nil
		}

		conf := Config{
			Path:          path,
			FuncName:      "",
			ModulePath:    modulePath,
			Dir:           dir,
			MockLib:       opt.MockLib,
			MockDir:       mockPath,
			ExistingTests: existingTests,
		}

		// parse the file
		_prompt, packageName, err := readFileToPrompt(conf)
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

		if opt.DebugPrompt {
			// debug prompt
			fmt.Println(promptStr)
			return nil
		}

		outputPath := strings.Replace(path, ".go", "_test.go", 1)

		err = generateAddWriteTestFile(promptStr, outputPath, packageName, opt)
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

func generateAddWriteTestFile(prompt, outputPath, packageName string, opt Option) error {
	codeCompletion, err := generator.Generate(generator.Option{
		Prompt: prompt,
		Model:  opt.Model,
		Llm:    opt.Llm,
		Token:  opt.LlmToken,
		Host:   opt.LlmHost,
	})
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

func checkGoImports(projectDir string) {
	cmd := exec.Command("goimports", "-w", projectDir)
	err := cmd.Run()
	if err != nil {

		fmt.Printf(`Dolpin need goimports for importing package after generate
You can install manualy by visit this site

https://pkg.go.dev/golang.org/x/tools/cmd/goimports
`)
		os.Exit(1)

	}
}
