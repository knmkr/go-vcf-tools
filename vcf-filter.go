package main

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"regexp"
	"errors"
	"strconv"
	"strings"
	"github.com/knmkr/go-vcf-tools/lib"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s <a list of SNPs to be kept> < in.vcf > out.vcf\n", os.Args[0])
		os.Exit(0)
	}

	// Get SNP IDs to be kept
	var fp *os.File
	var err error

	fp, err = os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	pattern := regexp.MustCompile(`rs(\d+)`)

	keep_ids := make(map[int]bool)

	ids_reader := bufio.NewReaderSize(fp, 128 * 1024)
	ids_line, err := lib.Readln(ids_reader)
	for err == nil {
		id_found := pattern.FindStringSubmatch(ids_line)
		if id_found != nil {
			keep_id, _ := strconv.Atoi(id_found[1])
			keep_ids[keep_id] = true
		}

		ids_line, err = lib.Readln(ids_reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Parse header lines
	reader := bufio.NewReaderSize(os.Stdin, 128 * 1024)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {
			fmt.Println(line)
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


	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		// Filter by id
		id_found := pattern.FindStringSubmatch(records[2])
		if id_found  != nil {
			id, _  := strconv.Atoi(id_found[1])
			if keep_ids[id] {
				fmt.Println(line)
			}
		}

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}
