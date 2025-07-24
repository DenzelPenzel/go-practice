package main

import (
	"container/list"
	"crypto/rand"
	"time"
)

type Side bool

const (
	BUY  Side = false
	SELL Side = true
)

type Order struct {
	ID        uint64
	Side      Side
	Price     uint64
	Quantity  uint64
	Timestamp int64
}

type OrderQueue struct {
	orders *list.List
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{orders: list.New()}
}
func (o *OrderQueue) Add(order *Order) *list.Element { return o.orders.PushBack(order) }
func (o *OrderQueue) Remove(order *list.Element)     { o.orders.Remove(order) }
func (o *OrderQueue) Front() *list.Element           { return o.orders.Front() }
func (o *OrderQueue) Len() int                       { return o.orders.Len() }

type PriceLevel struct {
	Price      uint64
	OrderQueue *OrderQueue
}

type Trade struct {
	TakerOrderID uint64
	MakerOrderID uint64
	Price        uint64
	Quantity     uint64
	Timestamp    int64
}

const maxLevel = 16

type skipListNode struct {
	price uint64
	value *PriceLevel
	next  []*skipListNode
}

type Comparator func(a, b uint64) bool

type SkipList struct {
	head       *skipListNode
	level      int
	comparator Comparator
}

func NewSkipList(comparator Comparator) *SkipList {
	return &SkipList{
		head: &skipListNode{
			next: make([]*skipListNode, maxLevel),
		},
		level:      0,
		comparator: comparator,
	}
}

func (s *SkipList) Front() *skipListNode {
	return s.head.next[0]
}

// Set inserts or updates a value for a given price. O(log N)
func (s *SkipList) Set(price uint64, value *PriceLevel) {
	update := make([]*skipListNode, maxLevel)
	node := s.head

	// Find the insertion points at each level
	for i := s.level; i >= 0; i-- {
		for node.next[i] != nil && s.comparator(price, node.next[i].price) {
			node = node.next[i]
		}
		update[i] = node
	}

	if node.next[0] != nil && node.next[0].price == price {
		node.next[0].value = value
		return
	}

	// Create a new node with a random level
	newLvl := 0
	for r := rand.Float64(); r < 0.5 && newLvl < maxLevel-1; r = rand.Float64() {
		newLvl++
	}

	if newLvl > s.level {
		for i := s.level + 1; i <= newLvl; i++ {
			update[i] = s.head
		}
		s.level = newLvl
	}

	newNode := &skipListNode{
		price: price,
		value: value,
		next:  make([]*skipListNode, newLvl+1),
	}

	for i := 0; i <= newLvl; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
}

func (s *SkipList) Remove(price uint64) {
	update := make([]*skipListNode, maxLevel)
	// FIX 2: Declare node variable locally
	node := s.head

	for i := s.level; i >= 0; i-- {
		for node.next[i] != nil && s.comparator(price, node.next[i].price) {
			node = node.next[i]
		}
		update[i] = node
	}

	node = node.next[0]

	if node != nil && node.price == price {
		for i := 0; i <= s.level; i++ {
			if update[i].next[i] != node {
				break
			}
			update[i].next[i] = node.next[i]
		}
		for s.level > 0 && s.head.next[s.level] == nil {
			s.level--
		}
	}
}

type OrderBook struct {
	bids   *SkipList
	asks   *SkipList
	orders map[uint64]*list.Element // O(1) order lookup for cancellation
}

func NewOrderBook() *OrderBook {
	// Bids are ordered from highest price to lowest (so we use >)
	bidsComparator := func(a, b uint64) bool { return a > b }

	// Asks are ordered from lowest price to highest (so we use <)
	asksComparator := func(a, b uint64) bool { return a < b }

	return &OrderBook{
		bids:   NewSkipList(bidsComparator),
		asks:   NewSkipList(asksComparator),
		orders: make(map[uint64]*list.Element),
	}
}

