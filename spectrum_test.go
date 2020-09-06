package spectrum

import (
	"math/big"
	"testing"
)

const len8 = 8
const len32 = 32
const len64 = 64
const len128 = 128

const bits8 = 0xFF
const bits32 = 0xFFFFFFFF
const bits64 uint64 = 0xFFFFFFFFFFFFFFFF
const bits128 = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF

const str32bit = "0b0000000000000000000000000000000011111111111111111111111111111111"
const str32dec = "4294967295"
const str32hex = "0x00000000ffffffff"

func TestLen(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	got := spctr.Len()
	t.Logf("Len()")
	if got != len64 {
		t.Errorf("Length expected %d, got %d", len64, got)
	}
}

func TestCopy(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: Copy()")
	_spctr := spctr.Copy()
	_spctr.bitVector = big.NewInt(1)
	if spctr.bitVector == _spctr.bitVector {
		t.Errorf("Expected call by value, but got call by reference.")
	}
}

// --- Input ---
func TestSet(t *testing.T) {
	spctr, _ := NewSpectrum(len32)

	// -- normal usecase --
	t.Logf("Exec: Set()")
	if _, err := spctr.Set(big.NewInt(bits32)); err != nil {
		t.Fatal(err)
	} else if spctr.bitVector.Cmp(big.NewInt(0)) == 0 {
		t.Errorf("bitVector expected %x, got %x", bits32, spctr.bitVector)
	}

	t.Logf("Exec: SetUint64()")
	if _, err := spctr.SetUint64(bits32); err != nil {
		t.Fatal(err)
	}

	t.Logf("Exec: SetString()")
	if _, err := spctr.SetString("FFFFFFFF", 16); err != nil {
		t.Fatal(err)
	}

	// -- exception usecase --
	t.Logf("Error handling: Set()")
	if _, err := spctr.Set(big.NewInt(0).SetUint64(bits64)); err == nil {
		t.Error("Error handling may not be appropriate.")
	}

	t.Logf("Error handling: SetUint64()")
	if _, err := spctr.SetUint64(bits64); err == nil {
		t.Error("Error handling may not be appropriate.")
	}

	t.Logf("Error handling: SetString()")
	if _, err := spctr.SetString("FFFFFFFFFFFFFFFF", 16); err == nil {
		t.Error("Error handling may not be appropriate.")
	}
}

// --- Output ---
func TestGet(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: Uint64()")
	spctr.SetUint64(bits64)
	if got := spctr.Uint64(); got != bits64 {
		t.Errorf("Uint64() expected %x, got %x", bits64, got)
	}

	t.Logf("Exec: BigInt()")
	spctr.SetUint64(bits32)
	if got := spctr.BigInt(); got.Cmp(big.NewInt(bits32)) != 0 {
		t.Errorf("BigInt() expected %x, got %x", bits32, got)
	} else if got.SetUint64(bits64); got.Cmp(spctr.BigInt()) == 0 {
		t.Logf("got %x, BigInt() %x", got, spctr.BigInt())
		t.Errorf("Expected call by value, but got call by reference.")
	}

	t.Logf("Exec: Bit()")
	if got := spctr.Bit(); got != str32bit {
		t.Errorf("Bit() expected %s, got %s", str32bit, got)
	}

	t.Logf("Exec: Text()")
	if got := spctr.Text(10); got != str32dec {
		t.Errorf("Text() expected %s, got %s", str32dec, got)
	}

	t.Logf("Exec: Hex()")
	if got := spctr.Hex(); got != str32hex {
		t.Errorf("Hex() expected %s, got %s", str32hex, got)
	} else {
		spctr7, _ := NewSpectrum(7)
		spctr7.SetUint64(127)
		if got := spctr7.Hex(); got != "0x7f" {
			t.Errorf("7bits Hex() expected %s, got %s", "0x7f", got)
		}
	}
}

// --- bits ---

func testOnesCount(t *testing.T, spctr *Spectrum, want uint) {
	if got := spctr.OnesCount(); got != want {
		t.Errorf("Case(0x%x) expected %d, got %d", spctr.bitVector, want, got)
	}
}

