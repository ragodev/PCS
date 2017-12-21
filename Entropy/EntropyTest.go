package main

import (
	"fmt"
	"math"
	"strings"
	"bytes"
	"os"
)

// read from file
func ReadFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	str := buf.String()
	//fmt.Println("successful read from file")
	return str
}

// calculate entropy
func H(data string) (entropy float64) {
	if data == "" {
		return 0
	}
	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}
	return entropy
}
func main(){
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Args error, please reinput")
		os.Exit(0)
	}
	fmt.Println(H(ReadFile(args[1])))
}

