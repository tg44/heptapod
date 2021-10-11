package utils

type List struct {
	Data string
	Next *List
}

func (p *List) AddAsHead(s string) *List {
	return &List{s, p}
}

func (p *List) Prepend(s []string) *List {
	r := p
	for _, d := range s {
		r = r.AddAsHead(d)
	}
	return r
}

func (p *List) Size() int {
	if p == nil {
		return 0
	}
	i := 1
	k := p
	for k.Next != nil {
		i += 1
		k = k.Next
	}
	return i
}

func (p *List) ToArray() []string {
	if p == nil {
		return []string{}
	}
	i := []string{}
	k := p
	for k.Next != nil {
		i = append(i, k.Data)
		k = k.Next
	}
	i = append(i, k.Data)
	return i
}

func (p *List) Contains(s string) bool {
	if p == nil {
		return false
	}
	k := p
	for k.Next != nil {
		if k.Data == s {
			return true
		}
		k = k.Next
	}
	if k.Data == s {
		return true
	} else {
		return false
	}
}
