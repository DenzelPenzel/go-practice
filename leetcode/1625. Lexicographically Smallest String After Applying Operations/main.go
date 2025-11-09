package main

func findLexSmallestString(s string, a int, b int) string {
	queue := make([]string, 0)
	seen := make(map[string]bool)
	queue = append(queue, s)
	seen[s] = true
	res := s

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v < res {
			res = v
		}
		sBytes := []byte(v)
		for i := 1; i < len(sBytes); i += 2 {
			digit := (sBytes[i] - '0') + byte(a)
			sBytes[i] = (digit % 10) + '0'
		}

		add_s := string(sBytes)

		if !seen[add_s] {
			seen[add_s] = true
			queue = append(queue, add_s)
		}

		rotate_s := v[len(v)-b:] + v[:len(v)-b]
		if !seen[rotate_s] {
			seen[rotate_s] = true
			queue = append(queue, rotate_s)
		}
	}

	return res
}
