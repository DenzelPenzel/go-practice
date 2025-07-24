package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testOrderID uint64

func newTestOrder(side Side, price, quantity uint64) *Order {
	testOrderID++
	return &Order{
		ID:        testOrderID,
		Side:      side,
		Price:     price,
		Quantity:  quantity,
		Timestamp: time.Now().UnixNano(),
	}
}

func setupTest() *OrderBook {
	testOrderID = 0
	return NewOrderBook()
}

func TestAddOrder_NoMatch(t *testing.T) {
	book := setupTest()

	// Add a buy order
	buyOrder := newTestOrder(BUY, 100, 10)
	trades := book.AddOrder(buyOrder)

	assert.Empty(t, trades, "Should be no trades when adding to an empty book")
	assert.Equal(t, 1, book.bids.Len(), "Bid heap should have one price level")
	assert.Equal(t, 0, book.asks.Len(), "Ask heap should be empty")
	assert.NotNil(t, book.orders[buyOrder.ID], "Order should be in the orders map")

	// Add a sell order that doesn't cross the spread
	sellOrder := newTestOrder(SELL, 101, 10)
	trades = book.AddOrder(sellOrder)
	assert.Empty(t, trades, "Should be no trades when spread is not crossed")
	assert.Equal(t, 1, book.asks.Len(), "Ask heap should have one price level")
	assert.Equal(t, uint64(101), book.asks.PriceLevels[0].Price, "Best ask should be 101")
	assert.Equal(t, uint64(100), book.bids.PriceLevels[0].Price, "Best bid should be 100")
}

func TestAddOrder_SimpleFullMatch(t *testing.T) {
	book := setupTest()

	// Add initial sell order
	sellOrder := newTestOrder(SELL, 100, 10)
	book.AddOrder(sellOrder)

	// Add a buy order that fully matches the sell order
	buyOrder := newTestOrder(BUY, 100, 10)
	trades := book.AddOrder(buyOrder)

	assert.Equal(t, 1, len(trades), "Should be exactly one trade")
	trade := trades[0]
	assert.Equal(t, buyOrder.ID, trade.TakerOrderID)
	assert.Equal(t, sellOrder.ID, trade.MakerOrderID)
	assert.Equal(t, uint64(100), trade.Price)
	assert.Equal(t, uint64(10), trade.Quantity)

	// The book should be empty now
	assert.Equal(t, 0, book.bids.Len(), "Bid heap should be empty after full match")
	assert.Equal(t, 0, book.asks.Len(), "Ask heap should be empty after full match")
	assert.Nil(t, book.orders[sellOrder.ID], "Maker order should be removed from the map")
	assert.Nil(t, book.orders[buyOrder.ID], "Taker order should not be in the map")
}

func TestAddOrder_PartialMatchTakerFilled(t *testing.T) {
	book := setupTest()

	// Add a large sell order
	makerOrder := newTestOrder(SELL, 100, 20)
	book.AddOrder(makerOrder)

	// Add a smaller buy order that partially matches
	takerOrder := newTestOrder(BUY, 100, 5)
	trades := book.AddOrder(takerOrder)

	assert.Equal(t, 1, len(trades), "Should be one trade")
	assert.Equal(t, uint64(5), trades[0].Quantity, "Trade quantity should be 5")

	// The book should still have the remaining part of the maker order
	assert.Equal(t, 0, book.bids.Len(), "Bids should be empty")
	assert.Equal(t, 1, book.asks.Len(), "Asks should still have one level")

	// Check that the maker order's quantity was reduced
	makerElement := book.orders[makerOrder.ID]
	assert.NotNil(t, makerElement)
	updatedMakerOrder := makerElement.Value.(*Order)
	assert.Equal(t, uint64(15), updatedMakerOrder.Quantity, "Maker order quantity should be reduced to 15")
}

