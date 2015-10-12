package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"regexp"
	"strconv"
	"compress/gzip"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s <RsMergeArch.bcp.gz> < in.vcf > out.vcf\n", os.Args[0])
		os.Exit(0)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	// [dbSNP Column Description for table: RsMergeArc](http://www.ncbi.nlm.nih.gov/projects/SNP/snp_db_table_description.cgi?t=RsMergeArch)
	//
	// - Table name and description
	//
	// | Table Description                                                                                                             |
	// |-------------------------------------------------------------------------------------------------------------------------------|
	// | "refSNP(rs) cluster is based on unique genome position. On new genome assembly, previously different contig may               |
	// | align. So different rs clusters map to the same location. In this case, we merge the rs. This table tracks this merging."     |
	//
	// - Table column and description
	//
	// | Column            | Description                                                                | Type          | Byte | Order |
	// |-------------------+----------------------------------------------------------------------------+---------------+------+-------|
	// | rsHigh            | Since rs# is assigned sequentially. Low number means the rs occurs         | int           |    4 |     1 |
	// |                   | early. So we always merge high rs number into low rs number.               |               |      |       |
	// | rsLow             |                                                                            | int           |    4 |     2 |
	// | build_id          | dbSNP build id when this rsHigh was merged into rsLow.                     | smallint      |    2 |     3 |
	// | orien             | The orientation between rsHigh and rsLow.                                  | tinyint       |    1 |     4 |
	// | create_time       |                                                                            | smalldatetime |    4 |     5 |
	// | last_updated_time |                                                                            | smalldatetime |    4 |     6 |
	// | rsCurrent         | rsCurrent is the current rs for rsHigh. If rs9 is merged into rs5 which is | int           |    4 |     7 |
	// |                   | later merged into rs2, then rsCurrent is 2 for rsHigh=9.                   |               |      |       |
	// | orien2Current     |                                                                            | tinyint       |    1 |     8 |
	//
	// This table/column description is last updated at: Mar 18 2015 02:51:00:000PM.

	// Get merge mappings of rs IDs
	rsHigh2current := make(map[int]int)

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")
		rsHigh, _  := strconv.Atoi(records[0])
		rsCurrent, _  := strconv.Atoi(records[6])
		rsHigh2current[rsHigh] = rsCurrent
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Parse header lines
	scanner = bufio.NewScanner(os.Stdin)

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

	pattern := regexp.MustCompile(`rs(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		records := strings.Split(line, "\t")

		// Update rs ID
		var id_updated_str string
		id_found := pattern.FindStringSubmatch(records[2])
		if id_found  != nil {
			id, _  := strconv.Atoi(id_found[1])
			id_updated := rsHigh2current[id]

			if id_updated != 0 {
				id_updated_str = "rs" + strconv.Itoa(id_updated)  // Map to current ID
			} else {
				id_updated_str = records[2]  // ID is not listed in merge history
			}
		} else {
			id_updated_str = records[2]  // ID is not rs ID
		}

		result := []string{}
		result = append(result, records[0:2]...)
		result = append(result, id_updated_str)
		result = append(result, records[3:]...)
		fmt.Println(strings.Join(result, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
