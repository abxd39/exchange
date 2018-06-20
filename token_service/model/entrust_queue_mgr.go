package model

type EntrustQueneMgr struct {
	dataMgr map[string]*EntrustQuene
}

func (s *EntrustQueneMgr) GetQuene(id string) (d *EntrustQuene, ok bool) {
	d, ok = s.dataMgr[id]
	if !ok {
		return
	}
	return
}

