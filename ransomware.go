package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"

	//"log"
	"crypto/aes"
	"encoding/hex"
	"math"
)

//addition, subtraction and multiplication
func MR_test(n *big.Int) bool {
	Big1 := big.NewInt(1)
	Big2 := big.NewInt(2)
	n_1 := new(big.Int).Sub(n, Big1)
	b_n := fmt.Sprintf("%b", n_1)

	d := new(big.Int)
	s := len(b_n) - strings.LastIndex(b_n, "1") - 1
	b_n = b_n[:len(b_n)-s]
	d.SetString(b_n, 2)

	//fmt.Println(s,d)
	k := 50
	for i := 0; i <= k; i++ {

		a, _ := rand.Int(rand.Reader, new(big.Int).Sub(n_1, Big2))
		//fmt.Println(i,a)
		a.Add(a, Big2)

		x := Fast_exp(a, d, n)

		if x.Cmp(Big1)*x.Cmp(n_1) == 0 {
			continue
		}

		for j := 1; j < s; j++ {
			x = Fast_exp(x, Big2, n)

			if x.Cmp(Big1) == 0 {
				return false
			}

			if x.Cmp(n_1) == 0 {
				break
			}
		}

		if x.Cmp(n_1) != 0 {
			return false
		}
	}
	return true
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

//ab = 1(mod n), ax + by = 1
func EEA(a_, n_ *big.Int) *big.Int {
	big0 := big.NewInt(0)
	a := new(big.Int).Add(a_, big0)
	n := new(big.Int).Add(n_, big0)
	//fmt.Println(a.String(),b)

	if a.Cmp(n) > 0 {
		a = new(big.Int).Mod(a, n)
	}

	x, lastx := big.NewInt(0), big.NewInt(1)
	y, lasty := big.NewInt(1), big.NewInt(0)
	var q, m *big.Int

	for n.Cmp(big0) != 0 {
		q, m = new(big.Int).DivMod(a, n, new(big.Int))
		//tmp = new.Sub(a,b)
		a, n = n, m
		x, lastx = new(big.Int).Sub(lastx, new(big.Int).Mul(q, x)), x
		y, lasty = new(big.Int).Sub(lasty, new(big.Int).Mul(q, y)), y
	}

	inverse := lastx
	//fmt.Println(inverse)
	if inverse.Cmp(big0) < 0 {
		return new(big.Int).Add(inverse, n_)
	}
	return inverse
}

func GCD(a_, b_ *big.Int) *big.Int {

	big0 := big.NewInt(0)
	a := new(big.Int).Add(a_, big0)
	b := new(big.Int).Add(b_, big0)

	//fmt.Println("before", a,b)
	if a.Cmp(b) < 0 {
		a, b = b, a
		//fmt.Println(a,b)
	}

	//big0 := big.NewInt(0)

	for b.Cmp(big0) != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}

	return a
}

func key_gen() *big.Int {
	bigs := new(big.Int)
	s := strings.Repeat("1", 1024)
	bigs.SetString(s, 2)

	//big := new(big.Int)strings.Repeat("1",1024)
	big, _ := rand.Int(rand.Reader, bigs)
	for !MR_test(big) {
		big, _ = rand.Int(rand.Reader, bigs)
	}

	return big

}

