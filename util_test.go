package main

import (
	//"log"
	"testing"
	//	"unicode"
)

func TestExtractExcludedTags_EmptyString(t *testing.T) {
	s := ""
	b := containsUnicodeSpace(s)
	if b != false {
		t.Error()
	}
}

func TestExtractExcludedTags_CorrectString(t *testing.T) {
	s := "a,b,c,d,e"
	b := containsUnicodeSpace(s)
	if b != false {
		t.Error()
	}
}

func TestExtractExcludedTags_StringWithSpaces(t *testing.T) {
	s := "a,b,c,d, e"
	b := containsUnicodeSpace(s)
	if b != true {
		t.Error()
	}
}

// Latin-1 spaces taken from https://golang.org/pkg/unicode/#IsSpace
func TestExtractExcludedTags_StringWithAllLatin1Spaces(t *testing.T) {
	s := "\t\n\v\f\r U+0085 U+00A0"
	b := containsUnicodeSpace(s)
	if b != true {
		t.Error()
	}
}
