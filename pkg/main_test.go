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
			repoPath:               "../testdata/readme",
			folder:                 "hasReadme",
			expectedReportMessages: nil,
		},
		{
			desc:                   "no readme",
			repoPath:               "../testdata/readme",
			folder:                 "noReadme",
			expectedReportMessages: []string{"../testdata/readme/noReadme/README.md does not exist!"},
		},
		{
			desc:                   "no readme",
			repoPath:               "../testdata/readme",
			folder:                 "emptyReadme",
			expectedReportMessages: []string{"../testdata/readme/emptyReadme/README.md exists, but has no content!"},
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

func TestCheckReadmeSubfolders(t *testing.T) {
	tcs := []struct {
		desc                   string
		repoPath               string
		folder                 string
		expectedReportMessages []string
	}{
		{
			desc:     "has readme",
			repoPath: "../testdata/",
			folder:   "readme",
			expectedReportMessages: []string{
				"../testdata/readme/noReadme/README.md does not exist!",
				"../testdata/readme/emptyReadme/README.md exists, but has no content!",
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
			checker.checkReadmeSubfolders(tc.folder)

			sort.Strings(tc.expectedReportMessages)
			sort.Strings(memReporter.message)

			assert.Equal(t, tc.expectedReportMessages, memReporter.message)
		})
	}
}

// helper

type MemReporter struct {
	message []string
}

func (r *MemReporter) Report(message string) {
	r.message = append(r.message, message)
}
