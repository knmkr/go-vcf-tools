package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s ID,ID,...(commna delimited sample IDs to be kept) < in.vcf > out.vcf\n", os.Args[0])
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Parse header lines
	var sample_ids []string
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "##") {
			fmt.Println(line)
		} else if strings.HasPrefix(line, "#CHROM") {
			fields := strings.Split(line, "\t")
			fmt.Print(strings.Join(fields[0:9], "\t"))
			fmt.Print("\t")
			sample_ids = fields[9:]
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Get indices of sample IDs to be kept
	keep_ids := strings.Split(os.Args[1], ",")
	var keep_idxs []int
	for i := range keep_ids {
		for j := range sample_ids {
			if keep_ids[i] == sample_ids[j] {
				keep_idxs = append(keep_idxs, j)
			}
		}
	}

	fmt.Println(strings.Join(subset(sample_ids, keep_idxs), "\t"))

	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")

		result := []string{}
		result = append(result, records[0:9]...)
		result = append(result, subset(records[9:], keep_idxs)...)
		fmt.Println(strings.Join(result, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func subset(records []string, keep_idxs []int) []string {
	result := []string{}
	for i := range keep_idxs {
		result = append(result, records[keep_idxs[i]])
	}
	return result
}
