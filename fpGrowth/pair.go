package fpGrowth

import "sort"

func rank(frequencies map[int]uint) PairList {
	pl := make(PairList, len(frequencies))
	i := 0
	for k, v := range frequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

type Pair struct {
	Item      int
	Frequency uint
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Frequency > p[j].Frequency }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
