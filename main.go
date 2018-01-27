package main

import (
	"fmt"
)

// Single point
type Point struct {
	timestamp uint64
	value     map[string]interface{}
	next      *Point
}

// Core linked list, main storage array.
type Metric struct {
	newest, oldest           *Point
	length                   uint64
	capacity                 uint64
	index                    map[uint64]*Point // coarseSearch index
	indexnewest, indexoldest uint64
	indexinterval            uint64 // in seconds > 0
}

func (m *Metric) Init() {
	m.index = make(map[uint64]*Point)
	if m.capacity == 0 {
		m.capacity = 86400 //day of seconds
	}
	if m.indexinterval == 0 {
		m.indexinterval = 3600 //seconds in hour
	}
}

// Add a pointer to the search index if needed, else do nothing and return.
func (m *Metric) indexPtr(p *Point) {
	bucket := uint64(int(p.timestamp / m.indexinterval))
	if _, ok := m.index[bucket]; ok {
		return // bail, point already indexed.
	}
	if p.timestamp < m.indexoldest { // We need to index a point further in the past than our index covers.
		m.indexoldest = bucket
	}
	if p.timestamp < m.indexnewest {
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
		nextBucket := uint64(int(m.oldest.next.timestamp / m.indexinterval))
		if nextBucket == m.indexoldest { // If the truncate does not push us into the next bucket.
			m.index[nextBucket] = m.oldest.next // shift right
		} else { // if it does.
			delete(m.index, m.indexoldest) // drop the empty bucket.
			m.indexoldest = nextBucket     // reset the indexoldest
		}
		m.oldest = m.oldest.next //reset metric oldest.
		m.length--
	}
}

// Given time(t) and value(v) create a new point and add it to the metric list.
func (m *Metric) Insert(timestamp uint64, value map[string]interface{}, overwriting bool) {
	var p = &Point{timestamp, value, nil}
	// check for overflow, deindex, and truncate oldest.
	if m.length >= m.capacity {
		m.Ltrim(int(m.length - m.capacity))
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
		m.newest.next = p // set the newest element next reference to the new point.
		m.newest = p      //set newest node to new point.
		m.length++
		m.indexPtr(p)
		return
	} else {
		// if our new point is newest in series.
		if m.newest.timestamp < p.timestamp {
			m.newest.next = p // set the newest element next reference to the new point.
			m.newest = p      //set newest node to new point.
			m.length++
			m.indexPtr(p)
			return
		}
		// do an evil past tense insert.
		if m.newest.timestamp >= p.timestamp {
			if m.newest.timestamp == p.timestamp {
				m.newest.value = value // Slight optimization for overwriting the newest point.
				return
			}
			// stub
			m.indexPtr(p)
			return
		}
	}
	// If we get here signal an error.
}

// Given time(t) return a pointer for that point or the next in list.
func (m *Metric) Search(t uint64) *Point {
	bucket := uint64(int(t / m.indexinterval))
	startpoint := m.index[bucket]
	current := startpoint
	if t == 0 {
		return m.oldest
	}
	for {
		if current.timestamp >= t {
			return current
		} else {
			if current.next == nil {
				return &Point{0, nil, nil} //return empty point
			}
			current = current.next
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
		if current.timestamp > to {
			break
		}
		result = append(result, current)
		if current.next == nil {
			break
		}
		current = current.next
	}
	return result
}

func main() {
	testmetric := Metric{capacity: 131487, indexinterval: 1440} // 3 months of per minute, indexed on days
	testmetric.Init()

	test := map[string]interface{}{
		"test": "mauro",
	}

	testmetric.Insert(1416585010, test, true)
	testmetric.Insert(1416585011, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585012, test, true)
	testmetric.Insert(1416585013, test, true)
	testmetric.Insert(1416585014, test, true)
	testmetric.Insert(1416585015, test, true)
	testmetric.Insert(1416585015, test, true)

	for _, pnt := range testmetric.GetRange(1416585013, 1416585015) {
		fmt.Printf("Found: %+v\n", pnt)
	}
}