func find_rel_prime(phi *big.Int) *big.Int {
	bigs := new(big.Int)

	big1 := big.NewInt(1)
	s := strings.Repeat("1", 1024)
	bigs.SetString(s, 2)

	bigi, _ := rand.Int(rand.Reader, bigs)

	for GCD(phi, bigi).Cmp(big1) != 0 {
		bigi, _ = rand.Int(rand.Reader, bigs)

	}
	return bigi
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

//return the padding byte array for the specific message M
func PKCS_5(M []byte) []byte {
	dst, _ := hex.DecodeString("")
	n := int(math.Mod(float64(len(M)), 32.0))

	if n == 0 {
		dst, _ = hex.DecodeString(strings.Repeat("20", 32))
		//fmt.Printf("padding: %x,%d\n",dst,dst)
		//return dst
	} else if n > 16 {
		s := fmt.Sprintf("0%x", 32-n)
		dst, _ = hex.DecodeString(strings.Repeat(s, 32-n))

	} else {
		s := fmt.Sprintf("%x", 32-n)
		dst, _ = hex.DecodeString(strings.Repeat(s, 32-n))

	}
	//fmt.Printf("need length: %d, padding: %x with length %d\n", 32-n, dst, len(dst))
	return dst

}

//use AES-256 in CBC mode to encrypt the message M
func AES_CBC_ENC(Kenc []byte, IV []byte, M []byte) []byte {
	var c []byte

	block, err := aes.NewCipher(Kenc)
	if err != nil {
		panic(err)
	}

	for len(M) > 0 {
		block.Encrypt(IV, XoR(IV, M[:32]))
		M = M[32:]
		c = append(c, IV...)
	}
	return c

}

//use AES-128 in CBC mode to decrypt the ciphertext C
func AES_CBC_DEC(Kenc []byte, IV []byte, C []byte) []byte {
	var m []byte
	temp := make([]byte, 32)

	block, err := aes.NewCipher(Kenc)
	if err != nil {
		panic(err)
	}

	for len(C) > 0 {

		block.Decrypt(temp, C)
		//fmt.Println(temp, len(C))
		m = append(m, XoR(IV, temp)...)
		IV = C[:32]
		C = C[32:]
	}

	//fmt.Println("IV:", IV)
	return m
}

//return the checksum based on SHA-512
func HMAC_SHA256(Kmac []byte, M []byte) []byte {
	hmac256 := hmac.New(sha256.New, Kmac)
	hmac256.Write(M)
	return hmac256.Sum(nil)

}

//encrypt the message by using Key
func encrypt(N *big.Int, e *big.Int, M []byte) []byte {
	Kmac := []byte("secret")
	T := HMAC_SHA256(Kmac, M)

	M_ := append(M, T...)
	PS := PKCS_5(M_)
	M__ := append(M_, PS...)
	fmt.Printf("Msg: %x\n", M__)
	IV := make([]byte, 32)
	rand.Read(IV)
	IV_ := make([]byte, len(IV))
	copy(IV_, IV)

	Kenc := make([]byte, 32)
	rand.Read(Kenc)
	C_ := AES_CBC_ENC(Kenc, IV_, M__)
	C := append(IV, C_...)
	fmt.Printf("MSG: %x\n", C)
	decrypt(Kenc, C)
	big_Kenc := new(big.Int).SetBytes(Kenc)
	//fmt.Printf("IV:%x with length:%d\n", IV, len(IV))
	//fmt.Printf("Kenc:%x with length:%d\n", Kenc, len(Kenc))
	tmp := Fast_exp(big_Kenc, e, N)
	//fmt.Printf("Big_enc:%d\n", big_Kenc)
	//fmt.Printf("encrption:%d\n", tmp)
	Kenc_encrypted := tmp.Bytes()
	//fmt.Printf("Header:%x\t with length:%d\n", Kenc_encrypted, len(Kenc_encrypted))
	return append(Kenc_encrypted, C...)
}

//decrypt the cipthertext by using the Key
func decrypt(Kenc []byte, C []byte) []byte {
	fmt.Println("running")
	Kmac := []byte("secret")
	IV := C[:32]
	C_ := C[32:]
	M__ := AES_CBC_DEC(Kenc, IV, C_)

	last := len(M__) - 1
	s := fmt.Sprintf("%d", M__[last])
	fmt.Printf("Msg: %x\n", M__)
	number, _ := strconv.Atoi(s)
	//padding := fmt.Sprintf("%d", M__[last-number:])
	//fmt.Println("checking padding", padding, number)
	for i := 1; i < number; i++ {
		if s != fmt.Sprintf("%d", M__[last-i]) {
			fmt.Println("INVALID PADDING")
			os.Exit(0)
		}
	}

	M_ := M__[:last-number+1]
	T := M_[len(M_)-32:]
	//fmt.Println(len(M_))
	M := M_[:len(M_)-32]
	T_ := HMAC_SHA256(Kmac, M)
	//fmt.Println(T_)
	fmt.Println("checking MAC")
	for i, hex := range T {
		if T_[i] != hex {
			fmt.Println("INVALID MAC")
			os.Exit(0)
		}
	}
	fmt.Println("GOOD!")
	return M

}

func read_file(input_file string) []byte {

	data, err := ioutil.ReadFile(input_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "guess_keylength: %v\n", err)
	}

	size := len(data)
	M := make([]byte, size)

	for i, bit_8 := range data {
		M[i] = bit_8
	}

	return M
}

