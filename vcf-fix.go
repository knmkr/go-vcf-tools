package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/knmkr/go-vcf-tools/lib"
	"io"
	"os"
	"regexp"
	"strings"
)

func doFix(c *cli.Context) {
	arg_remove_chr_string := c.Bool("remove-chr-string")
	arg_remove_qual := c.Bool("remove-qual")
	arg_remove_filter := c.Bool("remove-filter")
	arg_remove_info := c.Bool("remove-info")
	arg_keep_gt_only := c.Bool("keep-only-gt")

	// Parse header lines
	reader := bufio.NewReaderSize(os.Stdin, 128*1024)

	contig_pattern := regexp.MustCompile(`##contig=<(.+)>`)
	info_pattern := regexp.MustCompile(`##INFO=<(.+)>`)
	format_pattern := regexp.MustCompile(`##FORMAT=<(.+)>`)
	filter_pattern := regexp.MustCompile(`##FILTER=<(.+)>`)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {
			contig_founds := contig_pattern.FindStringSubmatch(line)
			info_founds := info_pattern.FindStringSubmatch(line)
			format_founds := format_pattern.FindStringSubmatch(line)
			filter_founds := filter_pattern.FindStringSubmatch(line)

			if arg_remove_chr_string && contig_founds != nil {
				// Remove 'chr' from contig meta-infos in header
				result := []string{}
				for _, x := range strings.Split(contig_founds[1], ",") {
					if strings.HasPrefix(x, "ID") {
						result = append(result, strings.Replace(x, "chr", "", 1))
					} else {
						result = append(result, x)
					}
				}
				fmt.Println("##contig=<" + strings.Join(result, ",") + ">")
			} else if arg_remove_info && info_founds != nil {
				// Skip INFO meta-info
			} else if arg_remove_filter && filter_founds != nil {
				// Skip FILTER meta-info
			} else if arg_keep_gt_only && format_founds != nil {
				// Skip FORMAT meta-info tags except GT
				for _, x := range strings.Split(format_founds[1], ",") {
					if x == "ID=GT" {
						fmt.Println(line)
						continue
					}
				}
			} else {
				fmt.Println(line)
			}
		} else if strings.HasPrefix(line, "#CHROM") {
			fmt.Println(line)
			break
		} else {
			err = errors.New("Invalid VCF header")
			break
		}

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Parse body lines

	// > 1.4.1 Fixed fields
	// > There are 8 fixed fields per record. All data lines are tab-delimited.
	// > In all cases, missing values are specified with a dot ('.').
	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		var chrom string
		if arg_remove_chr_string {
			chrom = strings.Replace(records[0], "chr", "", 1)
		} else {
			chrom = records[0]
		}

		// > 6. QUAL - quality: Phred-scaled quality score for the assertion made in ALT.
		// > ... If unknown, the missing value should be specified. (Numeric)
		var qual string
		if arg_remove_qual {
			qual = "."
		} else {
			qual = records[5]
		}

		// > 7. FILTER - filter status: PASS if this position has passed all filters, i.e. a call is made at this position.
		// > ... If filters have not been applied, then this field should be set to the missing value.
		// > (String, no white-space or semi-colons permitted)
		var filter string
		if arg_remove_filter {
			filter = "."
		} else {
			filter = records[6]
		}

		var info string
		if arg_remove_info {
			info = "."
		} else {
			info = records[7]
		}

		var format string
		genotypes := []string{}
		if arg_keep_gt_only {
			// > 1.4.2 Genotype fields
			// > ... The first sub-field must always be the genotype (GT) if it is present.
			format = "GT"
			for _, genotype := range records[9:] {
				genotypes = append(genotypes, strings.Split(genotype, ":")[0])
			}
		} else {
			format = records[8]
			genotypes = records[9:]
		}

		result := []string{}
		result = append(result, chrom)
		result = append(result, records[1:5]...)
		result = append(result, qual)
		result = append(result, filter)
		result = append(result, info)
		result = append(result, format)
		result = append(result, genotypes...)
		fmt.Println(strings.Join(result, "\t"))

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}
