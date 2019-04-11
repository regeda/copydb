package spatial

type cell map[*item]struct{}

func (c cell) add(it *item) {
	c[it] = struct{}{}
}

func (c cell) remove(it *item) {
	delete(c, it)
}
