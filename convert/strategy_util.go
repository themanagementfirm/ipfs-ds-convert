package convert

import (

	ds "gx/ipfs/QmVSase1JP7cq9QkPT46oNwdp9pT6kBkG3oqS14y3QcZjG/go-datastore"
	"sort"
)

type Spec map[string]interface{}

func (s *Spec) dsType() (string, bool) {
	return s.str("type")
}

func (s *Spec) str(key string) (string, bool) {
	t, ok := (*s)[key]
	if !ok {
		return "", false
	}
	ts, ok := t.(string)
	return ts, ok
}

func (s *Spec) Id() string {
	return DatastoreSpec(*s)
}

type SimpleMount struct {
	prefix ds.Key
	diskId string

	spec Spec
}

type SimpleMounts []SimpleMount

func (m *SimpleMounts) hasPrefixed(searched SimpleMount) int {
	for i, mnt := range *m {
		if mnt.prefix.Equal(searched.prefix) {
			return i
		}
	}

	return -1
}

func (m *SimpleMounts) hasMatching(searched SimpleMount) bool {
	i := m.hasPrefixed(searched)

	if i != -1 {
		return (*m)[i].diskId == searched.diskId
	}

	return false
}

func (m *SimpleMounts) filter(filter SimpleMounts) SimpleMounts {
	out := make([]SimpleMount, 0, len(*m))

	for _, mount := range *m {
		if !filter.hasMatching(mount) {
			out = append(out, mount)
		}
	}

	return out
}

func (m *SimpleMounts) sort() {
	sort.Slice(*m, func(i, j int) bool { return (*m)[i].prefix.String() > (*m)[j].prefix.String() })
}

func (m *SimpleMounts) spec() Spec {
	mounts := make([]interface{}, 0, len(*m))

	for _, mount := range *m {
		var spec map[string]interface{} = mount.spec
		mounts = append(mounts, spec)
	}

	return map[string]interface{}{
		"type": "mount",
		"mounts": mounts,
	}
}