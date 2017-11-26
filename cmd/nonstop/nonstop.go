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

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
	"unicode"

	"github.com/opennota/rwc"
)

var (
	numLines       = flag.Int("l", 0, "Number of lines to print (default 0 = infinity)")
	vocabularyFile = flag.String("v", "", "Vocabulary file")
)

const (
	X_PARAGRAPH = iota
	X_POINT
	X_THREEPOINT
	X_EXCLAM
	X_COMMA
	X_DASH
	X_COLON
	X_COMMADASH
	X_SEMICOLON
	X_OQUOTE
	X_CQUOTE
	X_OBRACKET
	X_CBRACKET
	X_SENTBRAK
	X_NAME
	X_QUESTION
	X_SPECIAL
	X_EOS
	X_OITALIC
	X_CITALIC
	X_DIRECT
)

const (
	BOS uint32 = 1 << iota
	EOS
	NOONE
	PUN
	QUOT
	BRAK
	ABZ
	WASPUN
	SENTBRAK
	WASONE
	WASCOLON
	ITAL
	WASDIRECT
)

var poss = []struct {
	value     float64
	initial   float64
	increment float64
}{
	{0, 15, 10},  // par
	{0, 15, 10},  // .
	{0, 0, 1.5},  // ...
	{0, 0, 2},    // !
	{0, 18, 15},  // ,
	{0, 0, 1},    // -
	{0, 2, 3},    // :
	{0, 2, 2},    // , -
	{0, 0, .5},   // ;
	{0, 6.5, 0},  // <<
	{0, 50, 20},  // >>
	{0, 4, 0},    // (
	{0, 50, 40},  // )
	{0, 10, -10}, // (...)
	{0, 4, 0},    // Name
	{0, 0, 2},    // ?
	{0, 100, 0},  // Спецсимвол
	{0, 15, 17},  // Конец предложения
	{0, .2, .05}, // {\it
	{0, 70, 30},  // }
	{0, 1, 1.5},  // : <
	{0, 0, 0},    // а
	{0, 0, 0},    // б
	{0, 0, 0},    // в
	{0, 0, 0},    // ж
	{0, 0, 0},    // и
	{0, 0, 0},    // к
	{0, 0, 0},    // о
	{0, 0, 0},    // с
	{0, 0, 0},    // у
	{0, 0, 0},    // я
}

func init() {
	for i, p := range poss {
		poss[i].value = p.initial
	}
}

func possibleX(x int) bool {
	return possible(poss[x].value)
}

func possible(p float64) bool {
	return p > float64(rand.Intn(256))
}

func resetX(x int) {
	poss[x].value = poss[x].initial
}

func advanceX(x int) {
	p := &poss[x]

	if p.increment == 0 {
		p.value += (float64(rand.Intn(int(p.initial*300))) - p.initial*150) / 1000
		if p.value < p.initial*.5 || p.value > p.initial*2 {
			p.value = p.initial
		}
		return
	}

	if p.initial >= 256 {
		decreaseX(x)
		return
	}

	if p.value < p.initial {
		p.value = p.initial
	}
	a := math.Sqrt((p.value-p.initial)/(256-p.initial))*256/p.increment + 1
	p.value = (256-p.initial)*p.increment*p.increment*a*a/(256*256) + p.initial
}

func decreaseX(x int) {
	poss[x].initial /= 1.5
	resetX(x)
}

func increaseX(x int) {
	if poss[x].initial*1.5 >= 256 {
		return
	}
	poss[x].initial *= 1.5
	advanceX(x)
}