func TestAddOrder_PartialMatchMakerFilledAndNewOrderAdded(t *testing.T) {
	book := setupTest()

	// Add a small sell order
	makerOrder := newTestOrder(SELL, 100, 5)
	book.AddOrder(makerOrder)

	// Add a larger buy order that fills the maker and becomes a new order
	takerOrder := newTestOrder(BUY, 100, 20)
	trades := book.AddOrder(takerOrder)

	assert.Equal(t, 1, len(trades), "Should be one trade")
	assert.Equal(t, uint64(5), trades[0].Quantity, "Trade quantity should be 5")

	// The original maker order should be gone, and a new bid should be in the book
	assert.Equal(t, 0, book.asks.Len(), "Asks should be empty")
	assert.Equal(t, 1, book.bids.Len(), "Bids should have one new level")

	// Check that the new order in the book is the remainder of the taker order
	newBidElement := book.orders[takerOrder.ID]
	assert.NotNil(t, newBidElement)
	newBidOrder := newBidElement.Value.(*Order)
	assert.Equal(t, uint64(15), newBidOrder.Quantity, "New bid order should have remaining quantity of 15")
	assert.Equal(t, uint64(100), newBidOrder.Price, "New bid order should have the correct price")
}

func TestAddOrder_MultiLevelMatch(t *testing.T) {
	book := setupTest()

	// Add multiple sell orders at different prices
	book.AddOrder(newTestOrder(SELL, 101, 5)) // ID 1
	book.AddOrder(newTestOrder(SELL, 102, 5)) // ID 2
	book.AddOrder(newTestOrder(SELL, 103, 5)) // ID 3

	// Add a large buy order that sweeps through the first two levels
	takerOrder := newTestOrder(BUY, 102, 12)
	trades := book.AddOrder(takerOrder)

	assert.Equal(t, 2, len(trades), "Should be two trades")

	// First trade should be at the best price (101)
	assert.Equal(t, uint64(1), trades[0].MakerOrderID)
	assert.Equal(t, uint64(101), trades[0].Price)
	assert.Equal(t, uint64(5), trades[0].Quantity)

	// Second trade should be at the next best price (102)
	assert.Equal(t, uint64(2), trades[1].MakerOrderID)
	assert.Equal(t, uint64(102), trades[1].Price)
	assert.Equal(t, uint64(5), trades[1].Quantity)

	// The book should have the remainder of the taker order and the last sell order
	assert.Equal(t, 1, book.asks.Len(), "Asks should have one level left")
	assert.Equal(t, uint64(103), book.asks.PriceLevels[0].Price)
	assert.Equal(t, 1, book.bids.Len(), "Bids should have the new remaining taker order")

	newBidOrder := book.orders[takerOrder.ID].Value.(*Order)
	assert.Equal(t, uint64(2), newBidOrder.Quantity, "Remaining taker quantity should be 2")
}

func TestCancelOrder(t *testing.T) {
	book := setupTest()

	// Add an order
	orderToCancel := newTestOrder(BUY, 100, 10)
	book.AddOrder(orderToCancel)
	assert.Equal(t, 1, book.bids.Len(), "Book should have 1 bid before cancellation")

	// Cancel the order
	book.CancelOrder(orderToCancel.ID)
	assert.Equal(t, 0, book.bids.Len(), "Book should be empty after cancellation")
	assert.Nil(t, book.orders[orderToCancel.ID], "Order should be removed from orders map")
	assert.Nil(t, book.prices[orderToCancel.Price], "Price level should be removed from prices map")

}

