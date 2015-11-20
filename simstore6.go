package simstore

type Storage interface {
	Add(sig, docid uint64)
	Find(sig uint64) []uint64
	Finish()
}

type Store6 struct {
	Store
}

func New6(hashes int) *Store6 {
	var s Store6
	s.tables = make([]table, 49)

	if hashes != 0 {
		for i := range s.tables {
			s.tables[i] = make([]entry, 0, hashes)
		}
	}

	return &s
}

// Add inserts a signature and document id into the store
func (s *Store6) Add(sig uint64, docid uint64) {
	t := 0

	var p uint64

	for i := 0; i < 6; i++ {
		p = sig
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff80007fffffffff) | (sig & 0x007f800000000000 >> 8) | (sig & 0x00007f8000000000 << 8)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff807f807fffffff) | (sig & 0x007f800000000000 >> 16) | (sig & 0x0000007f80000000 << 16)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff807fff807fffff) | (sig & 0x007f800000000000 >> 24) | (sig & 0x000000007f800000 << 24)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff807fffff807fff) | (sig & 0x007f800000000000 >> 32) | (sig & 0x00000000007f8000 << 32)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff807fffffff807f) | (sig & 0x007f800000000000 >> 40) | (sig & 0x0000000000007f80 << 40)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		p = (sig & 0xff80ffffffffff80) | (sig & 0x007f000000000000 >> 48) | (sig & 0x000000000000007f << 48)
		s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
		t++
		sig = (sig << 9) | (sig >> (64 - 9))
	}

	p = sig
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc0003fffffffff) | (sig & 0x003fc00000000000 >> 8) | (sig & 0x00003fc000000000 << 8)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc03fc03fffffff) | (sig & 0x003fc00000000000 >> 16) | (sig & 0x0000003fc0000000 << 16)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc03fffc03fffff) | (sig & 0x003fc00000000000 >> 24) | (sig & 0x000000003fc00000 << 24)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc03fffffc03fff) | (sig & 0x003fc00000000000 >> 32) | (sig & 0x00000000003fc000 << 32)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc07fffffffc07f) | (sig & 0x003f800000000000 >> 40) | (sig & 0x0000000000003f80 << 40)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
	t++
	p = (sig & 0xffc07fffffffff80) | (sig & 0x003f800000000000 >> 47) | (sig & 0x000000000000007f << 47)
	s.tables[t] = append(s.tables[t], entry{hash: p, docid: docid})
}

func (*Store6) unshuffle(sig uint64, t int) uint64 {

	t7 := t % 7
	shift := 8 * uint64(t7)

	var m2 uint64

	if t < 42 {
		m2 = 0x007f800000000000

		if t7 == 6 {
			m2 = 0x007f000000000000
		}
	} else {
		m2 = 0x003fc00000000000

		if t7 >= 5 {
			m2 = 0x003f800000000000

			if t7 == 6 {
				shift--
			}
		}
	}

	m3 := uint64(m2 >> shift)
	m1 := ^uint64(0) &^ (m2 | m3)

	sig = (sig & m1) | (sig & m2 >> shift) | (sig & m3 << shift)
	sig = (sig >> (9 * (uint64(t) / 7))) | (sig << (64 - (9 * (uint64(t) / 7))))
	return sig
}

const mask6_9_8 = 0xffff800000000000
const mask6_9_7 = 0xffff000000000000
const mask6_10_8 = 0xffffc00000000000
const mask6_10_7 = 0xffff800000000000

// Find searches the store for all hashes hamming distance 6 or less from the
// query signature.  It returns the associated list of document ids.
func (s *Store6) Find(sig uint64) []uint64 {
	var ids []uint64

	// TODO(dgryski): search in parallel

	t := 0

	var p uint64

	for i := 0; i < 6; i++ {
		p = sig
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff80007fffffffff) | (sig & 0x007f800000000000 >> 8) | (sig & 0x00007f8000000000 << 8)
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff807f807fffffff) | (sig & 0x007f800000000000 >> 16) | (sig & 0x0000007f80000000 << 16)
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff807fff807fffff) | (sig & 0x007f800000000000 >> 24) | (sig & 0x000000007f800000 << 24)
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff807fffff807fff) | (sig & 0x007f800000000000 >> 32) | (sig & 0x00000000007f8000 << 32)
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff807fffffff807f) | (sig & 0x007f800000000000 >> 40) | (sig & 0x0000000000007f80 << 40)
		ids = append(ids, s.tables[t].find(p, mask6_9_8, 6)...)
		t++
		p = (sig & 0xff80ffffffffff80) | (sig & 0x007f000000000000 >> 48) | (sig & 0x000000000000007f << 48)
		ids = append(ids, s.tables[t].find(p, mask6_9_7, 6)...)
		t++
		sig = (sig << 9) | (sig >> (64 - 9))
	}

	p = sig
	ids = append(ids, s.tables[t].find(p, mask6_10_8, 6)...)
	t++
	p = (sig & 0xffc0003fffffffff) | (sig & 0x003fc00000000000 >> 8) | (sig & 0x00003fc000000000 << 8)
	ids = append(ids, s.tables[t].find(p, mask6_10_8, 6)...)
	t++
	p = (sig & 0xffc03fc03fffffff) | (sig & 0x003fc00000000000 >> 16) | (sig & 0x0000003fc0000000 << 16)
	ids = append(ids, s.tables[t].find(p, mask6_10_8, 6)...)
	t++
	p = (sig & 0xffc03fffc03fffff) | (sig & 0x003fc00000000000 >> 24) | (sig & 0x000000003fc00000 << 24)
	ids = append(ids, s.tables[t].find(p, mask6_10_8, 6)...)
	t++
	p = (sig & 0xffc03fffffc03fff) | (sig & 0x003fc00000000000 >> 32) | (sig & 0x00000000003fc000 << 32)
	ids = append(ids, s.tables[t].find(p, mask6_10_8, 6)...)
	t++
	p = (sig & 0xffc07fffffffc07f) | (sig & 0x003f800000000000 >> 40) | (sig & 0x0000000000003f80 << 40)
	ids = append(ids, s.tables[t].find(p, mask6_10_7, 6)...)
	t++
	p = (sig & 0xffc07fffffffff80) | (sig & 0x003f800000000000 >> 47) | (sig & 0x000000000000007f << 47)
	ids = append(ids, s.tables[t].find(p, mask6_10_7, 6)...)
	t++

	return unique(ids)
}
