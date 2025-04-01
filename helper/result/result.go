package result

import "errors"

/**
 * result
 * 这么看起来还是 if err == nil 更方便一点 😂
 */

type Result[T any] struct {
	Value T
	Error error
}

func Ok[T any](value T) *Result[T] {
	return &Result[T]{Value: value}
}

func Err[T any](err error) *Result[T] {
	if err == nil {
		return &Result[T]{Error: errors.New("error without details")}
	}
	return &Result[T]{Error: err}
}

func (r *Result[any]) IsOk() bool {
	return r.Error == nil
}

func (r *Result[any]) IsErr() bool {
	return r.Error != nil
}
