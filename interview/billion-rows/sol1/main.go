package sol1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/DenzelPenzel/go-leetcode/interview/billion-rows/utils"
)

type node struct {
	min, max, sum, count int64
}

func Run(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open the file %v", err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	mapping := make(map[string]*node)

	for s.Scan() {
		line := s.Text()
		data := strings.Split(line, ";")
		key := data[0]
		val := convertStringToInt64(data[1])

		if item, ok := mapping[key]; ok {
			if val > item.max {
				item.max = val
			}
			if val < item.min {
				item.min = val
			}
			item.sum += val
			item.count++
		} else {
			mapping[key] = &node{min: val, max: val, sum: val, count: 1}
		}
	}

	cities := make([]string, 0, len(mapping))
	for city := range mapping {
		cities = append(cities, city)
	}

	sort.Strings(cities)

	var stringsBuilder strings.Builder

	stringsBuilder.WriteString(fmt.Sprintf("{"))
	for i, city := range cities {
		if i > 0 {
			stringsBuilder.WriteString(", ")
		}
		m := mapping[city]
		stringsBuilder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f", city,
			utils.Round(float64(m.min)/10.0),
			utils.Round(float64(m.sum)/10.0/float64(m.count)),
			utils.Round(float64(m.max)/10.0)))
	}
	stringsBuilder.WriteString(fmt.Sprintf("}\n"))

	return stringsBuilder.String()
}

func convertStringToInt64(s string) int64 {
	s = s[:len(s)-2] + s[len(s)-1:]
	res, _ := strconv.ParseInt(s, 10, 64)
	return res
}
