package k0yote3web

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
)

type EtherUnit int

const (
	EtherUnitNoEther EtherUnit = iota
	EtherUnitWei
	EtherUnitKWei
	EtherUnitMWei
	EtherUnitGWei
	EtherUnitSzabo
	EtherUnitFinney
	EtherUnitEther
)

var unitMap = map[EtherUnit]string{
	EtherUnitNoEther: "0",
	EtherUnitWei:     "1",
	EtherUnitKWei:    "1000",
	EtherUnitMWei:    "1000000",
	EtherUnitGWei:    "1000000000",
	EtherUnitSzabo:   "1000000000000",
	EtherUnitFinney:  "1000000000000000",
	EtherUnitEther:   "1000000000000000000",
}

func FromWei(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func FromGWei(wei *big.Int) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func FromWeiFloat(wei *big.Float) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	ret := bigval.Quo(wei, expF)
	return ret
}

func FromDecimals(wei *big.Int, decimals int64) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)
	ret := bigval.Quo(bigval, expF)
	return ret
}

func ToWeiInt(val int64, denominator int64) *big.Int {
	bigval := big.NewInt(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	result := bigval.Mul(bigval, exp)
	result = big.NewInt(1).Div(result, big.NewInt(denominator))
	return result
}

func FromWeiWithUnit(wei *big.Int, unit EtherUnit) *big.Float {
	if wei == nil {
		return big.NewFloat(0)
	}
	unitInt := 0
	switch unit {
	case EtherUnitNoEther:
		unitInt = 0
	case EtherUnitWei:
		unitInt = 1
	case EtherUnitKWei:
		unitInt = 3
	case EtherUnitMWei:
		unitInt = 6
	case EtherUnitGWei:
		unitInt = 9
	case EtherUnitSzabo:
		unitInt = 12
	case EtherUnitFinney:
		unitInt = 15
	case EtherUnitEther:
		unitInt = 18
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(unitInt)), nil)
	expF := new(big.Float)
	expF.SetInt(exp)

	bigval := new(big.Float)
	bigval.SetInt(wei)

	ret := bigval.Quo(bigval, expF)
	return ret
}

func ToHex(n *big.Int) string {
	return fmt.Sprintf("0x%x", n) // or %X or upper case
}

func HexToUint64(str string) (uint64, error) {
	return ParseUint64orHex(str)
}

func ToDecimals(val uint64, decimals int64) *big.Int {
	return convert(val, decimals)
}

func SameAddress(a, b common.Address) bool {
	return bytes.Equal(a[:], b[:])
}

func DifferentAddress(a, b common.Address) bool {
	return !bytes.Equal(a[:], b[:])
}

func RoundNWei(wei *big.Int, n int) (*big.Int, error) {
	af := FromWei(wei)
	aff, _ := af.Float64()

	roundfs := ""
	if n > 6 {
		return nil, fmt.Errorf("round n not support bigger than 6")
	}
	switch n {
	case 1:
		roundfs = fmt.Sprintf("%.1f", aff)

	case 2:
		roundfs = fmt.Sprintf("%.2f", aff)

	case 3:
		roundfs = fmt.Sprintf("%.3f", aff)

	case 4:
		roundfs = fmt.Sprintf("%.4f", aff)

	case 5:
		roundfs = fmt.Sprintf("%.5f", aff)

	case 6:
		roundfs = fmt.Sprintf("%.6f", aff)

	}

	r := ToWeiR(roundfs)

	return r, nil
}

// Ether converts a value to the ether unit with 18 decimals
func Ether(i uint64) *big.Int {
	return convert(i, 18)
}

func convert(val uint64, decimals int64) *big.Int {
	v := big.NewInt(int64(val))
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	return v.Mul(v, exp)
}

// ToBlockNumArg. Wrap blockNumber arg from big.Int to string
func ToBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}

func AddressChecksum(iaddress interface{}) bool {
	reg := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	switch v := iaddress.(type) {
	case string:
		return reg.MatchString(v)
	case common.Address:
		return reg.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}
	amount.Float64()
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

func ToWeiR(val string) *big.Int {

	if !strings.Contains(val, ".") {
		whole, ok := big.NewInt(0).SetString(val, 10)
		if !ok {
			return big.NewInt(0)
		}
		return big.NewInt(1).Mul(whole, big.NewInt(1e18))
	}
	comps := strings.Split(val, ".")
	if len(comps) != 2 {
		return big.NewInt(0)
	}

	whole := comps[0]
	fraction := comps[1]
	baseLength := len(unitMap[EtherUnitEther]) - 1
	fractionLength := len(fraction)
	if fractionLength > baseLength {
		return big.NewInt(0)
	}
	fraction += strings.Repeat("0", baseLength-fractionLength)
	wholeInt, ok := big.NewInt(0).SetString(whole, 10)
	if !ok {
		return big.NewInt(0)
	}
	fractionInt, ok := big.NewInt(0).SetString(fraction, 10)
	if !ok {
		return big.NewInt(0)
	}

	wholeMulBase := big.NewInt(1).Mul(wholeInt, big.NewInt(1e18))
	wholeAddFraction := big.NewInt(1).Add(wholeMulBase, fractionInt)

	return wholeAddFraction
}

func ToGWei(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)

	expF := new(big.Float)
	expF.SetInt(exp)

	bigval.Mul(bigval, expF)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func ToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func StringToBytes32(str string) [32]byte {
	var a *[32]byte
	b := crypto.Keccak256([]byte(str))
	if len(a) <= len(b) {
		a = (*[len(a)]byte)(unsafe.Pointer(&b[0]))
	}
	return *a
}

func EncodeUintToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

func ParseBigInt(str string) *big.Int {
	str = strings.TrimPrefix(str, "0x")
	num := new(big.Int)
	num.SetString(str, 16)
	return num
}

func ParseUint64orHex(str string) (uint64, error) {
	base := 10
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
		base = 16
	}
	return strconv.ParseUint(str, base, 64)
}

func EncodeToHex(b []byte) string {
	return hexutil.Encode(b)
}

func ParseHexBytes(str string) ([]byte, error) {
	return hexutil.Decode(str)
}
