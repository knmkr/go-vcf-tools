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
	arg_remove_info := c.Bool("remove-info")
	arg_keep_gt_only := c.Bool("keep-only-gt")

	// Parse header lines
	reader := bufio.NewReaderSize(os.Stdin, 128*1024)

	contig_fields_pattern := regexp.MustCompile(`##contig=<(.+)>`)
	info_fields_pattern := regexp.MustCompile(`##INFO=<(.+)>`)
	format_fields_pattern := regexp.MustCompile(`##FORMAT=<(.+)>`)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {

			if arg_remove_chr_string {
				// Remove 'chr' from contig fields in header
				contig_field_founds := contig_fields_pattern.FindStringSubmatch(line)
				if contig_field_founds != nil {
					contig_field := contig_field_founds[1]
					result := []string{}
					for _, x := range strings.Split(contig_field, ",") {
						if strings.HasPrefix(x, "ID") {
							result = append(result, strings.Replace(x, "chr", "", 1))
						} else {
							result = append(result, x)
						}
					}
					fmt.Println("##contig=<" + strings.Join(result, ",") + ">")
				} else {
					fmt.Println(line)
				}
			} else if arg_remove_info {
				// Skip INFO fields
				info_field_founds := info_fields_pattern.FindStringSubmatch(line)
				if info_field_founds == nil {
					fmt.Println(line)
				}
			} else if arg_keep_gt_only {
				// Skip FORMAT field tags except GT
				format_field_founds := format_fields_pattern.FindStringSubmatch(line)
				if format_field_founds != nil {
					format_field := format_field_founds[1]
					for _, x := range strings.Split(format_field, ",") {
						if x == "ID=GT" {
							fmt.Println(line)
							continue
						}
					}
				} else {
					fmt.Println(line)
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
	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		var chrom string
		if arg_remove_chr_string {
			chrom = strings.Replace(records[0], "chr", "", 1)
		} else {
			chrom = records[0]
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
		result = append(result, records[1:7]...)
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
