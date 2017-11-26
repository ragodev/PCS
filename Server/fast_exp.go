package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	//"log"
)

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

func load_key(input_file string) (*big.Int, *big.Int, *big.Int) {

	data, _ := ioutil.ReadFile(input_file)
	data_s := string(data[:])

	keys := strings.Split(data_s[strings.Index(data_s, "(")+1:strings.Index(data_s, ")")], ",")
	//e_s := data_s[strings.Index(data_s,",") + 1:strings.Index(data_s,")")]

	N, _ := new(big.Int).SetString(keys[0], 10)
	e, _ := new(big.Int).SetString(keys[1], 10)
	d, _ := new(big.Int).SetString(keys[2], 10)
	return N, e, d

}

func main() {
	N, e, d := load_key("key_pairs.txt")
	num, _ := new(big.Int).SetString(os.Args[1], 10)
	//fmt.Println(num)
	temp := Fast_exp(num, e, N)
	fmt.Println(temp)
	fmt.Println(Fast_exp(temp, d, N))
}
