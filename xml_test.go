package main

import (
	//"log"
	"testing"
)

func TestTagsContainHyphens(t *testing.T) {
	err := extractor([]string{tagsContainHyphens})
	if err != nil {
		t.Error(err)
	}
}

func TestTagsWithSameNameDifferentNameSpaceXML(t *testing.T) {
	err := extractor([]string{sameNameDifferentNameSpaceXML})
	if err != nil {
		t.Error(err)
	}
}

func TestMixedCaseSameNameTagsXML_Issue24(t *testing.T) {
	err := extractor([]string{mixedCaseSameNameXML_Issue24})
	if err != nil {
		t.Error(err)
	}
}

//https://github.com/gnewton/chidley/issues/14
func TestGithubIssue14(t *testing.T) {
	err := extractor([]string{githubIssue14})
	if err != nil {
		t.Error(err)
	}
}
