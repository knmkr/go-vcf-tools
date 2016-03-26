package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandFilter,
	commandSubset,
	commandFreq,
	commandUpdate,
	commandToTab,
	commandFillRsids,
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

var commandFreq = cli.Command{
	Name:   "freq",
	Usage:  "",
	Action: doFreq,
}

var commandUpdate = cli.Command{
	Name:   "update",
	Usage:  "",
	Action: doUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "rs-merge-arch",
			Usage: "",
		},
	},
}

var commandToTab = cli.Command{
	Name:   "to-tab",
	Usage:  "",
	Action: doToTab,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "without-header",
			Usage: "Output without header line.",
		},
		cli.BoolFlag{
			Name:  "without-chr-pos",
			Usage: "Output without CHROM and POS.",
		},
		cli.BoolFlag{
			Name:  "rs-id-as-int",
			Usage: "Output rs ID as integer.",
		},
		cli.BoolFlag{
			Name:  "genotype-as-pg-array",
			Usage: "Output genotype as PostgreSQL array. E.g., '{G,G}'",
		},
		cli.BoolFlag{
			Name:  "chrx-genotype-as-homo",
			Usage: "Output chrX genotype as homozygous",
		},
	},
}

var commandFillRsids = cli.Command{
	Name:   "fill-rsids",
	Usage:  "",
	Action: doFillRsids,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bucket",
			Usage: "Mappings of chrom/pos on reference genome. E.g., b142_SNPChrPosOnRef_105.bcp.gz",
		},
		cli.BoolFlag{
			Name:  "setup",
			Usage: "Setup local db.",
		},
		cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite rs ids if already exist in vcf. However, for loci not in local db, original records will be kept.",
		},
		cli.BoolFlag{
			Name:  "strict",
			Usage: "Along with '-overwrite' option, for loci not in local db will be filled as '.'",
		},
	},
}
