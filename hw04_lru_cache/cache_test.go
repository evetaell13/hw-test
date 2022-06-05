package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// логика выталкивания: из-за переполнения
		c := NewCache(3)
		c.Set("aaa", 100)
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
		// логика выталкивания: давно используемые элементы
		c.Set("bbb", 50)       // bbb вперед
		c.Get("ccc")           // ссс вперед
		c.Get("bbb")           // bbb вперед
		c.Set("ccc", 30)       // ccc вперед
		c.Set("mmm", 130)      // mmm новое
		val, ok = c.Get("ddd") //ddd вытолкнуло за давностью
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear logic", func(t *testing.T) {
		c := NewCache(1)
		c.Set("vvv", 200)
		val, ok := c.Get("vvv")
		require.True(t, ok)
		require.Equal(t, 200, val)
		c.Clear()
		_, ok = c.Get("vvv")
		require.False(t, ok)
		cap := c.(*lruCache).capacity
		len := c.(*lruCache).queue.Len()
		require.Equal(t, 0, len)
		require.Equal(t, 1, cap)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
