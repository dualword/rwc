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
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// LoadFrom loads binary representation of Constructor from r.
func (c *Constructor) LoadFrom(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, c.ng4[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, c.ng3[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, c.ng3beg[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, c.ng3end[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, c.ng2[:]); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &c.ng1); err != nil {
		return err
	}
	return nil
}

// WriteTo writes binary representation of Constructor to w.
func (c *Constructor) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.LittleEndian, c); err != nil {
		return 0, err
	}
	return int64(len(c.ng4) + len(c.ng3) + len(c.ng3beg) + len(c.ng3end) + len(c.ng2) + 4), nil
}

// LoadFromRWC loads binary representation of Constructor from an .RWC file.
func (c *Constructor) LoadFromRWC(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	const signature = "RTC 3.1 DataFile (C) 1992 OMS"
	buf := make([]byte, 143492)
	if _, err := io.ReadFull(r, buf[:len(signature)]); err != nil {
		return err
	}
	if string(buf[:len(signature)]) != signature {
		return errors.New("invalid file signature")
	}
	buf = buf[:0]

	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if b != 0xff {
			buf = append(buf, b)
		} else {
			b, err := r.ReadByte()
			if err != nil {
				return err
			}
			if b == 0xff {
				buf = append(buf, 0xff)
			} else {
				buf = append(buf, 0)
				b++
				for b > 0 {
					buf = append(buf, 0)
					b--
				}
			}
		}
	}

	return c.LoadFrom(bytes.NewReader(buf))
}