func load_key(input_file string) (*big.Int, *big.Int) {

	data, _ := ioutil.ReadFile(input_file)
	data_s := string(data[:])

	keys := strings.Split(data_s[strings.Index(data_s, "(")+1:strings.Index(data_s, ")")], ",")
	//e_s := data_s[strings.Index(data_s,",") + 1:strings.Index(data_s,")")]

	N, _ := new(big.Int).SetString(keys[0], 10)
	d, _ := new(big.Int).SetString(keys[1], 10)
	return N, d

}

func main() {
	N, _ := new(big.Int).SetString("27373176497986431932251251116443926663831806261071313412989844158826154814084495647335491185915655601666177981912915917513154667109854715685389256347698156167098286352963470578246399719078975590638921432718527013913639908860623319460907657961571096028464789764100067928297618269506856287203984333985334991680977184069499150590693164224823071048888338148038226506445470031874898694438134368521730060168857461554576831837880299721200508752202214917227818509186371000166250119754889508590004281188893160298827566989760562889798571576260058420518940759480974818966399516795337246357463835643099652311241207979328006473273", 10)
	e, _ := new(big.Int).SetString("114555885618321971166525289009247562530684912588490288375374382656653624817637873641873338444064827244125181436880632068365026025194100187802604323676970718988010331197271572511319987984765105343287840472134128084868761183464484207489195971498946203121897312631812237101119681888221283146033802960204006054855", 10)

	mode := os.Args[1]
	file := os.Args[2]
	if mode == "-e" {
		encrypted_data := encrypt(N, e, read_file(file))
		ioutil.WriteFile("encrypted.txt", encrypted_data, 0644)
	} else if mode == "-d" {
		//Decrypt
		_, d := load_key(os.Args[3])
		encrypted_data := read_file(os.Args[2])
		Header := encrypted_data[:256]
		//fmt.Printf("Header:%x\t with length:%d\n", Header, len(Header))

		cipher_text := encrypted_data[256:]
		big_header := new(big.Int).SetBytes(Header)
		Kenc := Fast_exp(big_header, d, N).Bytes()
		//fmt.Printf("Kenc:%x with length:%d\n", Kenc, len(Kenc))

		plain_text := decrypt(Kenc, cipher_text)
		ioutil.WriteFile("decrypted.txt", plain_text, 0644)
		//pri_key_file := os.Args[3]
	}
	// big1 := big.NewInt(1)
	// p := key_gen()
	// q := key_gen()

	// phi := new(big.Int).Mul(new(big.Int).Sub(p,big1),new(big.Int).Sub(q,big1))
	// N := new(big.Int).Mul(p,q)
	// e := find_rel_prime(phi)
	// d := EEA(e,phi)

	// pub_key_file := os.Args[1]
	// pri_key_file := os.Args[2]

	// data := fmt.Sprintf("Public key. (%s,%s)",N,e)
	// //fmt.Println(data)
	// ioutil.WriteFile(pub_key_file,[]byte(data),0644)

	// data = fmt.Sprintf("Private key. (%s,%s,%s,%s)",N,d,p,q)
	// //fmt.Println(data)
	// ioutil.WriteFile(pri_key_file,[]byte(data),0644)
	// //fmt.Println(pub_key_file,pri_key_file)

}
