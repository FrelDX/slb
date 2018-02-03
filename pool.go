package slb

import (
	"container/heap"

	"github.com/pkg/errors"
)

type server struct {
	host    string
	index   int
	pending int
}

func newServer(url string, index int) *server {
	return &server{
		host:    url,
		pending: 0,
		index:   index,
	}
}

type pool []*server

func newPool(urls []string) pool {
	pool := make(pool, len(urls))
	for i, url := range urls {
		pool[i] = newServer(url, i)
	}
	heap.Init(&pool)
	return pool
}

func (p pool) Server(host string) (*server, error) {
	for _, server := range p {
		if server.host == host {
			return server, nil
		}
	}
	return nil, errors.Errorf("coudln't find server with host: %s", host)
}

func (p pool) Len() int {
	return len(p)
}

func (p pool) Less(i, j int) bool {
	return p[i].pending > p[j].pending
}

func (p pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *pool) Push(x interface{}) {
	server := x.(*server)
	server.index = p.Len()
	*p = append(*p, server)
}

func (p *pool) Pop() interface{} {
	pool := *p
	last := p.Len() - 1
	elem := pool[last]
	*p = pool[:last]
	return elem
}