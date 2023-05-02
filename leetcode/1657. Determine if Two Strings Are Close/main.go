/*
Two strings are considered close if you can attain one from the other using the following operations:

Operation 1: Swap any two existing characters.
For example, abcde -> aecdb
Operation 2: Transform every occurrence of one existing character into another existing character,
and do the same with the other character.
For example, aacabb -> bbcbaa (all a's turn into b's, and all b's turn into a's)
You can use the operations on either string as many times as necessary.

Given two strings, word1 and word2, return true if word1 and word2 are close, and false otherwise.

Example 1:

	Input: word1 = "abc", word2 = "bca"
	Output: true
	Explanation: You can attain word2 from word1 in 2 operations.
	Apply Operation 1: "abc" -> "acb"
	Apply Operation 1: "acb" -> "bca"

Example 2:

	Input: word1 = "a", word2 = "aa"
	Output: false
	Explanation: It is impossible to attain word2 from word1, or vice versa, in any number of operations.

Example 3:

	Input: word1 = "cabbba", word2 = "abbccc"
	Output: true
	Explanation: You can attain word2 from word1 in 3 operations.
	Apply Operation 1: "cabbba" -> "caabbb"
	Apply Operation 2: "caabbb" -> "baaccc"
	Apply Operation 2: "baaccc" -> "abbccc"

Constraints:

	1 <= word1.length, word2.length <= 105
	word1 and word2 contain only lowercase English letters.
*/
package main

func closeStrings(word1 string, word2 string) bool {
	// Create maps to count character occurrences for both words
	m1 := make(map[rune]int)
	m2 := make(map[rune]int)

	// If lengths of both words are different, return false
	if len(word1) != len(word2) {
		return false
	}

	// Count character occurrences for both words
	for _, char := range word1 {
		m1[char] += 1
	}
	for _, char := range word2 {
		m2[char] += 1
	}

	// Check if all characters in word2 are present in word1
	for k := range m2 {
		if _, ok := m1[k]; !ok {
			return false
		}
	}

	// Create maps to count the frequency of character counts in both words
	cntOfCnt1 := make(map[int]int)
	cntOfCnt2 := make(map[int]int)

	for _, v := range m1 {
		cntOfCnt1[v] += 1
	}
	for _, v := range m2 {
		cntOfCnt2[v] += 1
	}

	// Check if the frequency of character counts are the same for both words
	for k := range cntOfCnt1 {
		if cntOfCnt1[k] != cntOfCnt2[k] {
			return false
		}
	}

	return true
}
