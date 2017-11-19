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

package rwc

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	rWord      = regexp.MustCompile(`(?i)[а-яё]+`)
	yoReplacer = strings.NewReplacer("ё", "е")
)

// LearnFrom learns n-grams from an UTF-8 text it reads from r.
// One can call LearnFrom multiple times with different readers.
func (c *Constructor) LearnFrom(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		c.process(s.Text())
	}
	return s.Err()
}

func (c *Constructor) process(s string) {
	// Ignore non-letter characters at the beginning and at the end of the string.
	for {
		r, size := utf8.DecodeRuneInString(s)
		if size == 0 || unicode.IsLetter(r) {
			break
		}
		s = s[size:]
	}
	for {
		r, size := utf8.DecodeLastRuneInString(s)
		if size == 0 || unicode.IsLetter(r) {
			break
		}
		s = s[:len(s)-size]
	}

	for _, m := range rWord.FindAllStringIndex(s, -1) {
		beg, end := m[0], m[1]
		if end-beg == 2 { // 2 bytes == one russian letter
			// Ignore initials and such.
			if beg > 0 {
				r, _ := utf8.DecodeLastRuneInString(s[:beg-1])
				if !unicode.IsSpace(r) {
					continue
				}
			}
			if end < len(s) {
				r, _ := utf8.DecodeRuneInString(s[end:])
				if !unicode.IsSpace(r) {
					continue
				}
			}
		}
		c.add(yoReplacer.Replace(strings.ToLower(s[beg:end])))
	}
}

func (c *Constructor) add(word string) {
	rr := []rune(word)
	n := len(rr)
	w := make([]byte, n)
	for i, r := range rr {
		w[i] = byte(r - 'а')
	}

	switch n {
	case 1:
		c.ng1 |= 1 << w[0]
	case 2:
		i := w[0]
		c.ng2[i] |= 1 << w[1]
	case 3:
		i := uint16(w[0])<<5 + uint16(w[1])
		c.ng3[i] |= 1 << w[2]
	default:
		i := uint16(w[0])<<5 + uint16(w[1])
		c.ng3beg[i] |= 1 << w[2]

		i = uint16(w[n-3])<<5 + uint16(w[n-2])
		c.ng3end[i] |= 1 << w[n-1]

		for i := 1; i < len(w)-4; i++ {
			i := uint16(w[0])<<10 + uint16(w[1])<<5 + uint16(w[2])
			c.ng4[i] |= 1 << w[3]
		}
	}
}
