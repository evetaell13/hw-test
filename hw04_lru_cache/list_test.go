package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)

		first := l.Front() // удалим первый элемент // [80, 60, 40, 10, 30, 50]
		l.Remove(first)
		require.Equal(t, 6, l.Len())
		require.Equal(t, 80, l.Front().Value)

		last := l.Back() // удалим последний элемент  // [80, 60, 40, 10, 30]
		l.Remove(last)
		require.Equal(t, 5, l.Len())
		require.Equal(t, 30, l.Back().Value)

		middle2 := l.Back().Prev // удалим элемент ближе к концу // [80, 60, 40, 30]
		l.Remove(middle2)
		require.Equal(t, 4, l.Len())

		l.PushFront(90) // [90, 80, 60, 40, 30]
		l.PushBack(100) // [90, 80, 60, 40, 30, 100]
		require.Equal(t, 90, l.Front().Value)
		require.Equal(t, 100, l.Back().Value)
		require.Equal(t, 6, l.Len())

		l.MoveToFront(l.Back().Prev)  // [30, 90, 80, 60, 40, 100]
		l.MoveToFront(l.Front().Next) // [90, 30, 80, 60, 40, 100]
		require.Equal(t, 90, l.Front().Value)

		require.Nil(t, l.Front().Prev)
		require.Nil(t, l.Back().Next)
	})
}
