package app

import zero "github.com/wdvxdr1123/ZeroBot"

type Users []int64

func (u Users) IsContains(userId int64) bool {
	for _, user := range u {
		if user == userId {
			return true
		}
	}
	return false
}

func (u Users) Rule() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		userId := ctx.Event.Sender.ID
		if userId <= 0 {
			return false
		}

		return u.IsContains(userId)
	}
}

func (u Users) RangeUser(yield func(user int64) bool) {
	for _, user := range u {
		if !yield(user) {
			return
		}
	}
}
