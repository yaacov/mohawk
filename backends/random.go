package backend

import (
	"fmt"
	"math/rand"
)

func FilterItems(vs []Item, f func(Item) bool) []Item {
	vsf := make([]Item, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

type Random struct {
	Items []Item
}

func (r *Random) Open() {
	r.Items = []Item{}

	seeds := []map[string]string{
		map[string]string{"type": "node", "group_id": "cpu/usage_rate", "units": "cpu", "issue": "42"},
		map[string]string{"type": "node", "group_id": "memory/usage_rate", "units": "byte"},
		map[string]string{"type": "node", "group_id": "cpu/usage_rate", "units": "cpu", "issue": "42"},
		map[string]string{"type": "node", "group_id": "memory/usage_rate", "units": "byte"},
		map[string]string{"type": "node", "group_id": "cpu/limit", "units": "cpu", "issue": "442"},
		map[string]string{"type": "node", "group_id": "memory/limit", "units": "byte", "issue": "442"},
		map[string]string{"type": "node", "group_id": "filesystem/usage_rate", "units": "byte"},
	}

	for i := 0; i < 120; i++ {
		seed := seeds[rand.Intn(len(seeds))]
		tags := map[string]string{
			"type":     seed["type"],
			"group_id": seed["group_id"],
			"units":    seed["units"],
			"issue":    seed["issue"],
			"hostname": fmt.Sprintf("example.%03d.com", i/4),
		}
		r.Items = append(r.Items, Item{
			Id:   fmt.Sprintf("hello_kitty_%3d", i),
			Tags: tags})
	}
}

func (r Random) GetItemList(tags map[string]string) []Item {
	res := r.Items

	if len(tags) > 0 {
		for key, value := range tags {
			res = FilterItems(res, func(i Item) bool { return i.Tags[key] == value })
		}
	}

	return res
}

func (r Random) GetRawData(id string, end int64, start int64, limit int64, order string) []DataItem {
	res := []DataItem{}
	delta := int64(5 * 60 * 1000)

	for i := limit; i > 0; i-- {
		res = append(res, DataItem{
			Timestamp: start - i*delta,
			Value:     float64(50 + rand.Intn(50))})
	}

	return res
}