package main

import (
	"testing"
)

const XML = `
<?xml version=\1.0" encoding="windows-1251"?>
<data>
  <order>43</order>
</data>
`

func TestHasNonDefaultEncoding(t *testing.T) {

}

const XMLHasTagCalledRoot = `
<data>
  <root>
  </root>
</data>
`

func TestHasTagCalledRoot(t *testing.T) {

}

const XMLHasCapitalizationTagCollisions = `
<data>
  <Entry>
    <value>123</value>
  </Entry>
  <entry>
    <op>add</op>
  </entry>
</data>
`

func TestHasCapitalizationTagCollisions(t *testing.T) {

}
