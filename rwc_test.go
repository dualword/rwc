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
	"regexp"
	"testing"
	"unicode/utf8"
)

func TestWord(t *testing.T) {
	for i := 1; i <= 20; i++ {
		w := Word(i)
		if utf8.RuneCountInString(w) != i {
			t.Errorf("want a word of length %d, got %q", i, w)
		}
		for _, r := range w {
			if r < 'а' || r > 'я' {
				t.Errorf("want the word contain only letters from '%c' to '%c', got %q", 'а', 'я', w)
				break
			}
		}
	}
}

func TestWordMask(t *testing.T) {
	rx := regexp.MustCompile("^[бвгджзйклмнпрстфхцчшщ][аеиоуыэюя][а-я]{3}ый$")
	for i := 1; i <= 20; i++ {
		const mask = "CV...ый"
		w := WordMask(mask)
		if !rx.MatchString(w) {
			t.Errorf("want a word matching the mask %q, got %q", mask, w)
		}
	}
}