func (ob *OrderBook) AddOrder(order *Order) []Trade {
	var trades []Trade

	if order.Side == BUY {
		trades = ob.match(order, ob.asks)
	} else {
		trades = ob.match(order, ob.bids)
	}

	// If the taker order is not fully filled, add it to the book
	if order.Quantity > 0 {
		ob.addOrderToBook(order)
	}

	return trades
}

func (ob *OrderBook) CancelOrder(orderID uint64) {
	el, ok := ob.orders[orderID]
	if !ok {
		return
	}

	order := el.Value.(*Order)
	var book *SkipList
	if order.Side == BUY {
		book = ob.bids
	} else {
		book = ob.asks
	}

	node := findNode(book, order.Price)
	if node == nil {
		return
	}

	priceLevel := node.value

	priceLevel.OrderQueue.Remove(el)
	delete(ob.orders, orderID)

	if priceLevel.OrderQueue.Len() == 0 {
		book.Remove(order.Price)
	}
}

func (ob *OrderBook) match(takerOrder *Order, bookToMatch *SkipList) []Trade {
	trades := make([]Trade, 0, 1)

	for bestPriceNode := bookToMatch.Front(); bestPriceNode != nil && takerOrder.Quantity > 0; {
		bestPriceLevel := bestPriceNode.value

		isMatch := (takerOrder.Side == BUY && takerOrder.Price >= bestPriceLevel.Price) ||
			(takerOrder.Side == SELL && takerOrder.Price <= bestPriceLevel.Price)

		if !isMatch {
			break
		}

		// Time-priority: iterate through orders at this price level's queue
		for bestPriceLevel.OrderQueue.Len() > 0 && takerOrder.Quantity > 0 {
			makerEl := bestPriceLevel.OrderQueue.Front()
			makerOrder := makerEl.Value.(*Order)

			tradeQuantity := min(takerOrder.Quantity, makerOrder.Quantity)

			trades = append(trades, Trade{
				TakerOrderID: takerOrder.ID,
				MakerOrderID: makerOrder.ID,
				Price:        makerOrder.Price,
				Quantity:     tradeQuantity,
				Timestamp:    time.Now().UnixNano(),
			})

			takerOrder.Quantity -= tradeQuantity
			makerOrder.Quantity -= tradeQuantity

			if makerOrder.Quantity == 0 {
				bestPriceLevel.OrderQueue.Remove(makerEl)
				delete(ob.orders, makerOrder.ID)
			}
		}

		// If the price level is now empty, remove it and get the next node.
		if bestPriceLevel.OrderQueue.Len() == 0 {
			// Get the next node *before* removing the current one.
			nextNode := bestPriceNode.next[0]
			bookToMatch.Remove(bestPriceLevel.Price)
			bestPriceNode = nextNode // Advance to the next price level for the next loop iteration.
		} else {
			// The level is not empty, so the taker order must be filled.
			// Break the loop, as no more matching is needed.
			break
		}
	}

	return trades
}

func (ob *OrderBook) addOrderToBook(order *Order) {
	var book *SkipList
	if order.Side == BUY {
		book = ob.bids
	} else {
		book = ob.asks
	}

	// Find the price level in the skip list. O(log N)
	node := findNode(book, order.Price)
	var priceLevel *PriceLevel

	if node != nil {
		priceLevel = node.value
	} else {
		// If it doesn't exist, create it and add it to the skip list. O(log N)
		priceLevel = &PriceLevel{
			Price:      order.Price,
			OrderQueue: NewOrderQueue(),
		}
		book.Set(order.Price, priceLevel)
	}

	// Add order to the queue and the global orders map. O(1)
	el := priceLevel.OrderQueue.Add(order)
	ob.orders[order.ID] = el
}

func findNode(sl *SkipList, price uint64) *skipListNode {
	node := sl.head
	// Traverse from the highest level down
	for i := sl.level; i >= 0; i-- {
		// Move right at the current level
		for node.next[i] != nil && sl.comparator(price, node.next[i].price) {
			node = node.next[i]
		}
		// After finding the correct path, the target node is at the bottom level.
		node = node.next[0]
		if node != nil && node.price == price {
			return node
		}
	}
	return nil
}

func main() {

}
