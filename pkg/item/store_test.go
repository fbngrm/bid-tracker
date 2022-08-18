package item

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// A note on compiler optimisations
// To be completely accurate, any benchmark should be careful to avoid compiler optimisations eliminating
// the function under test and artificially lowering the run time of the benchmark.

func TestRegister(t *testing.T) {
	type testcase struct {
		name      string
		testItems []*Item
		err       error
	}

	t.Run("expect successful concurrent registration", func(t *testing.T) {
		t.Parallel()

		tc := testcase{
			name: "expect no error",
			testItems: []*Item{
				{ID: uuid.New()}, // note, might panic,
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
			err: nil,
		}

		store := newStore()
		ctx := context.Background()

		wg := sync.WaitGroup{}
		for _, it := range tc.testItems {
			wg.Add(1)
			go func(t *testing.T, it *Item) {
				defer wg.Done()
				err := store.register(ctx, it)
				if err != nil {
					t.Logf("unexpected error: %v", err)
				}
			}(t, it)
		}
		wg.Wait()
		assert.Equal(t, len(tc.testItems), len(store.items))
	})

	t.Run("expect two misses", func(t *testing.T) {
		t.Parallel()

		tc := testcase{
			name: "expect no error",
			testItems: []*Item{
				{ID: uuid.New()}, // note, might panic,
				{ID: uuid.New()},
				nil,
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
				nil,
				{ID: uuid.New()},
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
			err: errors.New("could not register, item is nil"),
		}

		store := newStore()
		ctx := context.Background()

		wg := sync.WaitGroup{}
		for _, it := range tc.testItems {
			wg.Add(1)
			go func(t *testing.T, it *Item, terr error) {
				defer wg.Done()
				err := store.register(ctx, it)
				if err != nil && err.Error() != terr.Error() {
					t.Logf("unexpected error: %v", err)
				}
			}(t, it, tc.err)
		}
		wg.Wait()
		assert.Equal(t, len(tc.testItems)-2, len(store.items))
	})

}

func TestWrite(t *testing.T) {
	type testcase struct {
		name      string
		testItems []*Item
		testBids  []*bid.Bid
		err       error
	}

	itemID1 := uuid.New() // note, might panic,
	itemID2 := uuid.New()
	amount := 1.

	tc := testcase{
		name: "expect no error",
		testItems: []*Item{
			{ID: itemID1},
			{ID: itemID2},
		},
		testBids: []*bid.Bid{
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.5)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.5)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.6)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.6)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.7)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.7)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.8)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.8)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.9)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.9)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.1)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.1)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.2)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.2)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.3)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.3)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.4)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.4)},
			{ID: uuid.New(), ItemID: itemID1, Amount: float32(amount + 0.5)},
			{ID: uuid.New(), ItemID: itemID2, Amount: float32(amount + 0.5)},
		},
		err: nil,
	}

	store := newStore()
	ctx := context.Background()

	// prepare by registering items
	for _, it := range tc.testItems {
		err := store.register(ctx, it)
		if err != nil {
			t.Logf("unexpected error: %v", err)
			t.FailNow()
		}
	}

	wg := sync.WaitGroup{}
	for _, it := range tc.testBids {
		wg.Add(1)
		go func(t *testing.T, b *bid.Bid) {
			defer wg.Done()
			err := store.write(ctx, b)
			if err != nil {
				t.Logf("unexpected error: %v", err)
			}
		}(t, it)
	}
	wg.Wait()

	assert.Equal(t, 10, len(store.items[itemID1].bids))
	assert.Equal(t, 10, len(store.items[itemID2].bids))
}

func BenchmarkWrite(b *testing.B) {
	itemID := uuid.New()
	i := &Item{
		ID: itemID,
	}
	bi := &bid.Bid{ID: uuid.New(), ItemID: itemID, Amount: float32(0.5)}

	store := newStore()
	ctx := context.Background()
	// prepare by registering items
	err := store.register(ctx, i)
	if err != nil {
		b.Logf("unexpected error: %v", err)
		b.FailNow()
	}

	// run the write function b.N times
	for n := 0; n < b.N; n++ {
		_ = store.write(ctx, bi)
	}
}
