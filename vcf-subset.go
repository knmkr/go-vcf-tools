package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"flag"
	"bufio"
	"errors"
	"strings"
	"strconv"
	"github.com/knmkr/go-vcf-tools/lib"
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

	reader := bufio.NewReaderSize(os.Stdin, 64 * 1024)

	// Parse header lines
	var sample_ids []string

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {
			fmt.Println(line)
		} else if strings.HasPrefix(line, "#CHROM") {
			fields := strings.Split(line, "\t")
			fmt.Print(strings.Join(fields[0:9], "\t"))
			fmt.Print("\t")
			sample_ids = fields[9:]
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

	// Get indices of sample IDs to be kept
	keep_ids := []string{}
	keep_idxs := []int{}

	if *arg_keep_id != "" || *arg_keep_ids != "" {
		if *arg_keep_id != "" {
			// A sample ID to be kept. E.g., NA00001
			keep_ids = append(keep_ids, *arg_keep_id)
		} else {
			// Path to a file of sample IDs to be kept. Each line contains one sample ID.
			fp, err := os.Open(*arg_keep_ids)
			if err != nil {
				panic(err)
			}
			defer fp.Close()

			ids_reader := bufio.NewReaderSize(fp, 128 * 1024)
			ids_line, err := lib.Readln(ids_reader)
			for err == nil {
				keep_ids = append(keep_ids, ids_line)
			 	ids_line, err = lib.Readln(ids_reader)
			}
			if err != nil && err != io.EOF {
				panic(err)
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

		if len(keep_idxs) == 0 {
			fmt.Println()
			log.Fatal("No sample IDs matched.")
		}

	} else if *arg_keep_index != "" {
		// An index of sample ID field to be kept. E.g., to keep 1st sample, set: 0
		_keep_idx, _  := strconv.Atoi(*arg_keep_index)

		if _keep_idx > len(sample_ids) {
			fmt.Println()
			log.Fatal("No sample IDs matched.")
		}
		keep_idxs = append(keep_idxs, _keep_idx)
	}

	fmt.Println(strings.Join(subset(sample_ids, keep_idxs), "\t"))

	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		result := []string{}
		result = append(result, records[0:9]...)
		result = append(result, subset(records[9:], keep_idxs)...)
		fmt.Println(strings.Join(result, "\t"))

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func subset(records []string, keep_idxs []int) []string {
	result := []string{}
	for i := range keep_idxs {
		result = append(result, records[keep_idxs[i]])
	}
	return result
}
