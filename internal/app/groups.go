package app

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

type Groups []int64

func (g Groups) RangeGroup(yield func(group int64) bool) {
	for _, group := range g {
		if !yield(group) {
			return
		}
	}
}

func (g Groups) Rule(r zero.Rule) zero.Rule {
	return func(ctx *zero.Ctx) bool {
		if ctx == nil || ctx.Event == nil {
			return r(ctx)
		}
		target := ctx.Event.TargetID
		if target <= 0 {
			return r(ctx)
		}
		for _, group := range g {
			if group == target {
				return r(ctx)
			}
		}

		return true
	}
}
