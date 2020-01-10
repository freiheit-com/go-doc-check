package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckReadme(t *testing.T) {
	tcs := []struct {
		desc                   string
		repoPath               string
		folder                 string
		expectedReportMessages []string
	}{
		{
			desc:                   "has readme",
			repoPath:               "testdata/readme",
			folder:                 "hasReadme",
			expectedReportMessages: nil,
		},
		{
			desc:                   "no readme",
			repoPath:               "testdata/readme",
			folder:                 "noReadme",
			expectedReportMessages: []string{"testdata/readme/noReadme/README.md does not exist!"},
		},
		{
			desc:                   "no readme",
			repoPath:               "testdata/readme",
			folder:                 "emptyReadme",
			expectedReportMessages: []string{"testdata/readme/emptyReadme/README.md exists, but has no content!"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			memReporter := &MemReporter{}
			checker := Checker{
				repoPath: tc.repoPath,
				reporter: memReporter,
			}
			checker.checkReadme(tc.folder)

			assert.Equal(t, tc.expectedReportMessages, memReporter.message)
		})
	}
}

func TestCheckPackageDocSubfolders(t *testing.T) {
	tcs := []struct {
		desc                   string
		repoPath               string
		folder                 string
		expectedReportMessages []string
	}{
		{
			desc:     "has package doc",
			repoPath: "testdata/",
			folder:   "gopackagedoc",
			expectedReportMessages: []string{
				"testdata/gopackagedoc/noGoPackageDoc/doc.go does not exist!",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			memReporter := &MemReporter{}
			checker := Checker{
				repoPath: tc.repoPath,
				reporter: memReporter,
			}
			assert.Nil(t, checker.checkPackageDocSubfolders(tc.folder))

			sort.Strings(tc.expectedReportMessages)
			sort.Strings(memReporter.message)

			assert.Equal(t, tc.expectedReportMessages, memReporter.message)
		})
	}
}

func TestCheckGoFileDoc(t *testing.T) {
	memReporter := &MemReporter{}
	checker := Checker{
		repoPath: "testdata/filedoc",
		reporter: memReporter,
	}
	assert.Nil(t, checker.checkGoFileDoc("filedoc1"))

	expectedMessages := []string{
		"testdata/filedoc/filedoc1/top_without.go does not contain a file comment!",
		"testdata/filedoc/filedoc1/nested/nested_without.go does not contain a file comment!",
		"testdata/filedoc/filedoc1/nested/double/double_without.go does not contain a file comment!",
	}
	sort.Strings(expectedMessages)
	sort.Strings(memReporter.message)

	assert.Equal(t, expectedMessages, memReporter.message)
}

func TestCheckMonoRepo(t *testing.T) {
	memReporter := &MemReporter{}
	checker := Checker{
		repoPath: "testdata/monorepo",
		reporter: memReporter,
	}

	assert.Nil(t, runCheck(modeMonoRepo, &checker))

	expectedMessages := []string{
		"testdata/monorepo/README.md does not exist!",
		"testdata/monorepo/apps/app1/doc.go does not exist!",
		"testdata/monorepo/pkg/pkg1/doc.go does not exist!",
		"testdata/monorepo/pkg/pkg1/nodoc.go does not contain a file comment!",
		"testdata/monorepo/services/service1/doc.go does not exist!",
		"testdata/monorepo/services/service1/nodoc.go does not contain a file comment!",
	}
	sort.Strings(expectedMessages)
	sort.Strings(memReporter.message)

	assert.Equal(t, expectedMessages, memReporter.message)
}

func TestCheckApp(t *testing.T) {
	memReporter := &MemReporter{}
	checker := Checker{
		repoPath: "testdata/app",
		reporter: memReporter,
	}

	assert.Nil(t, runCheck(modeApp, &checker))

	expectedMessages := []string{
		"testdata/app/README.md does not exist!",
		"testdata/app/app/foo/doc.go does not exist!",
		"testdata/app/app/foo/foo_nodoc.go does not contain a file comment!",
		"testdata/app/app/nodoc.go does not contain a file comment!",
	}
	sort.Strings(expectedMessages)
	sort.Strings(memReporter.message)

	assert.Equal(t, expectedMessages, memReporter.message)
}

// helper

type MemReporter struct {
	message []string
}

func (r *MemReporter) Report(message string) {
	r.message = append(r.message, message)
}

func (r *MemReporter) FoundIssues() bool {
	return len(r.message) != 0
}
