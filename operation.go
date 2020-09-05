package spectrum

import "math/big"

// --- bitwise operation (ビット演算) ---

// And は，2つのSpectrumのbitVectorをAND比較します．
func And(source *Spectrum, target *Spectrum) *big.Int {
	return big.NewInt(0).And(source.bitVector, target.bitVector)
}

// OR は，2つのSpectrumのbitVectorをOR比較します．
func Or(source *Spectrum, target *Spectrum) *big.Int {
	return big.NewInt(0).Or(source.bitVector, target.bitVector)
}

// AndNot は，2つのSpectrumのbitVectorをANDNOT比較します．
func AndNot(source *Spectrum, target *Spectrum) *big.Int {
	return big.NewInt(0).AndNot(source.bitVector, target.bitVector)
}

// Xor は，2つのSpectrumのbitVectorをXOR比較します．
func Xor(source *Spectrum, target *Spectrum) *big.Int {
	return big.NewInt(0).Xor(source.bitVector, target.bitVector)
}

// --- shift operation (シフト演算) ---

// Rsh は，bitVectorを循環論理右シフトした新しいSpectrumを返します．
func Rsh(s *Spectrum, n uint) *Spectrum {
	b := s.BigInt()
	for i := 0; i < int(n); i++ {
		if big.NewInt(0).And(b, big.NewInt(1)) == big.NewInt(1) {
			b.SetBit(b, s.Len(), 1)
		}
		b.Rsh(b, 1)
	}

	sh := s.Copy()
	sh.Set(b)
	return sh
}

// Lsh は，bitVectorを循環論理左シフトした新しいSpectrumを返します．
func Lsh(s *Spectrum, n uint) *Spectrum {
	b := s.BigInt()
	for i := 0; i < int(n); i++ {
		b.Lsh(b, 1)
		if s.Len() < b.BitLen() {
			b.SetBit(b, s.Len(), 0)
			b.SetBit(b, 0, 1)
		}
	}

	sh := s.Copy()
	sh.Set(b)
	return sh
}
