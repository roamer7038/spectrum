package spectrum

import (
	"math/big"
	"testing"
)

const len32 = 32
const len64 = 64
const bits32 = 0xFFFFFFFF
const bits64 uint64 = 0xFFFFFFFFFFFFFFFF

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

// --- 入力系 ---
func TestSet(t *testing.T) {
	spctr, _ := NewSpectrum(len32)

	// -- 正常系 --
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

	// -- 異常系 --
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

// --- 出力系 ---
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

func TestCopy(t *testing.T) {
	spctr, _ := NewSpectrum(len64)

	t.Logf("Exec: Copy()")
	_spctr := spctr.Copy()
	_spctr.bitVector = big.NewInt(1)
	if spctr.bitVector == _spctr.bitVector {
		t.Errorf("Expected call by value, but got call by reference.")
	}
}
