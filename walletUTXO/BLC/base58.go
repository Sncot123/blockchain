package BLC

import (
	"bytes"
	"math/big"
)

var b58Alphabet = []byte("" + "123456789" + "abcdefghijkmnopqrstuvwxyz" + "ABCDEFGHJKLMNOPQRSTUVWXYZ")

// 编码函数
func Base58Encode(input []byte) []byte {
	var result []byte //结果
	//big.int
	//byte字节数组转换为big.int
	x := big.NewInt(0).SetBytes(input)
	//求余的基本长度
	base := big.NewInt(int64(len(b58Alphabet)))
	//求余数和商
	//判断条件，除掉的最终结果是否为0
	zero := big.NewInt(0)
	//设置余数，代表base58基数表的索引位置
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}
	//返回result切片
	Reverse(result)

	//添加一个前缀,代表某类地址

	result = append([]byte{b58Alphabet[0]}, result...)

	return result
}
func Reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// base58解码函数
// 1、去掉前缀
// 2、查找input中指定数字/字符在基数表中出现的索引(mod)
// 3、余数*58
// 4、乘积结果+mod（索引）
// 5、转换为byte字节数组

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	//前缀索引
	zeroByte := 1
	data := input[zeroByte:]
	for _, b := range data {

		charIndex := bytes.IndexByte(b58Alphabet, b)

		result.Mul(result, big.NewInt(59))

		result.Add(result, big.NewInt(int64(charIndex)))
	}
	decode := result.Bytes()
	return decode
}
