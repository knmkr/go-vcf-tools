package main

import (
	"errors"
	"flag"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/knmkr/go-vcf-tools/lib"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	arg_bucket := flag.String("bucket", "", "Mappings of chrom/pos on reference genome. E.g., b142_SNPChrPosOnRef_105.bcp.gz")
	arg_setup := flag.Bool("setup", false, "Setup local db.")
	flag.Parse()

	if len(os.Args) <=2 || len(os.Args) > 4 {
		fmt.Fprintln(os.Stderr)
		flag.Usage()
		os.Exit(0)
	}

	database_name := "bolt.db"
	bucket_name := []byte(path.Base(*arg_bucket))

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

		// Store chrpos <=> rsid mappings into bolt.db
		db, err := bolt.Open(database_name, 0644, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.Batch(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucket_name)
			if err != nil {
				return err
			}

			map_reader := bufio.NewReaderSize(gz, 128*1024)
			map_line, err := lib.Readln(map_reader)
			for err == nil {
				records := strings.Split(map_line, "\t")
				rsId := records[0]
				rsChr := records[1]
				rsPos, _ := strconv.ParseInt(records[2], 10, 64)

				if rsChr != "" && rsChr != "NotOn" && rsChr != "Multi" && rsChr != "Un" && rsChr != "PAR" {
					// | chrom id   |  0-filled pos  |
					// |------------|----------------|
					// |        xx  |     xxxxxxxxx  |
					// | (2 digits) |     (9 digits) |
					chrpos, _ := strconv.ParseInt(lib.ChromToId(rsChr)+fmt.Sprintf("%09d", rsPos), 10, 64)
					key := lib.Itob(chrpos)
					val := []byte(rsId)  // TODO: put/get rsId as byte(int)
					err = bucket.Put(key, val)
					fmt.Println(chrpos)
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

	// TODO:
	// Parse VCF body lines
	line, err = lib.Readln(reader)
	for err == nil {
		// records := strings.Split(line, "\t")

		// result := []string{}
		// result = append(result, records[0:9]...)
		// result = append(result, subset(records[9:], keep_idxs)...)
		// fmt.Println(strings.Join(result, "\t"))
		fmt.Println(line)

		line, err = lib.Readln(reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}

	// TODO:
	// retrieve the data
	db, err := bolt.Open(database_name, 0600, nil)  // FIXME: read-only?
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket_name)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket_name)
		}

		// Query example
		query := []int64{13032446841,
			13032447221,
			7091839109,
			7091747130,
			7091779556,
			7092408328,
			7092373453,
			7092383887,
			7011364200,
			7011337163,
			9111718105,
			16028025216}

		for _, key := range query {
			val := bucket.Get(lib.Itob(key))
			fmt.Printf("%s\n", val)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
