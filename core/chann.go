package core

import (
	"fmt"
)

type ChanN []Chan1

func MakeChanN(nComp int, tag, unit string, m *Mesh, blocks ...int) ChanN {
	c := make(ChanN, nComp)
	for i := range c {
		c[i] = MakeChan(fmt.Sprint(tag, i), unit, m, blocks...)
	}
	return c
}

// WriteNext locks and returns a slice of length n for 
// writing the next n elements to the Chan3.
// When done, WriteDone() should be called to "send" the
// slice down the Chan3. After that, the slice is not valid any more.
func (c ChanN) WriteNext(n int) [3][]float32 {
	var next [3][]float32
	for i := range c {
		c[i].WriteNext(n)
		a, b := c[i].mutex.WRange()
		next[i] = c[i].list[a:b]
	}
	return next
}

// WriteDone() signals a slice obtained by WriteNext() is fully
// written and can be sent down the Chan3.
func (c ChanN) WriteDone() {
	for i := range c {
		c[i].WriteDone()
	}
}

func (c ChanN) WriteDelta(Δstart, Δstop int) [3][]float32 {
	var next [3][]float32
	for i := range c {
		c[i].WriteDelta(Δstart, Δstop)
		a, b := c[i].mutex.WRange()
		next[i] = c[i].list[a:b]
	}
	return next
}

func (c ChanN) Mesh() *Mesh  { return c[0].Mesh }
func (c ChanN) Size() [3]int { return c[0].Size() }
func (c ChanN) Unit() string { return c[0].Unit() }
func (c ChanN) Tag() string  { return c[0].Tag() }
func (c ChanN) NComp() int  { return len(c) }

// UnsafeData returns the underlying storage without locking.
// Intended only for page-locking, not for reading or writing.
func (c ChanN) UnsafeData() [3][]float32 { return [3][]float32{c[0].list, c[1].list, c[2].list} }
func (c ChanN) UnsafeArray() [3][][][]float32 {
	return [3][][][]float32{c[0].array, c[1].array, c[2].array}
}
func (c ChanN) Comp(idx int) Chan1 { return c[idx] }
