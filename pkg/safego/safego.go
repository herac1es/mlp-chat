package safego

import (
	"context"
	"errors"
	"fmt"
	"runtime"
)

func Go(ctx context.Context, f func(ctx context.Context)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				for {
					n := runtime.Stack(buf, false)
					if n < len(buf) {
						break
					}
					buf = make([]byte, len(buf)*2)
				}
				err = errors.New(fmt.Sprintf("panic %s\n%s", err, buf))
				fmt.Printf("%v", err)
			}
		}()
		f(ctx)
	}()
}
