/*

Skip List vs. Heap

Feature	                      Skip List  	Heap 		Winner
Find Best Price				  O(1)			O(1)		Tie
Add/Remove Price Level		  O(log N)		O(log N)	Tie
Cancel Order				  O(log N)		O(1)		Skip List
Match Across Levels		  	  O(1)			O(log N)	Skip List
*/

package main

import (
	"container/heap"
	"container/list"
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

// OrderQueue is a queue of orders at a specific price level.
// Implemented as a doubly linked list for O(1) add/remove operations.
type OrderQueue struct {
	orders *list.List
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		orders: list.New(),
	}
}

func (o *OrderQueue) Add(order *Order) *list.Element {
	return o.orders.PushBack(order)
}

func (o *OrderQueue) Remove(order *list.Element) {
	o.orders.Remove(order)
}

func (o *OrderQueue) Front() *list.Element {
	return o.orders.Front()
}

func (o *OrderQueue) Len() int {
	return o.orders.Len()
}

type PriceLevel struct {
	Price      uint64
	OrderQueue *OrderQueue
	heapIndex  int // index of the item in the heap
}

type Trade struct {
	TakerOrderID uint64
	MakerOrderID uint64
	Price        uint64
	Quantity     uint64
	Timestamp    int64
}

// Heap implementation
type PriceLevels []*PriceLevel

func (pl PriceLevels) Len() int {
	return len(pl)
}

func (pl PriceLevels) Less(i, j int) bool {
	return pl[i].Price < pl[j].Price
}

func (pl PriceLevels) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
	pl[i].heapIndex = i
	pl[j].heapIndex = j
}

func (pl *PriceLevels) Push(x interface{}) {
	n := len(*pl)
	item := x.(*PriceLevel)
	item.heapIndex = n // Set the index for the new item
	*pl = append(*pl, item)
}

func (pl *PriceLevels) Pop() interface{} {
	old := *pl
	n := len(old)
	item := old[n-1]
	old[n-1] = nil      // avoid memory leak
	item.heapIndex = -1 // for safety
	*pl = old[0 : n-1]
	return item
}

// MaxPriceLevels implements heap.Interface for a max-heap (for bids).
type MaxPriceLevels struct {
	PriceLevels
}

func (mpl MaxPriceLevels) Less(i, j int) bool {
	return mpl.PriceLevels[i].Price > mpl.PriceLevels[j].Price
}

// MinPriceLevels implements heap.Interface for a min-heap (for asks).
type MinPriceLevels struct {
	PriceLevels
}

func (mpl MinPriceLevels) Less(i, j int) bool {
	return mpl.PriceLevels[i].Price < mpl.PriceLevels[j].Price
}

type OrderBook struct {
	bids *MaxPriceLevels
	asks *MinPriceLevels

	prices map[uint64]*PriceLevel   // Map price to PriceLevel
	orders map[uint64]*list.Element // For O(1) order lookup/cancellation
}

func NewOrderBook() *OrderBook {
	bids := &MaxPriceLevels{
		PriceLevels: make([]*PriceLevel, 0),
	}
	asks := &MinPriceLevels{
		PriceLevels: make([]*PriceLevel, 0),
	}
	heap.Init(bids)
	heap.Init(asks)

	return &OrderBook{
		bids:   bids,
		asks:   asks,
		orders: make(map[uint64]*list.Element),
		prices: make(map[uint64]*PriceLevel),
	}
}

func (ob *OrderBook) GetBestAsk() (uint64, bool) {
	if ob.asks.Len() == 0 {
		return 0, false
	}
	return ob.asks.PriceLevels[0].Price, true
}

func (ob *OrderBook) AddOrder(order *Order) []Trade {
	if order.Side == BUY {
		return ob.match(order, &ob.asks.PriceLevels, &ob.bids.PriceLevels)
	} else {
		return ob.match(order, &ob.bids.PriceLevels, &ob.asks.PriceLevels)
	}
}

func (ob *OrderBook) CancelOrder(orderId uint64) {
	el, ok := ob.orders[orderId]
	if !ok {
		return
	}

	order := el.Value.(*Order)
	priceLevel, ok := ob.prices[order.Price]
	if !ok {
		return
	}

	priceLevel.OrderQueue.Remove(el)
	delete(ob.orders, orderId)

	// If price level is now empty, remove it from the heap and the price map
	if priceLevel.OrderQueue.Len() == 0 {
		delete(ob.prices, order.Price)

		if order.Side == BUY {
			heap.Remove(&ob.bids.PriceLevels, priceLevel.heapIndex)
		} else {
			heap.Remove(&ob.asks.PriceLevels, priceLevel.heapIndex)
		}
	}
}

func (ob *OrderBook) match(takerOrder *Order, bookToMatch *PriceLevels, bookToAdd *PriceLevels) []Trade {
	trades := make([]Trade, 0, 1) // Pre-allocate slice capacity

	for bookToMatch.Len() > 0 && takerOrder.Quantity > 0 {
		bestPriceLevel := (*bookToMatch)[0]

		// Price-priority check: Can the taker order be matched at the best price?
		isMatch := (takerOrder.Side == BUY && takerOrder.Price >= bestPriceLevel.Price) || (takerOrder.Side == SELL && takerOrder.Price <= bestPriceLevel.Price)

		if !isMatch {
			break
		}

		// Time-priority: iterate through orders at the best price level
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

		// If the price level is empty, remove it from the heap
		if bestPriceLevel.OrderQueue.Len() == 0 {
			heap.Pop(bookToMatch)
			delete(ob.prices, bestPriceLevel.Price)
		}
	}

	// If taker order is not fully filled, add the remainder to the book.
	if takerOrder.Quantity > 0 {
		ob.addOrderToBook(takerOrder, bookToAdd)
	}

	return trades
}

func (ob *OrderBook) addOrderToBook(order *Order, book *PriceLevels) {
	priceLevel, ok := ob.prices[order.Price]
	if !ok {
		priceLevel = &PriceLevel{Price: order.Price, OrderQueue: NewOrderQueue()}
		ob.prices[order.Price] = priceLevel
		heap.Push(book, priceLevel)
	}
	el := priceLevel.OrderQueue.Add(order)
	ob.orders[order.ID] = el
}
