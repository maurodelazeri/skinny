package skinny

// Single point
type Point struct {
	Timestamp uint64
	Value     map[string]interface{}
	Next      *Point
}

// Core linked list, main storage array.
type Metric struct {
	newest, oldest           *Point
	length                   uint64
	Capacity                 uint64
	index                    map[uint64]*Point // coarseSearch index
	indexnewest, indexoldest uint64
	Indexinterval            uint64 // in seconds > 0
}

func (m *Metric) Init() {
	m.index = make(map[uint64]*Point)
	if m.Capacity == 0 {
		m.Capacity = 86400 //day of seconds
	}
	if m.Indexinterval == 0 {
		m.Indexinterval = 3600 //seconds in hour
	}
}

// Add a pointer to the search index if needed, else do nothing and return.
func (m *Metric) indexPtr(p *Point) {
	bucket := uint64(int(p.Timestamp / m.Indexinterval))
	if _, ok := m.index[bucket]; ok {
		return // bail, point already indexed.
	}
	if p.Timestamp < m.indexoldest { // We need to index a point further in the past than our index covers.
		m.indexoldest = bucket
	}
	if p.Timestamp < m.indexnewest {
		return // bail, don't need to index a point in the past.
	}
	if len(m.index) == 0 { //point is first in index.
		m.indexoldest = bucket
	}
	m.index[bucket] = p // create a new index entry with point
	m.indexnewest = bucket
}

// left trim for the metric stack.
func (m *Metric) Ltrim(t int) {
	for i := 0; i <= t; i++ {
		NextBucket := uint64(int(m.oldest.Next.Timestamp / m.Indexinterval))
		if NextBucket == m.indexoldest { // If the truncate does not push us into the Next bucket.
			m.index[NextBucket] = m.oldest.Next // shift right
		} else { // if it does.
			delete(m.index, m.indexoldest) // drop the empty bucket.
			m.indexoldest = NextBucket     // reset the indexoldest
		}
		m.oldest = m.oldest.Next //reset metric oldest.
		m.length--
	}
}

// Given time(t) and Value(v) create a new point and add it to the metric list.
func (m *Metric) Insert(Timestamp uint64, Value map[string]interface{}, overwriting bool) {
	var p = &Point{Timestamp, Value, nil}
	// check for overflow, deindex, and truncate oldest.
	if m.length >= m.Capacity {
		m.Ltrim(int(m.length - m.Capacity))
	}
	// If our first entry
	if m.newest == nil && m.oldest == nil {
		m.newest = p // Set newest and oldest to point
		m.oldest = p
		m.length++
		m.indexPtr(p)
		return
	}
	if overwriting {
		m.newest.Next = p // set the newest element Next reference to the new point.
		m.newest = p      //set newest node to new point.
		m.length++
		m.indexPtr(p)
		return
	} else {
		// if our new point is newest in series.
		if m.newest.Timestamp < p.Timestamp {
			m.newest.Next = p // set the newest element Next reference to the new point.
			m.newest = p      //set newest node to new point.
			m.length++
			m.indexPtr(p)
			return
		}
		// do an evil past tense insert.
		if m.newest.Timestamp >= p.Timestamp {
			if m.newest.Timestamp == p.Timestamp {
				m.newest.Value = Value // Slight optimization for overwriting the newest point.
				return
			}
			// stub
			m.indexPtr(p)
			return
		}
	}
	// If we get here signal an error.
}

// Given time(t) return a pointer for that point or the Next in list.
func (m *Metric) Search(t uint64) *Point {
	bucket := uint64(int(t / m.Indexinterval))
	startpoint := m.index[bucket]
	// prevent acessing nil point
	if startpoint == nil {
		return &Point{0, nil, nil}
	}
	current := startpoint
	if t == 0 {
		return m.oldest
	}
	for {
		if current.Timestamp >= t {
			return current
		} else {
			if current.Next == nil {
				return &Point{0, nil, nil} //return empty point
			}
			current = current.Next
		}
	}
}

// given unix time(from) and unix time(to) find all points in between(inclusively).
func (m *Metric) GetRange(from uint64, to uint64) []*Point { // Slice of points? Not sure.
	if to == 0 {
		to = 18446744073709551615
	}
	startpoint := m.Search(from)
	result := []*Point{}

	current := startpoint
	for {
		if current.Timestamp > to {
			break
		}
		result = append(result, current)
		if current.Next == nil {
			break
		}
		current = current.Next
	}
	return result
}
