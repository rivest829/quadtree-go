package quadtree

const (
	MaxObjects int = 10
	MaxLevels  int = 4
)

type Rect struct {
	X      float64 // 左上角X
	Y      float64 // 左上角Y
	Width  float64
	Height float64
}

type Quadtree interface {
	Split()
	GetIndex(r Rect) []int
	Insert(r Rect)
	Retrieve(r Rect) []Rect
	Clear()
}

type quadtreeNode struct {
	Level   int
	Bounds  Rect
	Objects []Rect
	Nodes   []Quadtree
}

func NewQuadtree(bounds Rect, level int) Quadtree {
	return &quadtreeNode{
		Level:   level,
		Bounds:  bounds,
		Objects: make([]Rect, 0),
		Nodes:   make([]Quadtree, 0),
	}
}

func (q *quadtreeNode) Split() {
	nextLevel := q.Level + 1
	subWidth := q.Bounds.Width / 2
	subHeight := q.Bounds.Height / 2
	x := q.Bounds.X
	y := q.Bounds.Y

	q.Nodes = append(q.Nodes, NewQuadtree(Rect{x + subWidth, y, subWidth, subHeight}, nextLevel))
	q.Nodes = append(q.Nodes, NewQuadtree(Rect{x, y, subWidth, subHeight}, nextLevel))
	q.Nodes = append(q.Nodes, NewQuadtree(Rect{x, y + subHeight, subWidth, subHeight}, nextLevel))
	q.Nodes = append(q.Nodes, NewQuadtree(Rect{x + subWidth, y + subHeight, subWidth, subHeight}, nextLevel))
}

func (q *quadtreeNode) GetIndex(pRect Rect) []int {
	indexes := make([]int, 0)
	verticalMidpoint := q.Bounds.X + (q.Bounds.Width / 2)
	horizontalMidpoint := q.Bounds.Y + (q.Bounds.Height / 2)

	startIsNorth := pRect.Y < horizontalMidpoint
	startIsWest := pRect.X < verticalMidpoint
	endIsEast := pRect.X+q.Bounds.Width > verticalMidpoint
	endIsSouth := pRect.Y+q.Bounds.Height > horizontalMidpoint

	if startIsNorth && endIsEast {
		indexes = append(indexes, 0)
	}

	if startIsWest && startIsNorth {
		indexes = append(indexes, 1)
	}

	if startIsWest && endIsSouth {
		indexes = append(indexes, 2)
	}

	if endIsEast && endIsSouth {
		indexes = append(indexes, 3)
	}

	return indexes
}

func (q *quadtreeNode) Insert(pRect Rect) {
	if len(q.Nodes) != 0 {
		indexes := q.GetIndex(pRect)
		for _, index := range indexes {
			q.Nodes[index].Insert(pRect)
		}
		return
	}

	q.Objects = append(q.Objects, pRect)
	if len(q.Objects) > MaxObjects && q.Level < MaxLevels {

		if len(q.Nodes) == 0 {
			q.Split()
		}

		i := 0
		for i < len(q.Objects) {
			indexes := q.GetIndex(q.Objects[i])
			for _, index := range indexes {
				q.Nodes[index].Insert(q.Objects[i])
			}
			i++
		}
		q.Objects = make([]Rect, 0)
	}
}

func (q *quadtreeNode) Retrieve(pRect Rect) []Rect {
	indexes := q.GetIndex(pRect)
	returnObjects := q.Objects

	if len(q.Nodes) != 0 {
		for _, index := range indexes {
			returnObjects = append(returnObjects, q.Nodes[index].Retrieve(pRect)...)
		}
	}

	return returnObjects
}

func (q *quadtreeNode) Clear() {
	q.Objects = make([]Rect, 0)

	for i := 0; i < len(q.Nodes); i++ {
		if len(q.Nodes) != 0 {
			q.Nodes[i].Clear()
		}
	}

	q.Nodes = make([]Quadtree, 0)
}

func (q *quadtreeNode) PrintAll() string {
	if len(q.Objects) == 0 && len(q.Nodes) == 0 {
		return ""
	}
	log := ""
	for i := uint8(0); i < q.Level; i++ {
		log += "   "
	}
	log += "==="
	log += fmt.Sprintf(" %.3f:%.3f ", q.Bounds.Width, q.Bounds.Height)
	for _, obj := range q.Objects {
		log += obj.Owner().Id() + " / "
	}
	log += "\n"
	for _, node := range q.Nodes {
		log += node.PrintAll()
	}
	return log
}

