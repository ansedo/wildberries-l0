package models

import (
	"context"
	"sync"
)

const lastOrdersCount = 6

type Cache struct {
	sync.RWMutex
	orders       map[string]Order
	itemsCount   int
	newOrderUIDs []string
	lastOrders   []string
}

func NewCache(_ context.Context) *Cache {
	return &Cache{
		orders: make(map[string]Order),
	}
}

func (c *Cache) GetOrderCount() int {
	return len(c.orders)
}
func (c *Cache) GetItemCount() int {
	return c.itemsCount
}

func (c *Cache) NewOrderCount() int {
	return len(c.newOrderUIDs)
}

func (c *Cache) AddOrder(order Order) {
	c.Lock()
	defer c.Unlock()
	if !order.InDB {
		c.newOrderUIDs = append(c.newOrderUIDs, order.OrderUID)
	}
	c.orders[order.OrderUID] = order
	if len(c.lastOrders) >= lastOrdersCount {
		c.lastOrders = c.lastOrders[1:]
	}
	c.lastOrders = append(c.lastOrders, order.OrderUID)
	c.itemsCount += len(order.Items)
}

func (c *Cache) AddItem(orderUID string, item Item) {
	c.Lock()
	defer c.Unlock()
	if order, ok := c.orders[orderUID]; ok {
		order.Items = append(order.Items, item)
	}
	c.itemsCount++
}

func (c *Cache) GetByOrderUID(orderUID string) Order {
	c.RLock()
	defer c.RUnlock()
	return c.orders[orderUID]
}

func (c *Cache) GetLastOrderUIDs() []string {
	return c.lastOrders
}

func (c *Cache) GetNewOrders() []Order {
	newOrders := make([]Order, len(c.newOrderUIDs))
	for i, orderUID := range c.newOrderUIDs {
		newOrders[i] = c.GetByOrderUID(orderUID)
	}
	return newOrders
}

func (c *Cache) DeleteNewOrderUIDs() {
	c.newOrderUIDs = []string{}
}
