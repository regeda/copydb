package copydb

import (
	"container/list"

	"github.com/pkg/errors"
	"github.com/regeda/copydb/internal/model"
)

type items map[string]*item

func (ii items) empty() bool {
	return len(ii) == 0
}

func (ii items) item(id string, pool Pool, lru *list.List) *item {
	i, ok := ii[id]
	if !ok {
		i = &item{
			Item: pool.Get(),
			elem: lru.PushBack(id),
		}
		ii[id] = i
	} else {
		lru.MoveToBack(i.elem)
	}
	return i
}

func (ii items) evict(deadline int64, pool Pool, lru *list.List) bool {
	elem := lru.Front()
	if elem == nil {
		return false
	}
	id := elem.Value.(string)
	i := ii[id]
	if i.unix > deadline {
		return false
	}
	lru.Remove(i.elem)
	pool.Put(i.Item)
	delete(ii, id)
	return true
}

type item struct {
	Item

	elem *list.Element

	unix    int64
	version int64
}

func (i *item) init(id string, unix, version int64, data map[string]string) {
	i.Remove()
	for k, v := range data {
		i.Set(k, []byte(v))
	}
	i.Set(keyID, []byte(id))
	i.unix = unix
	i.version = version
}

func (i *item) apply(u *model.Update) error {
	if u.Version-i.version > 1 {
		return errors.Wrapf(ErrVersionConflict, "update version (%d) is greater db version (%d) for %s", u.Version, i.version, u.ID)
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
		i.Set(keyID, []byte(u.ID))
	}
	i.unix = u.Unix
	i.version = u.Version
	return nil
}
