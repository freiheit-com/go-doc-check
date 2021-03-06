package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const modeMonoRepo = "monorepo"
const modeApp = "app"

func main() {

	if len(os.Args) != 3 {
		printUsageAndDie()
	}

	mode := os.Args[1]
	if mode != modeMonoRepo && mode != modeApp {
		printUsageAndDie()
	}

	repoPath := os.Args[2]
	if strings.TrimSpace(repoPath) == "" {
		printUsageAndDie()
	}

	checker := Checker{
		repoPath: repoPath,
		reporter: &CLIReporter{},
	}

	err := runCheck(mode, &checker)
	if err != nil {
		panic(err)
	}
	if checker.reporter.FoundIssues() > 0 {
		fmt.Printf("Found %d issues, please check the output\n", checker.reporter.FoundIssues())
		os.Exit(1)
	}
}

func runCheck(mode string, checker *Checker) error {
	switch mode {
	case modeMonoRepo:
		return checker.checkMonoRepo()
	case modeApp:
		return checker.checkApp()
	default:
		printUsageAndDie()
	}
	return nil
}

type Reporter interface {
	Report(string)
	FoundIssues() int
}

type CLIReporter struct {
	foundIssue int
}

func (r *CLIReporter) Report(message string) {
	r.foundIssue++
	fmt.Println(message)
}

func (r *CLIReporter) FoundIssues() int {
	return r.foundIssue
}

type Checker struct {
	repoPath string
	reporter Reporter
}

func (c *Checker) checkMonoRepo() error {
	c.checkReadme("")                      //top-level readme
	err := c.checkReadmeSubfolders("apps") //apps are probably not go projects and we require a Readme
	if err != nil {
		return err
	}
	err = c.checkPackageDocSubfolders("pkg")
	if err != nil {
		return err
	}
	err = c.checkPackageDocSubfolders("services")
	if err != nil {
		return err
	}

	err = c.checkGoFileDoc("pkg")
	if err != nil {
		return err
	}
	err = c.checkGoFileDoc("services")
	if err != nil {
		return err
	}
	return nil
}

func (c *Checker) checkApp() error {
	c.checkReadme("")
	err := c.checkPackageDocSubfolders("app")
	if err != nil {
		return err
	}
	err = c.checkGoFileDoc("app")
	if err != nil {
		return err
	}
	return nil
}

func (c *Checker) checkReadme(folder string) {
	err := fileWithContentExists(c.readme(folder))
	if err != nil {
		c.reporter.Report(err.Error())
	}
}

func (c *Checker) checkPackageDoc(folder string) {
	err := fileWithContentExists(c.packageDoc(folder))
	if err != nil {
		c.reporter.Report(err.Error())
	}
}

func (c *Checker) goFileHasDocComment(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < 2; i++ {
		if !strings.HasPrefix(lines[i], "//") {
			c.reporter.Report(fmt.Sprintf("%s does not contain a file comment!", path))
			return nil
		}
	}

	return nil
}

func (c *Checker) checkPackageDocSubfolders(folder string) error {
	checkFolder := filepath.Join(c.repoPath, folder)
	files, err := ioutil.ReadDir(checkFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			c.checkPackageDoc(filepath.Join(folder, file.Name()))
		}
	}

	return nil
}

func (c *Checker) checkReadmeSubfolders(folder string) error {
	checkFolder := filepath.Join(c.repoPath, folder)
	files, err := ioutil.ReadDir(checkFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			c.checkReadme(filepath.Join(folder, file.Name()))
		}
	}

	return nil
}

func (c *Checker) checkGoFileDoc(subfolder string) error {
	root := filepath.Join(c.repoPath, subfolder)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return fmt.Errorf("root folder %s does not exist!", root)
	}

	return filepath.Walk(root, func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") && file.Name() != "doc.go" {
			err := c.goFileHasDocComment(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *Checker) readme(folder string) string {
	return filepath.Join(c.repoPath, folder, "README.md")
}

func (c *Checker) packageDoc(folder string) string {
	return filepath.Join(c.repoPath, folder, "doc.go")
}

func fileWithContentExists(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist!", file)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("%s exists, but has no content!", file)
	}
	return nil
}

func printUsageAndDie() {
	fmt.Printf("Usage: go-doc-check {monorepo|app} <path-to-repo>\n")
	os.Exit(1)
}
