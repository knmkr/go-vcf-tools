package main

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"regexp"
	"errors"
	"strings"
	"strconv"
	"github.com/knmkr/go-vcf-tools/lib"
)

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 128 * 1024)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {
			// pass
		} else if strings.HasPrefix(line, "#CHROM") {
			fmt.Println("#CHROM\tPOS\tID\tAllele\tFreq")
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

	pattern := regexp.MustCompile(`[|/]`)

	line, err = lib.Readln(reader)
	for err == nil {
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

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func sum(vals []int) int {
	var result int
	for i:= range vals {
		result += vals[i]
	}
	return result
}
