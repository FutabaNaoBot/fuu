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

func (g Groups) IsContains(groupId int64) bool {
	for _, group := range g {
		if group == groupId {
			return true
		}
	}
	return false
}

func (g Groups) Rule() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		groupId := ctx.Event.GroupID
		if groupId <= 0 {
			return false
		}

		return g.IsContains(groupId)
	}
}
