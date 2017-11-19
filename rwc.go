// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package rwc provides a pseudo-Russian word constructor.
package rwc

var (
	vowels     = []rune("аеиоуыэюя")
	consonants = []rune("бвгджзйклмнпрстфхцчшщ")
	incv       = [32]byte{5, 0, 0, 0, 0, 8, 0, 0, 14, 0, 0, 0, 0, 0, 19, 0, 0, 0, 0, 27, 0, 0, 0, 0, 0, 0, 0, 29, 0, 30, 31, 0}
	incc       = [32]byte{0, 3, 1, 4, 6, 0, 7, 9, 0, 10, 11, 12, 13, 15, 0, 16, 17, 18, 20, 0, 21, 22, 23, 24, 25, 2, 0, 0, 0, 0, 0, 0}
)

// Constructor is a pseudo-Russian word constructor.
type Constructor struct {
	ng4    [32768]uint32
	ng3    [1024]uint32
	ng3beg [1024]uint32
	ng3end [1024]uint32
	ng2    [32]uint32
	ng1    uint32
}

// Word returns a pseudo-Russian word of the specified length.
func Word(length int) string {
	return DefaultConstructor.Word(length)
}

// WordMask returns a pseudo-Russian word matching the mask.
// The mask may contain lowercase letters from 'а' to 'я' (excluding 'ё')
// and the following special symbols:
// V - for a vowel;
// C - for a consonant;
// . (dot) - for any letter.
func WordMask(mask string) string {
	return DefaultConstructor.WordMask(mask)
}

// Word returns a pseudo-Russian word of the specified length.
func (c *Constructor) Word(n int) string {
	if n <= 0 {
		return ""
	}

	w := make([]byte, n)
	for i := 0; i < len(w); i++ {
		w[i] = byte(randA.Rand())
	}
	orig := make([]byte, n)
	copy(orig, w)

	for i := 0; i < n; {
		if c.check(w, i) {
			i++
			continue
		}

		w[i] = (w[i] + 1) % 32
		for i >= 0 && w[i] == orig[i] {
			if i == 0 {
				return ""
			}
			i--
			w[i] = (w[i] + 1) % 32
		}
	}

	return makeString(w)
}

// WordMask returns a pseudo-Russian word matching the mask.
// The mask may contain lowercase letters from 'а' to 'я' (excluding 'ё')
// and the following special symbols:
// V - for a vowel;
// C - for a consonant;
// . (dot) - for any letter.
func (c *Constructor) WordMask(mask string) string {
	if mask == "" {
		return ""
	}

	var bmask []byte
	for _, r := range mask {
		switch {
		case r == '.' || r == 'V' || r == 'C':
			bmask = append(bmask, byte(r))
		case 'а' <= r && r <= 'я':
			bmask = append(bmask, byte(r-'а'))
		case 'А' <= r && r <= 'Я':
			bmask = append(bmask, byte(r-'А'))
		case r == 'ё' || r == 'Ё':
			bmask = append(bmask, 'е'-'а')
		default:
			panic("invalid character in mask")
		}
	}

	n := len(bmask)
	w := make([]byte, n)
	for i := 0; i < n; i++ {
		switch bmask[i] {
		case '.':
			w[i] = byte(randA.Rand())
		case 'V':
			w[i] = byte(vowels[randV.Rand()] - 'а')
		case 'C':
			w[i] = byte(consonants[randC.Rand()] - 'а')
		default:
			w[i] = bmask[i]
		}
	}
	orig := make([]byte, n)
	copy(orig, w)

	for i := 0; i < n; {
		if c.check(w, i) {
			i++
			continue
		}

		w[i] = inc(w[i], bmask[i])
		for i >= 0 && w[i] == orig[i] {
			if i == 0 {
				return ""
			}
			i--
			w[i] = inc(w[i], bmask[i])
		}
	}

	return makeString(w)
}

func inc(b, how byte) byte {
	switch how {
	case '.':
		return (b + 1) % 32
	case 'V':
		return incv[b]
	case 'C':
		return incc[b]
	default:
		return b
	}
}

func makeString(w []byte) string {
	rr := make([]rune, len(w))
	for i, b := range w {
		rr[i] = 'а' + rune(b)
	}
	return string(rr)
}

func (c *Constructor) check(w []byte, i int) bool {
	good := true
	n := len(w)
	switch {
	case n == 1:
		good = c.ng1&(1<<w[0]) != 0
	case n == 2:
		if i == 1 {
			index := w[0]
			good = c.ng2[index]&(1<<w[1]) != 0
		}
	case i < 2:
	case i == 2:
		index := uint16(w[0])<<5 + uint16(w[1])
		if n == 3 {
			good = c.ng3[index]&(1<<w[2]) != 0
		} else {
			good = c.ng3beg[index]&(1<<w[2]) != 0
		}
	case i == n-1:
		index := uint16(w[n-3])<<5 + uint16(w[n-2])
		good = c.ng3end[index]&(1<<w[n-1]) != 0
	default:
		index := uint16(w[i-3])<<10 + uint16(w[i-2])<<5 + uint16(w[i-1])
		good = c.ng4[index]&(1<<w[i]) != 0
	}
	return good
}
