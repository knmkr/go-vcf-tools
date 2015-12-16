package main

import (
	"flag"
	"os"
	"io"
	"fmt"
	"bufio"
	"regexp"
	"errors"
	"strings"
	"github.com/knmkr/go-vcf-tools/lib"
)

func main() {
	arg_remove_chr_string := flag.Bool("remove-chr-string", false, "Remove 'chr' strings from vcf CHROM records and output only chromosome codes. E.g. 'chr1' will be outputed as '1'.")
	flag.Parse()

	// Parse header lines
	reader := bufio.NewReaderSize(os.Stdin, 128 * 1024)

	contig_fields_pattern := regexp.MustCompile(`##contig=<(.+)>`)

	line, err := lib.Readln(reader)
	for err == nil {
		if strings.HasPrefix(line, "##") {

			// Remove 'chr' from contig fields in header
			if *arg_remove_chr_string {
				contig_field_founds := contig_fields_pattern.FindStringSubmatch(line)
				if contig_field_founds != nil {
					contig_field := contig_field_founds[1]
					result := []string{}
					for _, x  := range strings.Split(contig_field, ",") {
						if strings.HasPrefix(x, "ID") {
							result = append(result, strings.Replace(x, "chr", "", 1))
						} else {
							result = append(result, x)
						}
					}
					fmt.Println("##contig=<" + strings.Join(result, ",") + ">")
				} else {
					fmt.Println(line)
				}
			} else {
				fmt.Println(line)
			}

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

		var chrom string
		if *arg_remove_chr_string {
			chrom = strings.Replace(records[0], "chr", "", 1)
		} else {
			chrom = records[0]
		}

		result := []string{}
		result = append(result, chrom)
		result = append(result, records[1:]...)
		fmt.Println(strings.Join(result, "\t"))

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}
