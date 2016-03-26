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
	"strconv"
	"strings"
)

func doToTab(c *cli.Context) {
	is_without_header := c.Bool("without-header")
	is_without_chr_pos := c.Bool("without-chr-pos")
	is_rs_id_as_int := c.Bool("rs-id-as-int")
	is_genotype_as_pg_array := c.Bool("genotype-as-pg-array")
	is_chrx_genotype_as_homo := c.Bool("chrx-genotype-as-homo")

	reader := bufio.NewReaderSize(os.Stdin, 128*1024)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {
			// pass
		} else if strings.HasPrefix(line, "#CHROM") {
			if !is_without_header {
				fields := strings.Split(line, "\t")
				if !is_without_chr_pos {
					fmt.Print("#CHROM\tPOS\tID\t")
				} else {
					fmt.Print("ID\t")
				}
				fmt.Println(strings.Join(fields[9:], "\t"))
			}
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

	pattern := regexp.MustCompile(`rs(\d+)`)

	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		chrom := records[0]
		pos := records[1]
		id := records[2]

		if is_rs_id_as_int {
			id_found := pattern.FindStringSubmatch(records[2])
			if id_found != nil {
				id = id_found[1]
			}
		}

		ref := records[3]
		alt := strings.Split(records[4], ",")
		format := strings.Split(records[8], ":")
		gts := records[9:]

		genotypes := []string{}

		for i := range gts {
			gt := strings.Split(gts[i], ":")

			for j := range gt {
				var genotype string
				if format[j] == "GT" {
					_gt := gt2genotype(ref, alt, gt[j])

					if is_chrx_genotype_as_homo && chrom == "X" {
						if len(_gt) == 1 {
							_gt = append(_gt, _gt...)
						}
					}

					if is_genotype_as_pg_array {
						genotype = "{" + strings.Join(_gt, ",") + "}"
					} else {
						genotype = strings.Join(_gt, "/")
					}
					genotypes = append(genotypes, genotype)
				}
			}
		}

		result := []string{}
		if !is_without_chr_pos {
			result = []string{chrom, pos, id}
		} else {
			result = []string{id}
		}
		result = append(result, genotypes...)
		fmt.Println(strings.Join(result, "\t"))

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func gt2genotype(ref string, alt []string, gt string) []string {
	pattern := regexp.MustCompile(`[|/]`)

	alleles := []string{}
	alleles = append(alleles, ref)
	alleles = append(alleles, alt...)

	gt_idxs := pattern.Split(gt, -1)

	genotype := []string{}
	for i := range gt_idxs {
		gt_idx, _ := strconv.Atoi(gt_idxs[i])
		genotype = append(genotype, alleles[gt_idx])
	}

	return genotype
}
