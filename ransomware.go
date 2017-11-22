package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"crypto/rand"
	"math/big"
)

//addition, subtraction and multiplication
func MR_test(n *big.Int) bool{
	Big1 := big.NewInt(1)
	Big2 := big.NewInt(2)
	n_1 := new(big.Int).Sub(n,Big1)
	b_n := fmt.Sprintf("%b", n_1)

	d := new(big.Int)
	s := len(b_n) - strings.LastIndex(b_n,"1") - 1
	b_n = b_n[:len(b_n) - s]
	d.SetString(b_n,2)

	//fmt.Println(s,d)
	k := 50
	for i:= 0; i <= k; i++{

		a, _:=rand.Int(rand.Reader, new(big.Int).Sub(n_1,Big2))
		//fmt.Println(i,a)
		a.Add(a,Big2)

		x := Fast_exp(a,d,n)

		if x.Cmp(Big1) * x.Cmp(n_1) == 0{
			continue
		}

		for j:= 1; j<s;j++{
			x = Fast_exp(x,Big2,n)

			if x.Cmp(Big1) == 0{
				return false
			}

			if x.Cmp(n_1) == 0{
				break
			}
		}

		if x.Cmp(n_1) != 0{
			return false
		}
	}
	return true
}

func Fast_exp(a,d,n *big.Int) *big.Int{
	b_d := fmt.Sprintf("%b", d)
	var x *big.Int
	x = a
	result := big.NewInt(1)

	length := len(b_d)
	//x = a
	//fmt.Println("l",length)
	for i:= len(b_d) - 2; i >= 0; i--{
		x = new(big.Int).Mod(new(big.Int).Mul(x,x),n)
		//fmt.Println(length-i,x)
		if b_d[i] == '1'{
			result.Mod(new(big.Int).Mul(result,x),n)
		}

	}
	if b_d[length-1] == '1'{
		//fmt.Println(a)
		result.Mod(new(big.Int).Mul(result,a),n)
	}
	//fmt.Println("result:", quo)
	return result
}

//ab = 1(mod n), ax + by = 1 
func EEA(a_,n_ *big.Int) *big.Int{
	big0 := big.NewInt(0)
	a := new(big.Int).Add(a_,big0)
	n := new(big.Int).Add(n_,big0)
	//fmt.Println(a.String(),b)

	if a.Cmp(n) > 0{
		a = new(big.Int).Mod(a,n)
	}

	x, lastx := big.NewInt(0), big.NewInt(1)
	y, lasty := big.NewInt(1), big.NewInt(0)
	var q,m *big.Int 

	for n.Cmp(big0) != 0 {
		q,m = new(big.Int).DivMod(a,n,new(big.Int))
		//tmp = new.Sub(a,b)
		a,n = n, m
		x, lastx = new(big.Int).Sub(lastx,new(big.Int).Mul(q,x)),x
		y, lasty = new(big.Int).Sub(lasty,new(big.Int).Mul(q,y)),y
	}

	inverse := lastx
	//fmt.Println(inverse)
	if inverse.Cmp(big0) < 0{
		return new(big.Int).Add(inverse,n_)
	}
	return inverse
}

func GCD(a_,b_ *big.Int) *big.Int{

	big0 := big.NewInt(0)
	a := new(big.Int).Add(a_,big0)
	b := new(big.Int).Add(b_,big0)

	//fmt.Println("before", a,b)
	if a.Cmp(b) < 0 {
		a,b = b,a
		//fmt.Println(a,b)
	}
	
	//big0 := big.NewInt(0)

	for b.Cmp(big0) != 0 {
		a,b = b, new(big.Int).Mod(a,b)
	}

	return a
}

func key_gen() *big.Int{
	bigs := new(big.Int)
	s := strings.Repeat("1",1024)
	bigs.SetString(s,2)

	//big := new(big.Int)strings.Repeat("1",1024)
	big, _:=rand.Int(rand.Reader, bigs)
	for !MR_test(big){
		big, _ = rand.Int(rand.Reader, bigs)
	}

	return big

}

func find_rel_prime(phi *big.Int) *big.Int{
	bigs := new(big.Int)

	big1 := big.NewInt(1)
	s := strings.Repeat("1",1024)
	bigs.SetString(s,2)

	bigi, _:= rand.Int(rand.Reader, bigs)

	for GCD(phi,bigi).Cmp(big1) != 0 {
		bigi, _ = rand.Int(rand.Reader, bigs)

	}
	return bigi
}

func main(){
	big1 := big.NewInt(1)
	p := key_gen()	
	q := key_gen()

	phi := new(big.Int).Mul(new(big.Int).Sub(p,big1),new(big.Int).Sub(q,big1))
	N := new(big.Int).Mul(p,q)
	e := find_rel_prime(phi)
	d := EEA(e,phi)

	pub_key_file := os.Args[1]
	pri_key_file := os.Args[2]

	data := fmt.Sprintf("Public key. (%s,%s)",N,e)
	//fmt.Println(data)
	ioutil.WriteFile(pub_key_file,[]byte(data),0644)

	data = fmt.Sprintf("Private key. (%s,%s,%s,%s)",N,d,p,q)
	//fmt.Println(data)
	ioutil.WriteFile(pri_key_file,[]byte(data),0644)
	//fmt.Println(pub_key_file,pri_key_file)

}