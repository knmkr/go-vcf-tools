package lib

import (
	"bufio"
)


func Readln(r *bufio.Reader) (string, error) {
	// http://stackoverflow.com/a/12206584
	//
	// Example:
	//
	// ```
	// reader := bufio.NewReaderSize(fp, 128 * 1024)  // 128KB
	//
	// line, err := Readln(reader)
	// for err == nil {
	//     fmt.Println(line)
	//     line, err = Readln(reader)
	// }
	// ```

	var (isPrefix bool = true
		err error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
