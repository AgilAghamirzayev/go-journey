package main

import (
	"fmt"
)

func lengthOfLongestSubstring(s string) int {
	lastSeen := make(map[rune]int)
	left := 0
	maxLen := 0

	for right, ch := range []rune(s) {
		// If character seen, move left pointer
		if idx, ok := lastSeen[ch]; ok && idx >= left {
			left = idx + 1
		}
		lastSeen[ch] = right

		// Update max length
		if right-left+1 > maxLen {
			maxLen = right - left + 1
		}
	}

	return maxLen
}

func main() {
	fmt.Println(lengthOfLongestSubstring("abcabcbb")) // Output: 3
	fmt.Println(lengthOfLongestSubstring("bbbbb"))    // Output: 1
	fmt.Println(lengthOfLongestSubstring("pwwkew"))   // Output: 3
}
