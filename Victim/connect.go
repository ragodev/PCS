package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//connect blockchain network
//create corresponding node, 'folder' in network
// func connect(addr string) {
// 	paths := addresses()
// 	find := false

// 	for _, path := range paths{
// 		addr_ := path[strings.LastIndex(temp, "/")+1:]
// 		if addr_ == addr{
// 			find = true
// 		}
// 	}

// 	if addr
// }

//verify if the host has enough balance
func ver_bal(addr string, bal float64) {
	if get_bal(addr) < bal {
		fmt.Printf("false")
	}
	fmt.Printf("true")
}

//need each node in network check the addr in the chain.csv
//then return the balance for that addr
func get_bal(addr string) float64 {

	return 10.0
}

//transfer BC from host id to target node id
func transfer(addr string, target string, value string) bool {
	addresses := addresses()
	var results []bool
	for _, addr := range addresses {

		cmd := fmt.Sprintf("%s/connect", addr)
		var out bytes.Buffer
		result := exec.Command(cmd, "vb", addr, value)
		result.Stdout = &out
		err := result.Run()

		if err != nil {
			log.Fatal(err)

		}
		//fmt.Println(out.String())
		if out.String() == "true" {
			results = append(results, true)
		} else {
			results = append(results, false)
		}
	}
	//fmt.Println(results)
	//fmt.Println(tolerance(results))
	return tolerance(results)
}

//provide host id and sha256 to attacker
//ask the dec key
func ask_dec(hash_value string) string {
	dir, _ := os.Getwd()
	parent_dir := dir[:strings.LastIndex(dir, "/")]
	cmd := fmt.Sprintf("%s/Attacker/retrieve", parent_dir)

	var out bytes.Buffer
	result := exec.Command(cmd, hash_value)
	result.Stdout = &out
	err := result.Run()

	if err != nil {
		log.Fatal(err)
	}
	return out.String()
	// keys := strings.Split(result[out.String().Index(out.String(), "(")+1:strings.Index(out.Sring(), ")")], ",")
	// N, _ := new(big.Int).SetString(keys[0], 10)
	// d, _ := new(big.Int).SetString(keys[1], 10)
}

//verify if the attacker provide correct key
func ver_key() {

}

func fin_tr() {

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
		temp, _ := os.Getwd()
		addr := temp[strings.LastIndex(temp, "/")+1:]
		//fmt.Println(addr)
		target := os.Args[2]
		value := os.Args[3]
		transfer(addr, target, value)
	}
	//value := os.Args[2]
	//fmt.Println(addresses())
	//fmt.Println(tolerance([]bool{true, false, true, false}))
	//fmt.Println(ver_bal("abc", 10.13))
}

//
