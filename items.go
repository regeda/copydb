package copydb

import (
	"container/list"

	"github.com/pkg/errors"
)

type items map[string]*item

func (ii items) item(id string, pool Pool) *item {
	i, ok := ii[id]
	if !ok {
		i = &item{
			Item: pool.New(),
			elem: &list.Element{Value: id},
		}
		ii[id] = i
	}
	return i
}

func (ii items) destroy(id string, pool Pool) bool {
	i, ok := ii[id]
	if !ok {
		return false
	}
	pool.Destroy(i.Item)
	delete(ii, id)
	return true
}

type item struct {
	Item

	elem    *list.Element
	version int64
}

func (i *item) init(version int64, data map[string]string) {
	i.version = version
	i.Remove()
	for k, v := range data {
		i.Set(k, []byte(v))
	}
}

func (i *item) apply(u *Update) error {
	if u.Version-i.version > 1 {
		return errors.Wrapf(ErrVersionConflict, "update version (%d) is greater db version (%d) for %s", u.Version, i.version, u.Id)
	}
	if u.Remove {
		i.Remove()
	} else {
		for _, f := range u.Set {
			i.Set(f.Name, f.Data)
		}
		for _, f := range u.Unset {
			i.Unset(f.Name)
		}
	}
	i.version = u.Version
	return nil
}
