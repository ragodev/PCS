package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//connect blockchain network
//create corresponding node, 'folder' in network
func connect(addr string) {

}

//verify if the host has enough balance
func ver_bal(addr string, bal float64) {
	if get_bal(addr) < bal {
		fmt.Printf("false")
	}
	fmt.Printf("true")
}

func get_bal(addr string) float64 {

	return 10.0
}

//verify if the attacker provide correct key
func ver_key() {

}

//provide host id and sha256 to attacker
//ask the dec key
func ask_dec() {

}

//transfer BC from host id to target node id
func transfer(addr string, target string, value string) {
	addresses := addresses()
	var results []bool
	for _, addr := range addresses {

		cmd := fmt.Sprintf("go run %s/connect.go", addr)
		var out bytes.Buffer
		result := exec.Command(cmd, addr, value).Stdout
		result = &out
		if out.String() == "true" {
			results = append(results, true)
		} else {
			results = append(results, false)
		}
	}

	fmt.Println(tolerance(results))
}

//fault torlerant
func tolerance(responses []bool) bool {
	count_y := 0
	count_no := 0

	for _, response := range responses {
		if response {
			count_y++
		} else {
			count_no++
		}
	}

	if count_y > count_no {
		return true
	} else {
		return false
	}

}

//test the entropy difference
func entropy() {

}

//update balance
func update() {

}

//
func addresses() []string {
	dir, _ := os.Getwd()
	parent_dir := dir[:strings.LastIndex(dir, "/")]
	network_loc := fmt.Sprintf("%s/Network/", parent_dir)

	var nodes_addr []string
	filepath.Walk(network_loc, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			nodes_addr = append(nodes_addr, path)
		}

		return nil
	})

	return nodes_addr[1:]
}

func main() {
	cmd := os.Args[1]

	if cmd == "vb" {
		value, _ := strconv.ParseFloat(os.Args[3], 64)
		ver_bal(os.Args[2], value)
	}
	if cmd == "tf" {
		addr := os.Args[2]
		target := os.Args[3]
		value := os.Args[4]
		transfer(addr, target, value)
	}
	//value := os.Args[2]
	//fmt.Println(addresses())
	//fmt.Println(tolerance([]bool{true, false, true, false}))
	//fmt.Println(ver_bal("abc", 10.13))
}

//
