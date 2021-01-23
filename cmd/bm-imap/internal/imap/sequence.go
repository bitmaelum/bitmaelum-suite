package imap

import (
	"math"
	"strconv"
	"strings"
)

type Sequence struct {
	Min int
	Max int
}

type SequenceSet struct {
	seq []Sequence
}

func (s *SequenceSet) InSet(i int) bool {
	for _, seq := range s.seq {
		if i >= seq.Min && i <= seq.Max {
			return true
		}
	}
	return false
}


func NewSequenceSet(s string) SequenceSet {
	parts := strings.Split(s, ",")

	set := SequenceSet{
		seq: []Sequence{},
	}

	for _, r := range parts {
		if strings.Contains(r, ":") {

			var (
				min int
				max int
			)

			parts := strings.Split(r, ":")
			min, _ = strconv.Atoi(parts[0])
			if parts[1] == "*" {
				max = math.MaxInt32
			} else {
				max, _ = strconv.Atoi(parts[0])
			}

			set.seq = append(set.seq, Sequence{
				Min: min,
				Max: max,
			})
		} else {
			uid, _ := strconv.Atoi(r)
			set.seq = append(set.seq, Sequence{
				Min: uid,
				Max: uid,
			})
		}
	}

	return set
}
