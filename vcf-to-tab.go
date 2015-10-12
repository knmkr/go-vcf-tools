package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
	"regexp"
	"flag"
)

func main() {
	is_without_header := flag.Bool("without-header", false, "Output without header line.")
	is_without_chr_pos := flag.Bool("without-chr-pos", false, "Output without CHROM and POS.")
	is_rs_id_as_int := flag.Bool("rs-id-as-int", false, "Output rs ID as integer.")
	is_genotype_as_pg_array := flag.Bool("genotype-as-pg-array", false, "Output genotype as PostgreSQL array. E.g., '{G,G}'")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "##") {
			continue
		} else if strings.HasPrefix(line, "#CHROM") {
			if ! *is_without_header {
				fields := strings.Split(line, "\t")
				if ! *is_without_chr_pos {
					fmt.Print("#CHROM\tPOS\tID\t")
				} else {
					fmt.Print("ID\t")
				}
				fmt.Println(strings.Join(fields[9:], "\t"))
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pattern := regexp.MustCompile(`rs(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")

		chrom := records[0]
		pos := records[1]
		id := records[2]

		if *is_rs_id_as_int {
			id_found := pattern.FindStringSubmatch(records[2])
			if id_found  != nil {
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
					if *is_genotype_as_pg_array {
						genotype = "{" + strings.Join(_gt, ",") + "}"
					} else {
						genotype = strings.Join(_gt, "/")
					}
					genotypes = append(genotypes, genotype)
				}
			}
		}

		result := []string{}
		if ! *is_without_chr_pos {
			result = []string{chrom, pos, id}
		} else {
			result = []string{id}
		}
		result = append(result, genotypes...)
		fmt.Println(strings.Join(result, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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
