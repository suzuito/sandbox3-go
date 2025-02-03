package article

type Articles []*Article

func (t Articles) GroupByID() map[ID]*Article {
	ret := make(map[ID]*Article, len(t))
	for _, a := range t {
		ret[a.ID] = a
	}
	return ret
}
