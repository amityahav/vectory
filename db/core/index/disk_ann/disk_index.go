package disk_ann

import "Vectory/db/core/index/utils"

type DiskIndex struct {
	dal *dal
}

func (di *DiskIndex) search(q []float32, k int, listSize int, onlySearch bool) ([]utils.Element, []utils.Element) {
	return nil, nil
}
