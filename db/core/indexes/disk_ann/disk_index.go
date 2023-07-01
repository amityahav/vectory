package disk_ann

import "Vectory/db/core/indexes/utils"

type DiskIndex struct {
	dal *dal
}

func (di *DiskIndex) search(q []float32, k int, listSize int, onlySearch bool) ([]utils.Element, []utils.Element) {
	return nil, nil
}