func LoadVocabulary(filename string) error {
	if err := rwc.DefaultConstructor.LoadFromRWC(filename); err != nil {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := rwc.DefaultConstructor.LoadFrom(f); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *numLines == 0 {
		*numLines = math.MaxInt64
	}
	if *vocabularyFile != "" {
		if err := LoadVocabulary(*vocabularyFile); err != nil {
			log.Fatal(err)
		}
	}
	rand.Seed(time.Now().UnixNano())

	flags, old_flags := BOS, BOS
	var lastopened, old_lastopened uint32

	str := make([]rune, 0, 255)
	var word []rune
	w := make([]rune, 35)

	l := 100

	for lines := 0; lines < *numLines; {
		word = word[:0]

		if flags&BOS != 0 && possibleX(X_SENTBRAK) {
			flags |= SENTBRAK
			word = append(word, '(')
		}

		if (flags&PUN != 0) && (flags&WASPUN == 0) && (flags&ITAL == 0) {
			if flags&BOS == 0 && flags&BRAK == 0 && flags&SENTBRAK == 0 && possibleX(X_OBRACKET) {
				word = append(word, '(')
				resetX(X_CBRACKET)
				flags |= BRAK | WASPUN
				lastopened = BRAK
			} else if flags&QUOT == 0 && possibleX(X_OQUOTE) {
				word = append(word, '«')
				resetX(X_CQUOTE)
				flags |= QUOT | WASPUN
				lastopened = QUOT
			}
		}

		if possible(35) {
			if flags&NOONE != 0 {
				flags, lastopened = old_flags, old_lastopened
				continue
			}
			l = 1

			flags |= NOONE
			if possible(128) {
				w[0] = 'и'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(80) && flags&WASPUN != 0 {
				w[0] = 'а'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(20) {
				w[0] = 'в'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(20) {
				w[0] = 'к'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(20) {
				w[0] = 'с'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(15) {
				w[0] = 'о'
				flags &= ^(PUN | EOS | WASPUN)
				if flags&BOS != 0 {
					flags |= PUN
				}
			} else if possible(5) {
				w[0] = 'у'
				flags &= ^(PUN | EOS | WASPUN)
			} else if possible(4) {
				w[0] = 'я'
				flags |= PUN | EOS
				flags &= ^WASPUN
			} else if possible(2) && flags&BOS == 0 && flags&WASPUN == 0 {
				w[0] = 'б'
				flags |= PUN | EOS
			} else if flags&BOS == 0 && flags&WASPUN == 0 {
				w[0] = 'ж'
				if possible(90) {
					w[l] = 'е'
					l++
				}
				flags |= EOS
			} else {
				flags, lastopened = old_flags, old_lastopened
				continue
			}
			flags |= WASONE
		} else if possible(29) {
			l = 2
			copy(w, []rune(rwc.Word(2)))
			flags &= ^(PUN | EOS | WASPUN | WASONE)
		} else if possible(24) {
			l = 3
			copy(w, []rune(rwc.Word(3)))
			flags &= ^(NOONE | WASPUN | WASONE)
			flags |= (PUN | EOS)
		} else {
			if possible(18) {
				l = 6 + rand.Intn(8)
			} else {
				l = 4 + rand.Intn(4)
			}
			copy(w, []rune(rwc.Word(l)))
			flags &= ^(NOONE | WASONE | WASPUN)
			flags |= (PUN | EOS)
		}

		if flags&BOS != 0 {
			flags &= ^BOS
			w[0] = unicode.ToUpper(w[0])
		} else if flags&WASONE == 0 && possible(poss[X_NAME].value*float64(l-2)*1.5) {
			w[0] = unicode.ToUpper(w[0])
		}

		if flags&ITAL == 0 && possibleX(X_OITALIC) {
			word = append(word, '_')
			flags |= ITAL
			resetX(X_CITALIC)
		}

		if flags&PUN != 0 && flags&WASPUN == 0 {
			if possibleX(X_COMMA) {
				w[l] = ','
				l++
				resetX(X_COMMA)
				flags |= WASPUN
			} else if possibleX(X_DASH) {
				w[l] = ' '
				l++
				w[l] = '—'
				l++
				resetX(X_DASH)
				flags |= WASPUN
			} else if flags&QUOT != 0 && lastopened != BRAK && flags&ITAL == 0 && possibleX(X_CQUOTE) {
				w[l] = '»'
				l++
				flags &= ^QUOT
				flags |= WASPUN
				if flags&BRAK != 0 {
					lastopened = BRAK
				} else {
					lastopened = 0
				}
				resetX(X_OQUOTE)
			} else if flags&BRAK != 0 && lastopened != QUOT && flags&ITAL == 0 && possibleX(X_CBRACKET) {
				w[l] = ')'
				l++
				flags &= ^BRAK
				flags |= WASPUN
				if flags&QUOT != 0 {
					lastopened = QUOT
				} else {
					lastopened = 0
				}
				resetX(X_OBRACKET)
			} else if flags&ITAL != 0 && possibleX(X_CITALIC) {
				w[l] = '_'
				l++
				flags &= ^ITAL
				resetX(X_OITALIC)
			} else if possibleX(X_SPECIAL) {
				if possibleX(X_SEMICOLON) {
					w[l] = ';'
					l++
					resetX(X_SEMICOLON)
					flags |= WASPUN
				} else if flags&WASCOLON == 0 && possibleX(X_COLON) {
					w[l] = ':'
					l++
					resetX(X_COLON)
					increaseX(X_COMMA)
					advanceX(X_COMMADASH)
					advanceX(X_COMMADASH)
					flags |= WASPUN | WASCOLON
				} else if possibleX(X_COMMADASH) {
					w[l] = ','
					l++
					w[l] = ' '
					l++
					w[l] = '—'
					l++
					resetX(X_COMMADASH)
					flags |= WASPUN
					if flags&WASCOLON != 0 {
						decreaseX(X_COMMA)
						flags &= ^WASCOLON
					}
				} else {
					flags, lastopened = old_flags, old_lastopened
					continue
				}
			} else if flags&EOS != 0 && flags&WASPUN == 0 && possibleX(X_EOS) {
				for lastopened != 0 {
					if flags&BRAK != 0 && lastopened == BRAK {
						w[l] = ')'
						l++
						flags &= ^BRAK
						if flags&QUOT != 0 {
							lastopened = QUOT
						} else {
							lastopened = 0
						}
					}
					if flags&QUOT != 0 && lastopened == QUOT {
						w[l] = '»'
						flags &= ^QUOT
						if flags&BRAK != 0 {
							lastopened = BRAK
						} else {
							lastopened = 0
						}
					}
				}

				if possibleX(X_POINT) {
					w[l] = '.'
					l++
					resetX(X_POINT)
				} else if possibleX(X_EXCLAM) {
					w[l] = '!'
					l++
					resetX(X_EXCLAM)
				} else if possibleX(X_THREEPOINT) {
					w[l] = '…'
					l++
					resetX(X_THREEPOINT)
				} else if possibleX(X_QUESTION) {
					w[l] = '?'
					l++
					resetX(X_QUESTION)
				} else {
					flags, lastopened = old_flags, old_lastopened
					continue
				}

				if flags&ITAL != 0 {
					w[l] = '_'
					l++
					flags &= ^ITAL
					resetX(X_OITALIC)
				}
				if flags&SENTBRAK != 0 {
					w[l] = ')'
					l++
					flags &= ^SENTBRAK
					resetX(X_SENTBRAK)
				}
				if flags&WASCOLON != 0 {
					decreaseX(X_COMMA)
					flags &= ^WASCOLON
				}
				resetX(X_EOS)

				flags |= BOS | WASPUN
				flags &= ^EOS | PUN
				advanceX(X_PARAGRAPH)
				if possibleX(X_PARAGRAPH) {
					flags |= ABZ
					resetX(X_PARAGRAPH)
				}
			}
		}

		if flags&WASDIRECT == 0 {
			w[l] = ' '
			l++
		}

		word = append(word, w[:l]...)
		fmt.Print(string(word))
		str = append(str, word...)

		if flags&ABZ != 0 {
			flags &= ^ABZ
			fmt.Print("\n\n")
			str = str[:0]
			lines++
		} else if len(str) > 55 && flags&WASDIRECT == 0 {
			fmt.Print("\n")
			str = str[:0]
			lines++
		}

		flags &= ^WASDIRECT

		advanceX(X_COMMA)
		advanceX(X_EOS)
		advanceX(X_POINT)
		advanceX(X_THREEPOINT)
		advanceX(X_EXCLAM)
		advanceX(X_QUESTION)
		advanceX(X_DASH)
		advanceX(X_COLON)
		advanceX(X_COMMADASH)
		advanceX(X_SEMICOLON)
		advanceX(X_OQUOTE)
		advanceX(X_OBRACKET)
		advanceX(X_SENTBRAK)
		advanceX(X_NAME)
		advanceX(X_OITALIC)
		if flags&QUOT != 0 {
			advanceX(X_CQUOTE)
		}
		if flags&BRAK != 0 {
			advanceX(X_CBRACKET)
		}
		if flags&ITAL != 0 {
			advanceX(X_CITALIC)
		}

		old_flags, old_lastopened = flags, lastopened
	}
}
