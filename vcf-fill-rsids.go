package main

import (
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
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s <SNPChrPosOnRef.bcp.gz> < in.vcf > out.vcf\n", os.Args[0])
		os.Exit(0)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bucket_name := []byte(path.Base(os.Args[1]))

	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	// Store chrpos <=> rsid mappings into bolt.db
	db, err := bolt.Open("bolt.db", 0644, nil)
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
			rsChr := lib.ChromToId(records[1])
			rsPos, _ := strconv.ParseInt(records[2], 10, 64)

			// | chrom id   |  0-filled pos  |
			// |------------|----------------|
			// |        xx  |     xxxxxxxxx  |
			// | (2 digits) |     (9 digits) |
			chrpos, _ := strconv.ParseInt(rsChr+fmt.Sprintf("%09d", rsPos), 10, 64)
			key := lib.Itob(chrpos)
			val := []byte(rsId)
			err = bucket.Put(key, val)

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

	// retrieve the data
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
			7011337163}

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
