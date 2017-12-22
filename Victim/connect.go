package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//connect blockchain network
//create corresponding node, 'folder' in network
func connect(addr string) {

}

//verify if the host has enough balance
func ver_bal() {

}

//verify if the attacker provide correct key
func ver_key() {

}

//provide host id and sha256 to attacker
//ask the dec key
func ask_dec() {

}

//transfer BC from host id to target node id
func transfer() {

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
	//fmt.Println(addresses())
	fmt.Println(tolerance([]bool{true, false, true, false}))
}

//
