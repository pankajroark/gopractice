package main

import (
	"bufio"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	buf := make([]byte, 1024)
	f, err := os.Open("/etc/passwd")
	check(err)
	defer f.Close()

	/*
		for {
			n, err := f.Read(buf)
			if n == 0 {
				break
			}
			check(err)
			os.Stdout.Write(buf[:n])
		}
	*/
	r := bufio.NewReader(f)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for {
		n, _ := r.Read(buf)
		if n == 0 {
			break
		}
		w.Write(buf[:n])
	}
}
