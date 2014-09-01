package main

type FQN struct {
	space string
	name  string
}

type FQNAbbr struct {
	FQN
	abbr string
}
