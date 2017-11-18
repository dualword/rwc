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
	"math/rand"
	"sync"
	"time"

	"github.com/opennota/vose"
)

type lockedSource struct {
	m   sync.Mutex
	src rand.Source64
}

func (s *lockedSource) Int63() int64 {
	s.m.Lock()
	v := s.src.Int63()
	s.m.Unlock()
	return v
}

func (s *lockedSource) Uint64() uint64 {
	s.m.Lock()
	v := s.src.Uint64()
	s.m.Unlock()
	return v
}

func (s *lockedSource) Seed(seed int64) {
	s.m.Lock()
	s.src.Seed(seed)
	s.m.Unlock()
}

var (
	_seed = time.Now().UnixNano()
	_src  = rand.NewSource(_seed).(rand.Source64)
	_rand = rand.New(&lockedSource{src: _src})

	// https://ru.wikipedia.org/wiki/Частотность
	// буква 	употреблений
	// а 		40487008
	// б 		8051767
	// в 		22930719
	// г 		8564640
	// д 		15052118
	// е 		42691213
	// ё 		184928
	// ж 		4746916
	// з 		8329904
	// и 		37153142
	// й 		6106262
	// к 		17653469
	// л 		22230174
	// м 		16203060
	// н 		33838881
	// о 		55414481
	// п 		14201572
	// р 		23916825
	// с 		27627040
	// т 		31620970
	// у 		13245712
	// ф 		1335747
	// х 		4904176
	// ц 		2438807
	// ч 		7300193
	// ш 		3678738
	// щ 		1822476
	// ъ 		185452
	// ы 		9595941
	// ь 		8784613
	// э 		1610107
	// ю 		3220715
	// я 		10139085

	randA = vose.New(
		_rand,
		[]float64{
			40487008, 8051767, 22930719, 8564640, 15052118, 42691213 + 184928,
			4746916, 8329904, 37153142, 6106262, 17653469, 22230174, 16203060,
			33838881, 55414481, 14201572, 23916825, 27627040, 31620970, 13245712,
			1335747, 4904176, 2438807, 7300193, 3678738, 1822476, 185452, 9595941,
			8784613, 1610107, 3220715, 10139085,
		},
	)

	randV = vose.New(
		_rand,
		[]float64{
			40487008, 42691213 + 184928, 37153142, 55414481, 13245712, 9595941,
			1610107, 3220715, 10139085,
		},
	)

	randC = vose.New(
		_rand,
		[]float64{
			8051767, 22930719, 8564640, 15052118, 4746916, 8329904, 6106262, 17653469,
			22230174, 16203060, 33838881, 14201572, 23916825, 27627040, 31620970,
			1335747, 4904176, 2438807, 7300193, 3678738, 1822476,
		},
	)
)
