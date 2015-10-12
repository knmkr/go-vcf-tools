package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"flag"
	"strconv"
)

func main() {
	arg_keep_id := flag.String("keep-id", "", "A sample ID to be kept. E.g., NA00001")
	arg_keep_ids := flag.String("keep-ids", "", "Path to a file of sample IDs to be kept. Each line contains one sample ID.")
	arg_keep_index := flag.String("keep-index", "", "An index of sample ID field to be kept. E.g., to keep 1st sample, set: 0")
	flag.Parse()

	if len(os.Args) <=2 || len(os.Args) > 4 {
		fmt.Fprintln(os.Stderr, "Set only one of --keep-id/--keep-ids/--keep-index")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
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
	keep_ids := []string{}
	keep_idxs := []int{}

	if *arg_keep_id != "" || *arg_keep_ids != "" {
		if *arg_keep_id != "" {
			// A sample ID to be kept. E.g., NA00001
			keep_ids = append(keep_ids, *arg_keep_id)
		} else {
			// Path to a file of sample IDs to be kept. Each line contains one sample ID.
			f, err := os.Open(*arg_keep_ids)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				keep_ids = append(keep_ids, line)
			}
		}

		for i := range keep_ids {
			for j := range sample_ids {
				if keep_ids[i] == sample_ids[j] {
					keep_idxs = append(keep_idxs, j)
					break
				}
			}
		}
	} else if *arg_keep_index != "" {
		// An index of sample ID field to be kept. E.g., to keep 1st sample, set: 0
		_keep_idx, _  := strconv.Atoi(*arg_keep_index)
		keep_idxs = append(keep_idxs, _keep_idx)
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
