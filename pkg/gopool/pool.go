package gopool

import "github.com/panjf2000/ants/v2"

var defaultPool, _ = ants.NewPool(-1)

func Go(f func()) {
	_ = defaultPool.Submit(f)
}
