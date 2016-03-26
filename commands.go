package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandFilter,
	commandSubset,
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

var commandSubset = cli.Command{
	Name:   "subset",
	Usage:  "",
	Action: doSubset,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "keep-id",
			Usage: "A sample ID to be kept. E.g., NA00001",
		},
		cli.StringFlag{
			Name:  "keep-ids",
			Usage: "Path to a file of sample IDs to be kept. Each line contains one sample ID.",
		},
		cli.StringFlag{
			Name:  "keep-index",
			Usage: "An index of sample ID field to be kept. E.g., to keep 1st sample, set: 0",
		},
	},
}
