package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//verify if the host has enough balance
func ver_bal(addr string, target string, bal float64) {
	if get_bal(addr, target) < bal {
		fmt.Printf("false")
		return
	}
	fmt.Printf("true")
}

//need each node in network check the addr in the chain.csv
//then return the balance for that addr
func get_bal(addr string, target string) float64 {
	result := 0.0
	path := fmt.Sprintf("%s/chain.txt", addr)

	//fmt.Println(path)

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatal(err)

	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')
		//		fmt.Println(line)
		lineString := strings.Split(line, " ")
		//fmt.Println(lineString[0], target)
		if lineString[0] == target {
			result, err = strconv.ParseFloat(lineString[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(lineString[1])
		}
		if err != nil {
			break
		}
	}

	return result

}

//transfer BC from host id to target node id
func transfer(addr string, value string) bool {
	addresses := addresses()
	self_dir, _ := os.Getwd()
	target := self_dir[strings.LastIndex(self_dir, "/")+1:]

	var results []bool
	for _, addr := range addresses {
		//fmt.Println(addr)
		cmd := fmt.Sprintf("%s/connect", addr)
		var out bytes.Buffer
		result := exec.Command(cmd, "vb", addr, target, value)
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
	return tolerance(results)
}

//provide host id and sha256 to attacker
//ask the dec key
func ask_key(hash_value string) string {
	dir, _ := os.Getwd()
	parent_dir := dir[:strings.LastIndex(dir, "/")]
	cmd := fmt.Sprintf("%s/Attacker/retrieve", parent_dir)

	var out bytes.Buffer
	//fmt.Println(len(hash_value))
	result := exec.Command(cmd, hash_value)
	result.Stdout = &out
	err := result.Run()

	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

//verify if the attacker provide correct key
func ver_key(N string, d string) bool {
	ext := []string{".txt", ".exe", ".doc", ".pdf", ".csv"}
	paths := checkExt(ext)
	var results []bool

	addresses := addresses()
	//var results []bool
	for _, addr_ := range addresses {
		nBig, _ := rand.Int(rand.Reader, big.NewInt(int64(len(paths)-1)))
		n := nBig.Int64()
		encrypted_data, _ := ioutil.ReadFile(paths[n])

		addr := fmt.Sprintf("%s/test.enc", addr_)
		ioutil.WriteFile(addr, encrypted_data, 0644)

		var out bytes.Buffer
		cmd := fmt.Sprintf("%s/connect", addr_)
		//fmt.Println("addr:", addr, "N:", N, "d:", d)
		result := exec.Command(cmd, "et", addr, N, d)
		result.Stdout = &out
		result.Run()

		os.Remove(addr)

		if out.String() == "true" {
			results = append(results, true)
		} else {
			results = append(results, false)
		}
	}

	return tolerance(results)
}

//return all possible encrypted files
func checkExt(ext []string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var paths []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if f.IsDir() {
			return nil
		}
		if contains(ext, filepath.Ext(path)) {
			paths = append(paths, path)
		}

		return nil
	})

	return paths
}

func contains(ss []string, s string) bool {
	for _, a := range ss {
		if a == s {
			return true
		}
	}
	return false
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
func e_test(path string, N_s string, d_s string) string {
	data, _ := ioutil.ReadFile(path)
	data_s := string(data[272:])

	h1 := H(data_s)
	Header := data[:256]

	cipher_text := data[256:]
	big_header := new(big.Int).SetBytes(Header)

	N, _ := new(big.Int).SetString(N_s, 10)
	d, _ := new(big.Int).SetString(d_s, 10)
	Kenc := Fast_exp(big_header, d, N).Bytes()
	IV := cipher_text[:16]
	C_ := cipher_text[16:]

	M := AES_CBC_DEC(Kenc, IV, C_)
	M_s := string(M[:])
	h2 := H(M_s)
	//fmt.Println("before:", h1, "after:", h2)

	if h1 > h2 {
		return "true"
	}
	return "false"
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

func Fast_exp(a, d, n *big.Int) *big.Int {
	b_d := fmt.Sprintf("%b", d)
	var x *big.Int
	x = a
	result := big.NewInt(1)

	length := len(b_d)
	//x = a
	//fmt.Println("l",length)
	for i := len(b_d) - 2; i >= 0; i-- {
		x = new(big.Int).Mod(new(big.Int).Mul(x, x), n)
		//fmt.Println(length-i,x)
		if b_d[i] == '1' {
			result.Mod(new(big.Int).Mul(result, x), n)
		}

	}
	if b_d[length-1] == '1' {
		//fmt.Println(a)
		result.Mod(new(big.Int).Mul(result, a), n)
	}
	//fmt.Println("result:", quo)
	return result
}

//return the XOR on bitwise for two byte array
func XoR(b1 []byte, b2 []byte) []byte {
	xor := make([]byte, len(b1))
	for i := 0; i < len(b1); i++ {
		xor[i] = b1[i] ^ b2[i]
	}
	//fmt.Printf("xor: %x\n", xor)
	return xor
}

//use AES-256 in CBC mode to decrypt the ciphertext C
func AES_CBC_DEC(Kenc []byte, IV []byte, C []byte) []byte {
	var m []byte
	temp := make([]byte, 16)
	//fmt.Println(Kenc)
	block, err := aes.NewCipher(Kenc)
	if err != nil {
		panic(err)
	}

	for len(C) > 0 {

		block.Decrypt(temp, C)
		//fmt.Println(temp, len(C))
		m = append(m, XoR(IV, temp)...)
		IV = C[:16]
		C = C[16:]
	}
	return m
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
		value, _ := strconv.ParseFloat(os.Args[4], 64)
		ver_bal(os.Args[2], os.Args[3], value)
	}
	if cmd == "tf" {
		temp, _ := os.Getwd()
		addr := temp[strings.LastIndex(temp, "/")+1:]
		//fmt.Println(addr)
		//target := os.Args[2]
		value := os.Args[2]
		hash_value := os.Args[3]
		//fmt.Println("before transfer")
		if transfer(addr, value) == false {
			fmt.Println("no enough balance")
			os.Exit(0)
		}
		fmt.Println("balance checking passed")
		keys_ := ask_key(hash_value)
		keys := strings.Split(keys_[strings.Index(keys_, "(")+1:strings.Index(keys_, ")")], ",")

		if ver_key(keys[0], keys[1]) == false {
			fmt.Println("key not valid")
			os.Exit(0)
		}
		pri := fmt.Sprintf("Private key: (%s,%s)", keys[0], keys[1])
		ioutil.WriteFile("keys", []byte(pri), 0644)
		fmt.Println("success")

	}
	if cmd == "ak" {
		hash_value := os.Args[2]
		keys_ := ask_key(hash_value)

		keys := strings.Split(keys_[strings.Index(keys_, "(")+1:strings.Index(keys_, ")")], ",")
		fmt.Println(ver_key(keys[0], keys[1]))
	}
	if cmd == "et" {
		fmt.Print(e_test(os.Args[2], os.Args[3], os.Args[4]))
	}

}