func TestOnesCount(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: OnesCount()")
	spctr.SetUint64(bits64)
	testOnesCount(t, spctr, 64)

	spctr.SetUint64(bits32)
	testOnesCount(t, spctr, 32)
}

func TestAdjustOnesCount(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: AdjustOnesCount()")
	spctr.AdjustOnesCount(8)
	testOnesCount(t, spctr, 8)

	spctr.AdjustOnesCount(4)
	testOnesCount(t, spctr, 4)
}

// --- rand ---

func TestSeed(t *testing.T) {
	var cp []int
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: Seed()")
	spctr.Seed(1)
	for i := 0; i < 5; i++ {
		cp = append(cp, spctr.rnd.Intn(10))
	}

	spctr.Seed(1)
	for i := 0; i < 5; i++ {
		if got := spctr.rnd.Intn(10); got != cp[i] {
			t.Errorf("rand value expected %d, got %d", cp[i], got)
		}
	}
}

// --- bitwise operation ---

func TestAnd(t *testing.T) {
	spctr32, _ := NewSpectrum(len64)
	spctr64, _ := NewSpectrum(len64)

	spctr32.SetUint64(bits32)
	spctr64.SetUint64(bits64)

	t.Logf("Exec: And()")
	if got := And(spctr64, spctr32); got.Cmp(big.NewInt(bits32)) != 0 {
		t.Errorf("Expected %x, got %x", bits32, got)
	}
}

func TestOr(t *testing.T) {
	spctr32, _ := NewSpectrum(len64)
	spctr64, _ := NewSpectrum(len64)

	spctr32.SetUint64(bits32)
	spctr64.SetUint64(bits64)

	t.Logf("Exec: Or()")
	if got := Or(spctr64, spctr32); got.Cmp(big.NewInt(0).SetUint64(bits64)) != 0 {
		t.Errorf("Expected %x, got %x", bits64, got)
	}
}

func TestAndNot(t *testing.T) {
	spctr32, _ := NewSpectrum(len64)
	spctr64, _ := NewSpectrum(len64)

	spctr32.SetUint64(bits32)
	spctr64.SetUint64(bits64 - 1)

	t.Logf("Exec: AndNot()")
	want, _ := big.NewInt(0).SetString("00000001", 16)
	if got := AndNot(spctr32, spctr64); got.Cmp(want) != 0 {
		t.Errorf("Expected %x, got %x", want, got)
	}

	want, _ = big.NewInt(0).SetString("ffffffff00000000", 16)
	if got := AndNot(spctr64, spctr32); got.Cmp(want) != 0 {
		t.Errorf("Expected %x, got %x", want, got)
	}

}

func TestXor(t *testing.T) {
	spctr32, _ := NewSpectrum(len64)
	spctr64, _ := NewSpectrum(len64)

	spctr32.SetUint64(bits32)
	spctr64.SetUint64(bits64)

	t.Logf("Exec: Xor()")
	want, _ := big.NewInt(0).SetString("ffffffff00000000", 16)
	if got := Xor(spctr64, spctr32); got.Cmp(want) != 0 {
		t.Errorf("Expected %x, got %x", want, got)
	}
}

// --- shift operation ---

func TestRsh(t *testing.T) {
	spctr, _ := NewSpectrum(len8)

	pattern := []string{
		"11111111",
		"10101010",
		"01010101",
		"00000000",
	}

	t.Logf("Exex: Rsh()")
	for _, s := range pattern {
		spctr.SetString(s, 2)
		if got := Rsh(spctr, 2); got.Bit() != "0b"+s {
			t.Errorf("Expected 0b%v, got %v", s, got.Bit())
		}
	}
}

func TestLsh(t *testing.T) {
	spctr, _ := NewSpectrum(len8)

	pattern := []string{
		"11111111",
		"10101010",
		"01010101",
		"00000000",
	}

	t.Logf("Exec: Lsh()")
	for _, s := range pattern {
		spctr.SetString(s, 2)
		if got := Lsh(spctr, 2); got.Bit() != "0b"+s {
			t.Errorf("Expected 0x%v, got %v", s, got.Bit())
		}
	}
}
