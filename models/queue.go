package models

type Queue struct {
	adresses []string
	urls     []string
	index    int
	length   int
}

func (q *Queue) Pop() (string, string) {
	if q.index >= 10000 {
		q.Unload()
	}
	address, url := q.adresses[q.index], q.urls[q.index]
	q.index++
	q.length--
	return address, url
}
func (q *Queue) Push(adress string, url string) {
	q.length++
	q.adresses = append(q.adresses, adress)
	q.urls = append(q.urls, url)
}

func NewQueue() Queue {
	que := Queue{}
	que.adresses = make([]string, 0)
	que.urls = make([]string, 0)
	que.index = 0
	que.length = 0
	return que
}

func (q *Queue) Unload() {
	if q.length == 0 {
		q.adresses = make([]string, 0)
		q.urls = make([]string, 0)
	} else {
		q.adresses = q.adresses[q.index:]
		q.urls = q.urls[q.index:]
	}
	q.index = 0
}

func (q *Queue) IsEmpty() bool {
	return q.length == 0
}