func TestCancelOrder_LeavesPriceLevelIntact(t *testing.T) {
	book := setupTest()

	// Add two orders at the same price
	order1 := newTestOrder(SELL, 100, 10)
	order2 := newTestOrder(SELL, 100, 10)
	book.AddOrder(order1)
	book.AddOrder(order2)

	assert.Equal(t, 1, book.asks.Len(), "Should be one price level")
	assert.Equal(t, 2, book.prices[uint64(100)].OrderQueue.Len(), "Price level should have two orders")

	// Cancel one of them
	book.CancelOrder(order1.ID)
	assert.Equal(t, 1, book.asks.Len(), "Price level should still exist")
	assert.Equal(t, 1, book.prices[uint64(100)].OrderQueue.Len(), "Price level should now have one order")
	assert.Nil(t, book.orders[order1.ID], "Cancelled order should be gone")
	assert.NotNil(t, book.orders[order2.ID], "Other order should still exist")
}

func TestCancelOrder_NonExistent(t *testing.T) {
	book := setupTest()
	book.AddOrder(newTestOrder(BUY, 100, 10))

	// Attempt to cancel an order ID that doesn't exist
	book.CancelOrder(999)

	// Ensure the book state is unchanged
	assert.Equal(t, 1, book.bids.Len(), "Book state should not change when cancelling a non-existent order")
}

func TestPriceTimePriority(t *testing.T) {
	book := setupTest()

	// Add two orders at the same price, the first one has time priority
	firstOrder := newTestOrder(SELL, 100, 5)
	time.Sleep(1 * time.Millisecond) // Ensure timestamps are different
	secondOrder := newTestOrder(SELL, 100, 5)

	book.AddOrder(firstOrder)
	book.AddOrder(secondOrder)

	// Add a taker order that will match one of them
	takerOrder := newTestOrder(BUY, 100, 5)
	trades := book.AddOrder(takerOrder)

	assert.Equal(t, 1, len(trades), "Should be one trade")

	// The trade should be with the first order placed
	assert.Equal(t, firstOrder.ID, trades[0].MakerOrderID, "Trade should match the first order in time")

	// The second order should still be in the book
	assert.NotNil(t, book.orders[secondOrder.ID], "The second order should remain in the book")
	assert.Nil(t, book.orders[firstOrder.ID], "The first order should have been filled and removed")
}

// =================================================================
//
//	BENCHMARK TESTS
//
// =================================================================

// BenchmarkAddOrder_NoMatch tests the performance of adding orders that do not cross the spread.
// This benchmark focuses on the cost of inserting new price levels and orders into the book.
func BenchmarkAddOrder_NoMatch(b *testing.B) {
	// Explanation:
	// We add a large number of buy and sell orders at prices that won't match each other.
	// Buys are added at decreasing prices (e.g., 99, 98, 97...).
	// Sells are added at increasing prices (e.g., 101, 102, 103...).
	// This isolates the performance of the `addOrderToBook` function, which involves
	// map lookups and heap pushes.

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ob := NewOrderBook()
		testOrderID = 0

		b.StartTimer()

		for j := 0; j < 1000; j++ {
			// Add buy orders below the spread
			ob.AddOrder(newTestOrder(BUY, uint64(100-j), 10))
			// Add sell orders above the spread
			ob.AddOrder(newTestOrder(SELL, uint64(101+j), 10))
		}
	}
}

// BenchmarkAddOrder_FullMatch tests the performance of adding orders that are completely filled.
// This benchmark focuses on the matching engine's performance under high liquidity.
func BenchmarkAddOrder_FullMatch(b *testing.B) {
	// Explanation:
	// The book is pre-filled with a deep stack of sell orders at a single price.
	// Then, we benchmark adding buy orders at a higher price that will consume
	// the existing sell orders completely. This tests the efficiency of the matching loop,
	// order removal, and heap operations (Pop) when a price level is depleted.
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ob := NewOrderBook()
		testOrderID = 0

		for j := 0; j < 1000; j++ {
			ob.AddOrder(newTestOrder(SELL, 100, 10))
		}
		b.StartTimer()

		for j := 0; j < 1000; j++ {
			ob.AddOrder(newTestOrder(BUY, 100, 10))
		}
	}
}

