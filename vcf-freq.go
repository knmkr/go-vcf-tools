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
			fmt.Println("#CHROM\tPOS\tID\tAllele\tFreq")
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pattern := regexp.MustCompile(`[|/]`)

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

		alleles := []string{}
		alleles = append(alleles, ref)
		alleles = append(alleles, alt...)

		var count []int
		for i := 0; i < len(alleles); i++ {
			count = append(count, 0)
		}

		for i := range gts {
			gt := strings.Split(gts[i], ":")

			for j := range gt {
				if format[j] == "GT" {
					gt_idxs := pattern.Split(gt[j], -1)

					for i := range gt_idxs {
						gt_idx, _ := strconv.Atoi(gt_idxs[i])
						count[gt_idx] += 1
					}
				}
			}
		}

		total := float64(sum(count))  // TODO: decimal?
		freqs := []string{}
		for i := range count {
			freqs = append(freqs, fmt.Sprintf("%.4f", float64(count[i]) / total))  // TODO:
		}

		result := []string{chrom, pos, id, strings.Join(alleles, ","), strings.Join(freqs, ",")}
		fmt.Println(strings.Join(result, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func sum(vals []int) int {
	var result int
	for i:= range vals {
		result += vals[i]
	}
	return result
}
