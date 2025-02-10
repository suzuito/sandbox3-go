package article

type Articles []*Article

func (t Articles) GroupByID() map[ID]Articles {
	ret := make(map[ID]Articles, len(t))
	for _, a := range t {
		ret[a.ID] = append(ret[a.ID], a)
	}
	return ret
}
