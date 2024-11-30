// Package io implements the IO type.
package io

import "github.com/tetsuo/warp"

func Map[A, B any](fa warp.IO[A], f func(A) B) warp.IO[B] {
	return func() B {
		return f(fa())
	}
}

func Ap[A, B any](fab warp.IO[func(A) B], fa warp.IO[A]) warp.IO[B] {
	return func() B {
		return fab()(fa())
	}
}

func Chain[A, B any](ma warp.IO[A], f func(A) warp.IO[B]) warp.IO[B] {
	return func() B {
		return f(ma())()
	}
}

func ApFirst[A, B any](fa warp.IO[A], fb warp.IO[B]) warp.IO[A] {
	return Ap(Map(fa, func(a A) func(B) A {
		return func(_ B) A {
			return a
		}
	}), fb)
}

func ApSecond[A, B any](fa warp.IO[A], fb warp.IO[B]) warp.IO[B] {
	return Ap(Map(fa, func(_ A) func(B) B {
		return func(b B) B {
			return b
		}
	}), fb)
}

func ChainFirst[A, B any](ma warp.IO[A], f func(A) warp.IO[B]) warp.IO[A] {
	return Chain(ma, func(a A) warp.IO[A] {
		return Map(f(a), func(_ B) A {
			return a
		})
	})
}

func ChainRec[A, B any](init A, f func(A) warp.IO[func() (A, B, bool)]) warp.IO[B] {
	return func() B {
		var (
			a  A
			b  B
			ok bool
		)

		a, b, ok = f(init)()()

		for {
			if ok {
				break
			}
			a, b, ok = f(a)()()
		}

		return b
	}
}

func Of[A any](a A) warp.IO[A] {
	return func() A {
		return a
	}
}
