package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lck      *sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		lck:      new(sync.Mutex),
	}
}

// добавлениe элемента кеша
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.lck.Lock()
	defer c.lck.Unlock()
	// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*cacheItem).value = value
		return true
	}
	// если элемента нет в словаре, то добавить в словарь и в начало очереди
	// (при этом, если размер очереди больше ёмкости кэша, то необходимо удалить последний элемент из очереди и его значение из словаря);
	if c.queue.Len() == c.capacity {
		c.purge()
	}

	cacheItem := &cacheItem{
		key:   key,
		value: value,
	}

	element := c.queue.PushFront(cacheItem)
	c.items[cacheItem.key] = element
	return false
}

// получение элемента кеша
// - если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
// - если элемента нет в словаре, то вернуть nil и false
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lck.Lock()
	defer c.lck.Unlock()
	element, exists := c.items[key]
	if !exists {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(*cacheItem).value, true
}

// выталкивание последнего элемента кэша
func (c *lruCache) purge() {
	if element := c.queue.Back(); element != nil {
		item := element.Value.(*cacheItem)
		c.queue.Remove(element)
		delete(c.items, item.key)
	}
}

// очистка кэша
func (c *lruCache) Clear() {
	c.lck.Lock()
	defer c.lck.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