// BenchmarkAddOrder_PartialMatch tests a more realistic scenario of partial fills.
// This benchmark measures the performance of matching, creating trades, and adding
// the remaining order quantity back to the book.
func BenchmarkAddOrder_PartialMatch(b *testing.B) {
	// Explanation:
	// The book is pre-filled with sell orders of a certain size (e.g., quantity 10).
	// We then add larger buy orders (e.g., quantity 15) that will partially fill against
	// the existing orders, create a trade, and then the remainder of the buy order
	// will be added to the book. This tests the full cycle of the `match` function.
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ob := NewOrderBook()
		testOrderID = 0

		for j := 0; j < 1000; j++ {
			ob.AddOrder(newTestOrder(SELL, uint64(100+j), 10))
		}
		b.StartTimer()

		// Add larger buy orders that will partially match and then be added to the book
		for j := 0; j < 1000; j++ {
			ob.AddOrder(newTestOrder(BUY, 101, 15))
		}
	}
}

// BenchmarkCancelOrder measures the performance of canceling existing orders.
// This tests the efficiency of the O(1) order lookup and subsequent removal
// from the order queue and potentially the heap.
func BenchmarkCancelOrder(b *testing.B) {
	// Explanation:
	// We first populate the order book with a large number of orders and store their IDs.
	// Then, we benchmark the time it takes to cancel each of these orders.
	// This measures the speed of map lookups (`ob.orders`) and list/heap removals.
	b.ReportAllocs()
	b.StopTimer()

	ob := NewOrderBook()
	testOrderID = 0
	orderIDs := make([]uint64, 0, 2000)

	for j := 0; j < 1000; j++ {
		buyOrder := newTestOrder(BUY, uint64(100-j), 10)
		sellOrder := newTestOrder(SELL, uint64(101+j), 10)
		ob.AddOrder(buyOrder)
		ob.AddOrder(sellOrder)
		orderIDs = append(orderIDs, buyOrder.ID, sellOrder.ID)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ob.CancelOrder(orderIDs[i%len(orderIDs)])
	}
}

// BenchmarkAddAndCancel_MixedWorkload simulates a more realistic market scenario.
// This tests the data structures' resilience and performance under interleaved
// add and cancel operations.
func BenchmarkAddAndCancel_MixedWorkload(b *testing.B) {
	// Explanation:
	// This benchmark simulates a dynamic market by randomly adding or canceling orders.
	// A list of active order IDs is maintained. On each iteration, we either add a new
	// order or cancel a randomly chosen existing one. This provides a robust test
	// of the order book's overall performance.
	b.ReportAllocs()
	rand.Seed(time.Now().UnixNano())

	b.StopTimer()
	ob := NewOrderBook()
	testOrderID = 0
	activeOrderIDs := make([]uint64, 0, 10000)

	for j := 0; j < 5000; j++ {
		order := newTestOrder(BUY, uint64(90+rand.Intn(10)), 10)
		ob.AddOrder(order)
		activeOrderIDs = append(activeOrderIDs, order.ID)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		// 70% chance to add an order, 30% chance to cancel
		if rand.Float32() < 0.7 || len(activeOrderIDs) == 0 {
			// Add a new order
			price := uint64(90 + rand.Intn(20)) // Prices around the spread
			var order *Order
			if rand.Float32() < 0.5 {
				order = newTestOrder(BUY, price, 10)
			} else {
				order = newTestOrder(SELL, price+10, 10) // Sell orders at a higher price range
			}
			ob.AddOrder(order)
			activeOrderIDs = append(activeOrderIDs, order.ID)
		} else {
			// Cancel an existing order
			cancelIdx := rand.Intn(len(activeOrderIDs))
			orderToCancelID := activeOrderIDs[cancelIdx]

			ob.CancelOrder(orderToCancelID)

			// Remove the canceled ID from our list
			activeOrderIDs[cancelIdx] = activeOrderIDs[len(activeOrderIDs)-1]
			activeOrderIDs = activeOrderIDs[:len(activeOrderIDs)-1]
		}
	}
}
