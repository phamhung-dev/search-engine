package utils

type Stack[T any] struct {
	vals []T
}

func (stack *Stack[T]) Top() T {
	if stack.IsEmpty() {
		var zero T

		return zero
	}

	lenStack := len(stack.vals)

	top := stack.vals[lenStack-1]

	return top
}

func (stack *Stack[T]) Push(val T) {
	stack.vals = append(stack.vals, val)
}

func (stack *Stack[T]) BatchPush(vals ...T) {
	stack.vals = append(stack.vals, vals...)
}

func (stack *Stack[T]) Pop() (T, bool) {
	if stack.IsEmpty() {
		var zero T

		return zero, false
	}

	lenStack := len(stack.vals)

	top := stack.vals[lenStack-1]

	stack.vals = stack.vals[:lenStack-1]

	return top, true
}

func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.vals) == 0
}
