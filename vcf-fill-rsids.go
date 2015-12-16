package main

import (
	"regexp"
	"errors"
	"flag"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/knmkr/go-vcf-tools/lib"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	arg_bucket    := flag.String("bucket", "", "Mappings of chrom/pos on reference genome. E.g., b142_SNPChrPosOnRef_105.bcp.gz")
	arg_setup     := flag.Bool("setup", false, "Setup local db.")
	arg_overwrite := flag.Bool("overwrite", false, "Overwrite rs ids if already exist in vcf. However, for loci not in local db, original records will be kept.")
	arg_strict    := flag.Bool("strict", false, "Along with '-overwrite' option, for loci not in local db will be filled as '.'")
	flag.Parse()

	if len(os.Args) <=2 {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "[FATAL] too few arguments")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(0)
	} else if ! *arg_overwrite && *arg_strict {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "[FATAL] -strict option is only effective with -overwrite option")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(0)
	} else if *arg_bucket == "" {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "[FATAL] -bucket is required")
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(0)
	}

	databaseName := "bolt.db"
	bucketName := []byte(path.Base(*arg_bucket))

	// Store chrpos <=> rsid mappings into bolt.db
	db, err := bolt.Open(databaseName, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if *arg_setup {
		f, err := os.Open(*arg_bucket)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		gz, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer gz.Close()

		// TODO: workaround for non-uniq chrpos. skip high rs numbers?
		err = db.Batch(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}

			map_reader := bufio.NewReaderSize(gz, 128*1024)
			map_line, err := lib.Readln(map_reader)
			for err == nil {
				// chrom/pos to rs id mapping resource file ([TAB] delimited)
				//
				// | rs id  | chrom  | pos    |
				// |--------|--------|--------|
				// |  xxxxx |     xx |  xxxxx |
				//
				records := strings.Split(map_line, "\t")
				rsId := strings.Replace(records[0], "rs", "", 1)
				rsChr := records[1]
				rsPos, _ := strconv.ParseInt(records[2], 10, 64)

				if rsChr != "" && rsChr != "NotOn" && rsChr != "Multi" && rsChr != "Un" && rsChr != "PAR" {
					// | chrom id   |  0-filled pos  |
					// |------------|----------------|
					// |        xx  |     xxxxxxxxx  |
					// | (2 digits) |     (9 digits) |
					chrpos := lib.ChrPos(rsChr, rsPos)
					key := lib.Itob(chrpos)
					val := []byte(rsId)  // TODO: put/get rsId as byte(int)
					err = bucket.Put(key, val)
				}

				map_line, err = lib.Readln(map_reader)
			}
			if err != nil && err != io.EOF {
				return err
			}

			return nil
		})
		if err != nil {
			panic(err)
		}

		//
		os.Exit(0)
	}

	// Parse VCF header lines
	reader := bufio.NewReaderSize(os.Stdin, 64 * 1024)

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

	// Parse VCF body lines
	pattern := regexp.MustCompile(`rs(\d+)`)

	line, err = lib.Readln(reader)
	for err == nil {
		records := strings.Split(line, "\t")

		chrom := records[0]
		pos, _ := strconv.ParseInt(records[1], 10, 64)
		snpId := records[2]

		rsIdFound := pattern.FindStringSubmatch(snpId)
		// Skip or fill rs id. Switch by '-overwrite' option.
		//
		// | input     | overwrite = t | overwrite = f |
		// |-----------|---------------|---------------|
		// | "rsxxxx"  | fill          | skip          |
		// | "."       | fill          | fill          |
		if rsIdFound != nil && ! *arg_overwrite {
			// Skip
			fmt.Println(line)
		} else if (rsIdFound != nil && *arg_overwrite) || rsIdFound == nil {
			// Fill
			result := []string{}
			result = append(result, records[0:2]...)

			err = db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket(bucketName)
				if bucket == nil {
					return fmt.Errorf("Bucket %q not found!", bucketName)
				}

				val := bucket.Get(lib.Itob(lib.ChrPos(chrom, pos)))

				if val != nil {
					// Fill rs id if locus is found.
					result = append(result, "rs" + string(val))
				} else {

					if *arg_strict {
						// Fill '.' if locus in not found ('-strict' option).
						result = append(result, ".")
					} else {
						// Keep original record (including '.') if locus is not found.
						result = append(result, snpId)
					}
				}

				return nil
			})
			if err != nil {
				panic(err)
			}

			result = append(result, records[3:]...)
			fmt.Println(strings.Join(result, "\t"))
		}

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}
