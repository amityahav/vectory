package disk_ann

import (
	"math"
	"os"
	"sort"
)

const (
	pageSize        = 0x1000
	metadataPageNum = 0
)

type page struct {
	pageNum int
	data    []byte
}

type dal struct {
	file     *os.File
	nextPage int
}

func newDal(path string) (*dal, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	d := dal{
		file:     file,
		nextPage: 1, // 0 is reserved for metadata
	}

	return &d, nil
}

func (d *dal) newMemoryEmptyPage() *page {
	p := page{
		data: make([]byte, pageSize),
	}

	return &p
}

func (d *dal) readPage(pageNum int) (*page, error) {
	p := d.newMemoryEmptyPage()

	offset := int64(pageNum * pageSize)

	_, err := d.file.ReadAt(p.data, offset)
	if err != nil {
		return nil, err
	}

	p.pageNum = pageNum

	return p, nil
}

func (d *dal) writePage(p *page) error {
	offset := int64(p.pageNum * pageSize)

	_, err := d.file.WriteAt(p.data, offset)
	if err != nil {
		return err
	}

	return nil
}

func (d *dal) getNextPage() int {
	curr := d.nextPage
	d.nextPage++

	return curr
}

/*
Index disk layout:
- Each page is of pageSize size.
- First page contains size's metadata: vectors dimensions, graph's size, graph's max degree, firstId, index size, s id.
- All other pages contain vertices, each vertex contains: objId, vector, neighbors.
- The number of vertices in each page is calculated by: floor(pageSize / (4 * (1 + dim + degree))).
 ------------------------------------------------
| metadata	| vertices	| vertices	| vertices	|
|			|			|			|			|
 ------------------------------------------------
*/
func (d *dal) writeIndex(mi *MemoryIndex) error {
	// writing metadata page
	metadataPage := d.newMemoryEmptyPage()
	metadataPage.pageNum = metadataPageNum
	mi.serializeMetadata(metadataPage.data)

	err := d.writePage(metadataPage)
	if err != nil {
		return err
	}

	// writing vertices pages
	sortedIds := make([]uint32, 0, len(mi.graph.vertices))

	for id := range mi.graph.vertices {
		sortedIds = append(sortedIds, id)
	}

	sort.Slice(sortedIds, func(i, j int) bool {
		return sortedIds[i] < sortedIds[j]
	})

	numOfVerticesInPage := math.Floor(float64(pageSize / (4 * (1 + mi.dim + mi.maxDegree))))
	numOfFullPages := len(sortedIds) / int(numOfVerticesInPage)
	remainder := len(sortedIds) % int(numOfVerticesInPage)

	var offset int

	for i := 0; i < numOfFullPages; i++ {
		p := d.newMemoryEmptyPage()
		p.pageNum = d.getNextPage()

		mi.graph.serializeVertices(p.data, sortedIds[offset:offset+int(numOfVerticesInPage)], mi.maxDegree)

		err = d.writePage(p)
		if err != nil {
			return err
		}

		offset += int(numOfVerticesInPage)
	}

	if remainder != 0 {
		p := d.newMemoryEmptyPage()
		p.pageNum = d.getNextPage()

		mi.graph.serializeVertices(p.data, sortedIds[offset:offset+remainder], mi.maxDegree)

		err = d.writePage(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *dal) readIndex() (*MemoryIndex, error) {
	mi := MemoryIndex{graph: newGraph()}

	// reading metadata page
	p, err := d.readPage(metadataPageNum)
	if err != nil {
		return nil, err
	}

	mi.deserializeMetadata(p.data)

	// reading vertices pages
	numOfVerticesInPage := uint32(math.Floor(float64(pageSize / (4 * (1 + mi.dim + mi.maxDegree)))))
	numOfFullPages := mi.Size() / numOfVerticesInPage
	remainder := mi.Size() % numOfVerticesInPage

	currId := mi.firstId
	for i := 1; i <= int(numOfFullPages); i++ {
		p, err = d.readPage(i)
		if err != nil {
			return nil, err
		}

		mi.graph.deserializeVertices(p.data, mi.dim, mi.maxDegree, numOfVerticesInPage, currId)
		currId += numOfVerticesInPage
	}

	if remainder != 0 {
		p, err = d.readPage(int(numOfFullPages) + 1)
		if err != nil {
			return nil, err
		}

		mi.graph.deserializeVertices(p.data, mi.dim, mi.maxDegree, remainder, currId)
		currId += remainder
	}

	return &mi, nil
}
