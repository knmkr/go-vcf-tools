package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandFilter,
}

var commandFilter = cli.Command{
	Name:   "filter",
	Usage:  "Filter out VCF records by conditions",
	Action: doFilter,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "keep-ids",
			Usage: "Path to a file of rs IDs to be kept. Each line contains one rs ID. E.g. rs123",
		},
		cli.StringFlag{
			Name:  "keep-pos",
			Usage: "Path to a file of loci to be kept. Each line contains one TAB delimited loci (chromosome and position). E.g. 1[TAB]100",
		},
		cli.BoolFlag{
			Name:  "keep-only-pass",
			Usage: "Keep only FILTER == PASS records",
		},
	},
}
