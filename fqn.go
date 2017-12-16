package main

type FQN struct {
	space     string
	name      string
	maxLength int
}

type FQNAbbr struct {
	FQN
	abbr string
}
