package validator

import (
	"sync"

	"github.com/Rhymond/go-money"
	"github.com/gookit/validate"
)

func ConfigureDefaultValidator() {
	sync.OnceFunc(func() {
		validate.Config(func(opt *validate.GlobalOption) {
			opt.StopOnError = false
			opt.SkipOnEmpty = false
		})

		validate.AddValidator("money_amount", func(val any) bool {
			v, ok := val.(money.Amount)

			if !ok || v <= 0 {
				return false
			}

			return true
		})
	})()
}
