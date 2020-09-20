// Package spectrum implements the bit-vector generation and manipluation functions.
// spectrum パッケージは，ビット配列を生成，操作する関数を実装します．
package spectrum

import (
	"errors"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"time"
)

// --- Spectrum ---
//
// keywords:
// spectrum, spectrum-bit, spectrum-length, OnesCount(hamming-weight, population count), one-hot, bit-vector
//
// Spectrum は，宣言時にビット長を指定し，
// また変数を隠蔽することでビット配列を意図しない変更から保護します．

// Spectrum は，spectrum情報を保持するビット配列または関数を提供する構造体です．
type Spectrum struct {
	bitVector *big.Int
	length    int
	rnd       *rand.Rand
}

// NewSpectrum は，Spectrumインターフェースを満たす構造体を宣言して返します．
func NewSpectrum(length uint) (*Spectrum, error) {
	var err error

	return &Spectrum{
		bitVector: big.NewInt(0),
		length:    int(length),
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}, err
}

// Copy は，Spectrumを複製します．
func (s *Spectrum) Copy() *Spectrum {
	ns, _ := NewSpectrum(uint(s.length))
	ns.Set(s.bitVector)

	return ns
}

// Len は，bitVectorの長さを返します．
func (s *Spectrum) Len() int {
	return s.length
}

// OnesCount は，1ビット数（hamming-weight）を返します．
func (s *Spectrum) OnesCount() uint {
	var count uint
	for _, v := range s.bitVector.Bits() {
		count += uint(bits.OnesCount(uint(v)))
	}

	return count
}

// AdjustOnesCount は，指定した1ビット数になるまでビットフラグを増減させます．
func (s *Spectrum) AdjustOnesCount(n uint) *Spectrum {
	var set uint = 1
	if uint(s.length/2) < n {
		s.bitVector.SetUint64(1)
		s.bitVector.Lsh(s.bitVector, uint(s.length))
		s.bitVector.Sub(s.bitVector, big.NewInt(1))
	}

	oc := s.OnesCount()
	if n < oc {
		set = 0
	}

	for oc != n {
		s.bitVector.SetBit(s.bitVector, s.rnd.Intn(s.length), set)
		oc = s.OnesCount()
	}

	return s
}

// Source は，Spectrumが扱う疑似乱数のSeed値を変更して再宣言します．
func (s *Spectrum) Seed(seed int64) {
	s.rnd.Seed(seed)
}

// Set は，bitVectorに値xを設定します．
func (s *Spectrum) Set(x *big.Int) (*Spectrum, error) {
	if s.length < x.BitLen() {
		return nil, errors.New("Error: bitVector is too big for length of Spectrum.")
	}

	s.bitVector.Set(x)
	return s, nil
}

// SetUint64 は，bitVectorに値xを設定します．
func (s *Spectrum) SetUint64(x uint64) (*Spectrum, error) {
	if s.length < bits.Len(uint(x)) {
		return nil, errors.New("Error: bitVector is too big for length of Spectrum.")
	}

	s.bitVector.SetUint64(x)
	return s, nil
}

// SetUint64 は，bitVectorに文字列で表現される値xを設定します．
func (s *Spectrum) SetString(str string, base int) (*Spectrum, error) {
	v, ok := big.NewInt(0).SetString(str, base)
	if !ok {
		return nil, errors.New("Error: Failed to convert string.")
	}

	return s.Set(v)
}

// IsUint64 は，bitVectorがuint64型で表現できるかを返します．
func (s *Spectrum) IsUint64() bool {
	return s.bitVector.IsUint64()
}

// Uint64 は，bitVectorを10進数のuint64型で返します．uint64で表せない場合は未定義です．
func (s *Spectrum) Uint64() uint64 {
	return s.bitVector.Uint64()
}

// BigInt は，bitVectorを10進数のbig.Int型で返します．
func (s *Spectrum) BigInt() *big.Int {
	return big.NewInt(0).Set(s.bitVector)
}

// Bit は，bitVectorを2進数表記の文字列で返します．プレフィックに"0b"が追加されます．
func (s *Spectrum) Bit() string {
	return "0b" + fmt.Sprintf("%0*s", s.length, s.bitVector.Text(2))
}

// Text は，bitVectorを指定した進数での文字列を返します．プレフィックスは追加されません．
func (s *Spectrum) Text(base int) string {
	return s.bitVector.Text(base)
}

// String は，bitVectorを指定した進数での文字列を返します．プレフィックスは追加されません．
// 実体は Text() のエイリアスです．
func (s *Spectrum) String(base int) string {
	return s.Text(base)
}

// Hex は，bitVectorを16進数表記の文字列で返します．プレフィックスに"0x"が追加されます．
func (s *Spectrum) Hex() string {
	l := s.length / 4
	if 0 < s.length%4 {
		l++
	}

	return "0x" + fmt.Sprintf("%0*s", l, s.bitVector.Text(16))
}

// Uint64n は，指定した1ビット数を持つbitVectorをuint64で返します．フラグ位置はランダムです．uint64で表せない場合は未定義です．
func (s *Spectrum) Uint64n(n uint) uint64 {
	return s.Copy().AdjustOnesCount(n).Uint64()
}

// BigIntn は，指定した1ビット数を持つbitVectorをbig.Int型で返します．フラグ位置はランダムです．
func (s *Spectrum) BigIntn(n uint) *big.Int {
	return s.Copy().AdjustOnesCount(n).BigInt()
}
