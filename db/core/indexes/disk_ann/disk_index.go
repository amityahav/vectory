package disk_ann

import "Vectory/db/core/indexes/utils"

type DiskIndex struct {
}

func (di *DiskIndex) Search(q []float32, k int, listSize int, onlySearch bool) ([]utils.Element, []utils.Element) {
	return nil, nil
}
