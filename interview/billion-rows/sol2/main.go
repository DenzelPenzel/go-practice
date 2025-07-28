package sol2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"slices"
	"strings"
	"sync"
	"unsafe"

	"github.com/DenzelPenzel/go-leetcode/interview/billion-rows/utils"
)

const (
	// offset64  and prime64 are taken from fnv.go
	offset64   = 14695981039346656037
	prime64    = 1099511628211
	bucketSize = 1 << 25
)

const (
	shift1 = 8 * 1
	shift2 = 8 * 2
	shift3 = 8 * 3
	shift4 = 8 * 4

	charMask0 = uint64(255)
	charMask1 = uint64(255) << shift1
	charMask2 = uint64(255) << shift2
	charMask3 = uint64(255) << shift3
	charMask4 = uint64(255) << shift4

	dot1 = uint64('.') << 8
	dot2 = uint64('.') << 16
)

func Run(fileName string) string {
	workers := runtime.NumCPU()
	data := utils.Mmap(fileName)
	chunkSize := len(data) / workers

	if chunkSize == 0 {
		chunkSize = len(data)
	}

	chunks := make([]int, 0, chunkSize)
	s := 0

	for s < len(data) {
		e := s + chunkSize
		if e >= len(data) {
			chunks = append(chunks, len(data))
			break
		}
		nl := bytes.IndexByte(data[e:], '\n')
		if nl < 0 {
			chunks = append(chunks, len(data))
			break
		}
		e += nl + 1
		chunks = append(chunks, e)
		s = e
	}

	var wg sync.WaitGroup
	groups := make([]*Bucket, len(chunks))
	start := 0

	for i, end := range chunks {
		wg.Add(1)

		go func(workerId int, start, end uint64) {
			defer wg.Done()
			var b Bucket

			for start < end {
				w := binary.LittleEndian.Uint64(data[start : start+8])
				// w := *(*uint64)(unsafe.Pointer(&data[start]))

				// fmt.Printf("workerId: %d, partition: %v -> %v, semi: %v\n", workerId, start, end, findSemicolon(uint64(firstBytes)))

				var city []byte

				// check the presence of a semicolon within the initial 8 bytes
				if idx := findSemicolon(w); idx >= 0 {
					city = data[start : start+uint64(idx)]
					start += uint64(idx) + 1
				} else {
					// presence of the a semicolon within the first 8 bytes not found
					// move the pointer and check the next 8 bytes
					for i := start + 8; i < end; i += 8 {
						u := binary.LittleEndian.Uint64(data[i : i+8])
						if idx = findSemicolon(u); idx >= 0 {
							city = data[start : i+uint64(idx)]
							start = i + uint64(idx) + 1
							break
						}
					}
				}

				h := createHash(w, len(city))
				w = binary.LittleEndian.Uint64(data[start : start+8])
				temp, pos := parseNumber(w)

				node := b.Insert(h, city)
				node.min = min(node.min, temp)
				node.max = max(node.max, temp)
				node.sum += int64(temp)
				node.count++

				start += pos
			}

			groups[workerId] = &b
		}(i, uint64(start), uint64(end))

		start = end
	}

	wg.Wait()

	total := 0
	for _, b := range groups {
		total += len(b.Keys())
	}

	cities := make([]string, total)

	o := 0
	for _, b := range groups {
		ks := b.Keys()
		copy(cities[o:], ks)
		o += len(ks)
	}

	slices.Sort(cities)
	cities = slices.Compact(cities)

	var sb strings.Builder
	sb.WriteByte('{')

	for i, city := range cities {
		n := Node{
			key: city,
			min: math.MaxInt16,
			max: math.MinInt16,
		}

		u := *(*uint64)(unsafe.Pointer(unsafe.StringData(city)))
		key := createHash(u, len(city))

		for k := range groups {
			if item := groups[k].Find(key, city); item != nil {
				n.max = max(n.max, item.max)
				n.min = min(n.min, item.min)
				n.sum += item.sum
				n.count += item.count
			}
		}

		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f", city,
			utils.Round(float64(n.min)/10.0),
			utils.Round(float64(n.sum)/10.0/float64(n.count)),
			utils.Round(float64(n.max)/10.0)))
	}

	sb.WriteString("}\n")
	return sb.String()
}

