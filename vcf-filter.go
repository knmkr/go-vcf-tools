package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"regexp"
	"strconv"
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

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		id_found := pattern.FindStringSubmatch(line)
		if id_found != nil {
			keep_id, _ := strconv.Atoi(id_found[1])
			keep_ids[keep_id] = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	scanner = bufio.NewScanner(os.Stdin)

	// Parse header lines
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "##") {
			fmt.Println(line)
		} else if strings.HasPrefix(line, "#CHROM") {
			fmt.Println(line)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")

		// Filter by id
		id_found := pattern.FindStringSubmatch(records[2])
		if id_found  != nil {
			id, _  := strconv.Atoi(id_found[1])
			if keep_ids[id] {
				fmt.Println(line)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
