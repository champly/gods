package dag

import (
	"context"
	"sync"
	"sync/atomic"
)

type Edge struct {
	From *Vertex
	To   *Vertex
}

type Vertex struct {
	Dependency   []*Edge
	DepCompleted int32
	Task         Runnable
	Children     []*Edge
}

type Runnable interface {
	Run(i interface{})
}

type WorkFlow struct {
	done        chan struct{}
	doneOnce    *sync.Once
	alreadyDone bool
	root        *Vertex
	end         *Vertex
	edges       []*Edge
}

type EndWorkFlowAction struct {
	done chan struct{}
	s    *sync.Once
}

// Run implements Runnable.
func (e *EndWorkFlowAction) Run(i interface{}) {
	e.s.Do(func() {
		e.done <- struct{}{}
	})
}

func NewVertex(task Runnable) *Vertex {
	return &Vertex{
		Task: task,
	}
}

func AddEdge(from *Vertex, to *Vertex) *Edge {
	edge := &Edge{
		From: from,
		To:   to,
	}

	from.Children = append(from.Children, edge)
	to.Dependency = append(to.Dependency, edge)

	return edge
}

func NewWorkFlow() *WorkFlow {
	wf := &WorkFlow{
		root:     &Vertex{Task: nil},
		done:     make(chan struct{}),
		doneOnce: &sync.Once{},
	}

	end := &EndWorkFlowAction{
		done: wf.done,
		s:    wf.doneOnce,
	}
	wf.end = NewVertex(end)

	return wf
}

func (w *WorkFlow) AddStartVertex(vertex *Vertex) {
	w.edges = append(w.edges, AddEdge(w.root, vertex))
}

func (w *WorkFlow) AddEdge(from, to *Vertex) {
	w.edges = append(w.edges, AddEdge(from, to))
}

func (w *WorkFlow) ConnectToEnd(vertex *Vertex) {
	w.edges = append(w.edges, AddEdge(vertex, w.end))
}

func (w *WorkFlow) Start(ctx context.Context, i interface{}) {
}

func (w *WorkFlow) WaitDone() {
	<-w.done
	close(w.done)
}

func (w *WorkFlow) interrupDone() {
	w.alreadyDone = true
	w.doneOnce.Do(func() { w.done <- struct{}{} })
}

func (v *Vertex) Execute(ctx context.Context, w *WorkFlow, i interface{}) {
	if !v.dependencyHasDone() {
		return
	}

	if ctx.Err() != nil {
		w.interrupDone()
		return
	}

	if v.Task != nil {
		v.Task.Run(i)
	}

	if len(v.Children) > 0 {
		for idx := 1; idx < len(v.Children); idx++ {
			go func(child *Edge) {
				child.To.Execute(ctx, w, i)
			}(v.Children[idx])
		}
	}

	v.Children[0].To.Execute(ctx, w, i)
}

func (v *Vertex) dependencyHasDone() bool {
	if len(v.Dependency) <= 1 {
		return true
	}

	atomic.AddInt32(&v.DepCompleted, 1)

	return v.DepCompleted == int32(len(v.Dependency))
}
