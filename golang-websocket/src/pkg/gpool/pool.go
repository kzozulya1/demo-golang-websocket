package gpool

import "fmt"

//GPool maintain max N items
type GPool struct {
	work chan func()
	sem  chan bool
}

//New creates new pool
func New(size int) *GPool {
	return &GPool{
		work: make(chan func()),
		sem:  make(chan bool, size),
	}
}

//Schedule adds task to list
func (p *GPool) Schedule(task func()) {
	//fmt.Println("GPool: Schedule called")
	select {
	case p.work <- task:
	case p.sem <- true:
		fmt.Println("GPool: allocate new worker: p.worker(task)")
		go p.worker(task)
	}
}

func (p *GPool) worker(task func()) {
	for {
		//fmt.Println("GPool: run task")
		task()
		//fmt.Println("GPool: run task end, wait for data in p.work")
		task = <-p.work
		//fmt.Println("GPool: got data! done and quit")
	}
}
