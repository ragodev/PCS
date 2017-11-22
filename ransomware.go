package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"crypto/rand"
	"math/big"
	"log"
	"math"
	"encoding/hex"
	"strconv"
	"crypto/sha256"
	"crypto/aes"
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

//return the XOR on bitwise for two byte array
func XoR(b1 []byte, b2 []byte) []byte{
	xor := make([]byte, len(b1))
	for i:= 0; i < len(b1); i ++{
		xor[i] = b1[i] ^ b2[i]
	}
	//fmt.Printf("xor: %x\n", xor)
	return xor
}

//return the padding byte array for the specific message M
func PKCS_5(M []byte) []byte{

	n := int(math.Mod(float64(len(M)),16.0))

	if n == 0{
		dst, _ := hex.DecodeString(strings.Repeat("10",16))
		//fmt.Printf("padding: %x,%d\n",dst,dst)
		return dst
	}else{
		s := fmt.Sprintf("0%x", 16-n)
		//s := "0f"
		dst, _ := hex.DecodeString(strings.Repeat(s,16-n))
		//fmt.Printf("padding: %x,%d\n",dst,dst)
		return dst
	}
}
//use AES-128 in CBC mode to encrypt the message M
func AES_CBC_ENC(Kenc []byte, IV []byte, M []byte) []byte{
	var c []byte

	block, err := aes.NewCipher(Kenc)
	if err != nil{
		panic(err)
	}

	for len(M) > 0{
		block.Encrypt(IV, XoR(IV,M[:16]))
		M = M[16:]
		c = append(c, IV...)
	}
	return c
	
}
//use AES-128 in CBC mode to decrypt the ciphertext C
func AES_CBC_DEC(Kenc []byte, IV []byte, C []byte)[]byte{
	var m []byte
	temp := make([]byte, 16)

	block, err := aes.NewCipher(Kenc)
	if err != nil{
		panic(err)
	}

	for len(C) > 0{

		block.Decrypt(temp, C)
		//fmt.Println(temp, len(C))
		m = append(m, XoR(IV,temp)...)
		IV = C[:16]
		C = C[16:]
	}

	fmt.Println("IV:", IV)
	return m
}
//return the checksum from HMAC_SHA256 algorithm
func HMAC_SHA256(Kmac []byte, M []byte, B int) []byte{
	var K0 []byte

	ipad, _ := hex.DecodeString(strings.Repeat("36",B)) //inner pad
	opad, _ := hex.DecodeString(strings.Repeat("5c",B))	//outer pad

	if len(Kmac) < B{
		dst, _ := hex.DecodeString(strings.Repeat("00", B - len(Kmac)))
		K0 = append(Kmac, dst...)
	}

	h1 := sha256.New()
	h1.Write(append(XoR(K0,ipad),M...))

	h2 := sha256.New()
	h2.Write(append(XoR(K0,opad),h1.Sum(nil)...))

	return h2.Sum(nil)
}

//encrypt the message by using Key
func encrypt(Kenc []byte, Kmac []byte, M []byte) []byte{
	B := 64 //the Block size of SHA256
	T := HMAC_SHA256(Kmac,M,B)

	M_ := append(M, T...)
	PS := PKCS_5(M_)
	M__ := append(M_, PS...)

	IV := make([]byte, 16)
	rand.Read(IV)
	IV_ := make([]byte, len(IV))
	copy(IV_,IV)

	C_ := AES_CBC_ENC(Kenc, IV_, M__)
	C := append(IV, C_...)

	return C
}
//decrypt the cipthertext by using the Key
func decrypt(Kenc []byte, Kmac []byte, C []byte)[]byte{
	IV := C[:16]
	C_ := C[16:]
	M__ := AES_CBC_DEC(Kenc, IV, C_)

	last := len(M__) - 1
	s := fmt.Sprintf("%d", M__[last])
	number, _ := strconv.Atoi(s)
	//h := M__[last-number + 1:]

	for i := 1; i< number; i ++{
		if s != fmt.Sprintf("%d", M__[last - i]){
			fmt.Printf("INVALID PADDING")
			os.Exit(0)
		}
	}

	M_ := M__[:last - number + 1]
	T := M_[len(M_)-32:]
	//fmt.Println(len(M_))
	M := M_[:len(M_)-32]
	T_ := HMAC_SHA256(Kmac,M,64)
	//fmt.Println(T_)
	for i,hex := range T{
		if T_[i] != hex{
			fmt.Println("INVALID MAC")
			os.Exit(0)
		}
	}
	return M

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