func findSemicolon(word uint64) int {
	// See "determine if a word has a byte equal to n", from "Bit Twiddling
	// Hacks",
	// https://graphics.stanford.edu/~seander/bithacks.html.
	maskedInput := word ^ 0x3B3B3B3B3B3B3B3B
	maskedInput = (maskedInput - 0x0101010101010101) & ^maskedInput & 0x8080808080808080
	if maskedInput == 0 {
		return -1
	}
	// Divide by 8.
	return bits.TrailingZeros64(maskedInput) >> 3
}

type Hash uint64

// Index Compute a simple hash based on FNV
func (k Hash) Index() uint64 {
	var h uint64
	h = offset64
	h ^= uint64(k)
	h *= prime64
	// Compute the modulus of the hash with the table size.
	return h & (bucketSize - 1)
}

func createHash(u uint64, p int) Hash {
	if p >= 8 {
		return Hash(u)
	}
	var m uint64
	m = 1<<(p<<3) - 1
	return Hash(u & m)
}

type Node struct {
	key   string
	hash  Hash
	next  *Node
	sum   int64
	count int64
	min   int16
	max   int16
}

type Bucket struct {
	keys   []string
	bucket [bucketSize]*Node
}

func (b *Bucket) Keys() []string {
	return b.keys
}

func (b *Bucket) Find(h Hash, key string) *Node {
	cb := b.bucket[h.Index()]
	for cb != nil {
		if h == cb.hash && (len(key) <= 8 || key == cb.key) {
			return cb
		}
		cb = cb.next
	}
	return nil
}

func (b *Bucket) Insert(h Hash, key []byte) *Node {
	idx := h.Index()
	cb := b.bucket[idx]
	prev := cb
	for cb != nil {
		if h == cb.hash && (len(key) <= 8 || string(key) == cb.key) {
			return cb
		}
		prev = cb
		cb = cb.next
	}
	node := &Node{
		key:  string(key),
		hash: h,
		min:  math.MaxInt16,
		max:  math.MinInt16,
	}

	if prev != nil {
		prev.next = node
	} else {
		b.bucket[idx] = node
	}
	b.keys = append(b.keys, node.key)
	return node
}

func parseNumber(u uint64) (_ int16, advance uint64) {
	// Gather the key and fetch the corresponding record. We can do this without
	// scanning the line because there are only four possible sequences for
	// the temperature. The valid formats are:
	//
	//     0.0
	//    00.0
	//    -0.0
	//   -00.0
	//
	// We use the limited locations of semicolons and minus characters to avoid
	// conditional expressions and loops.
	switch {
	case (u & charMask1) == dot1:
		// Case: "0.0".
		ones := (u&charMask0 - '0') * 10
		tenths := u&charMask2>>shift2 - '0'
		// Advance past the newline, which is the fourth character.
		return int16(ones + tenths), 4

	case (u & charMask2) == dot2:
		// Case: "00.0" and "-0.0".
		v0 := u & charMask0
		tens := (v0 - '0') * 100
		ones := (u&charMask1>>shift1 - '0') * 10
		tenths := u&charMask3>>shift3 - '0'

		// neg is 1 if the first character is the minus charater, and 0
		// otherwise.
		//
		// NOTE: The Go compiler eliminates jumps when using this form of
		// conditional.
		var neg uint64
		if v0 == '-' {
			neg = 1
		}

		// Clear tens if there was a minus character in that position.
		tens = tens &^ -neg

		// Add the ones and tenths digits.
		temp := ones + tenths

		// Add the tens digit.
		temp += tens

		// Negate the value if we found a minus character.
		//
		// See "conditionally negate a value without branching", from "Bit
		// Twiddling Hacks",
		// https://graphics.stanford.edu/~seander/bithacks.html.
		temp = (temp ^ -neg) + neg

		// Advance past the newline, which is the fifth character.
		return int16(temp), 5
	default:
		// Case: "-00.0".
		tens := (u&charMask1>>shift1 - '0') * 100
		ones := (u&charMask2>>shift2 - '0') * 10
		tenths := u&charMask4>>shift4 - '0'

		t := int16(tens + ones + tenths)
		t *= -1
		// Advance past the newline, which is the sixth character.
		return t, 6
	}
}
