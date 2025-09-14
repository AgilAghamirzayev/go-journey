package main

import (
	"sort"
	"testing"
)

/*
Problem:
- Given a list of unsorted, independent meetings, returns a list of a merged
  one.

Example:
- Input: []meeting{{1, 2}, {2, 3}, {4, 5}}
  Output: []meeting{{1, 3}, {4, 5}}
- Input: []meeting{{1, 5}, {2, 3}}
  Output: []meeting{{1, 5}}

Approach:
- Sort the list in ascending order so that meetings that might need to be merged are next to each other.
- Can merge two meetings together if the first one's end time is greater or equal than the second one's start time.
*/

func TestMergeMeetings(t *testing.T) {
	tests := []struct {
		in       []meeting
		expected []meeting
	}{
		{[]meeting{}, []meeting{}},
		{[]meeting{{1, 2}}, []meeting{{1, 2}}},
		{[]meeting{{1, 2}, {2, 3}}, []meeting{{1, 3}}},
		{[]meeting{{1, 5}, {2, 3}}, []meeting{{1, 5}}},
		{[]meeting{{1, 2}, {4, 5}}, []meeting{{1, 2}, {4, 5}}},
		{[]meeting{{1, 5}, {2, 3}, {4, 5}}, []meeting{{1, 5}}},
		{[]meeting{{1, 2}, {2, 3}, {4, 5}}, []meeting{{1, 3}, {4, 5}}},
		{[]meeting{{1, 6}, {2, 3}, {4, 5}}, []meeting{{1, 6}}},
		{[]meeting{{4, 5}, {2, 3}, {1, 6}}, []meeting{{1, 6}}},
	}

	for _, tt := range tests {
		result := mergeMeetings(tt.in)
		Equal(t, tt.expected, result)
	}
}

type meeting struct {
	start int
	end   int
}

func mergeMeetings(meetings []meeting) []meeting {

	sort.Slice(meetings, func(i, j int) bool {
		return meetings[i].start < meetings[j].start
	})

	var out []meeting

	for i := range meetings {
		if i == 0 {
			out = append(out, meetings[i])
			continue
		}

		if out[len(out)-1].end >= meetings[i].start {
			out[len(out)-1].end = Max(meetings[i].end, out[len(out)-1].end)
		} else {
			out = append(out, meetings[i])
		}
	}

	return out

}
