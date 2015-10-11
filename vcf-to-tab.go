package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"strconv"
	"regexp"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "##") {
			continue
		} else if strings.HasPrefix(line, "#CHROM") {
			fields := strings.Split(line, "\t")
			fmt.Print("#CHROM\tPOS\tID\tREF\t")
			fmt.Println(strings.Join(fields[9:], "\t"))
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")

		chrom := records[0]
		pos := records[1]
		id := records[2]
		ref := records[3]
		alt := strings.Split(records[4], ",")
		format := strings.Split(records[8], ":")
		gts := records[9:]

		genotypes := []string{}

		for i := range gts {
			var genotype string

			gt := strings.Split(gts[i], ":")

			for j := range gt {
				if format[j] == "GT" {
					genotype = strings.Join(gt2genotype(ref, alt, gt[j]), "/")
				}
			}
			genotypes = append(genotypes, genotype)
		}

		result := []string{chrom, pos, id, ref}
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
