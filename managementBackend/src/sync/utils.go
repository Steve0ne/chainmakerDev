/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package sync

type Task struct {
	f func() error
}

func NewTask(f func() error) *Task {
	return &Task{
		f: f,
	}
}
func (t *Task) execute() {
	err := t.f()
	if err != nil {
		log.Error("Task execute err :", err.Error())
	}
}

type Pool struct {
	workerNum  int
	EntryChan  chan *Task
	workerChan chan *Task
}

func NewPool(num int) *Pool {
	return &Pool{
		workerNum:  num,
		EntryChan:  make(chan *Task),
		workerChan: make(chan *Task),
	}
}

func (p *Pool) worker() {
	for task := range p.workerChan {
		task.execute()
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.workerNum; i++ {
		go p.worker()
	}
	for task := range p.EntryChan {
		p.workerChan <- task
	}
}

type SingletonSync struct {
	SyncStart bool
}
