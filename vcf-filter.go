package main

import (
	"flag"
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
	arg_keep_ids := flag.String("keep-ids", "", "Path to a file of rs IDs to be kept. Each line contains one rs ID. E.g. rs123")
	arg_keep_pos := flag.String("keep-pos", "", "Path to a file of loci to be kept. Each line contains one TAB delimited loci (chromosome and position). E.g. 1[TAB]100")
	flag.Parse()

	if len(os.Args) != 3 && len(os.Args) != 5 {
		flag.Usage()
		os.Exit(0)
	}

	pattern := regexp.MustCompile(`rs(\d+)`)

	keep_ids := make(map[int]bool)
	keep_pos := make(map[int64]bool)

	// Get SNP IDs to be kept if exists
	if *arg_keep_ids != "" {
		var ids_fp *os.File
		var err error

		ids_fp, err = os.Open(*arg_keep_ids)
		if err != nil {
			panic(err)
		}
		defer ids_fp.Close()

		ids_reader := bufio.NewReaderSize(ids_fp, 128 * 1024)
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
	}

	// Get loci to be kept if exists
	if *arg_keep_pos != "" {
		var pos_fp *os.File
		var err error

		pos_fp, err = os.Open(*arg_keep_pos)
		if err != nil {
			panic(err)
		}
		defer pos_fp.Close()

		pos_reader := bufio.NewReaderSize(pos_fp, 128 * 1024)
		pos_line, err := lib.Readln(pos_reader)
		for err == nil {
			records := strings.Split(pos_line, "\t")
			chrom  := records[0]
			pos, _ := strconv.ParseInt(records[1], 10, 64)
			chrpos := lib.ChrPos(chrom, pos)
			keep_pos[chrpos] = true

			pos_line, err = lib.Readln(pos_reader)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
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

	// Parse body lines
	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")
		var is_pass bool

		// Filter by id
		if *arg_keep_ids != "" {
			id_found := pattern.FindStringSubmatch(records[2])
			if id_found  != nil {
				id, _  := strconv.Atoi(id_found[1])

				if keep_ids[id] {
					is_pass = true
				}
			}
		}

		// Filter by loci
		if *arg_keep_pos != "" {
			chrom  := records[0]
			pos, _ := strconv.ParseInt(records[1], 10, 64)
			chrpos := lib.ChrPos(chrom, pos)

			if keep_pos[chrpos] {
				is_pass = true
			}
		}

		if is_pass {
			fmt.Println(line)
		}

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}
