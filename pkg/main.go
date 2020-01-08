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

	switch mode {
	case modeMonoRepo:
		checker.checkMonoRepo()
	case modeApp:
		checker.checkApp()
	default:
		printUsageAndDie()
	}
}

type Reporter interface {
	Report(string)
}

type CLIReporter struct{}

func (r *CLIReporter) Report(message string) {
	fmt.Println(message)
}

type Checker struct {
	repoPath string
	reporter Reporter
}

func (c *Checker) checkMonoRepo() {
	c.checkReadme("") //top-level readme
	c.checkReadmeSubfolders("app")
	c.checkReadmeSubfolders("pkg")
	c.checkReadmeSubfolders("services")

	c.checkGoFileDoc("pkg")
	c.checkGoFileDoc("services")
}

func (c *Checker) checkApp() {
	c.checkReadme("")
	c.checkReadmeSubfolders("app")
	c.checkGoFileDoc("app")
}

func (c *Checker) checkReadme(folder string) {
	err := fileWithContentExists(c.readme(folder))
	if err != nil {
		c.reporter.Report(err.Error())
	}
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

func (c *Checker) checkGoFileDoc(subfolder string) {

}

func (c *Checker) readme(folder string) string {
	return filepath.Join(c.repoPath, folder, "README.md")
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
