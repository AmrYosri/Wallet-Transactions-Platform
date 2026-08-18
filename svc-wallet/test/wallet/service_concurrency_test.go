package wallet_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"svc-wallet/external/mongodb"
	"svc-wallet/internal/wallet"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func TestApplyBalanceChange_ConcurrentWithdrawals(t *testing.T) {
	// .env lives at svc-wallet/, this test file lives at svc-wallet/test/wallet/,
	// so it's two levels up. Adjust if your actual layout differs.
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}

	// isolate this test's writes from your real dev database
	os.Setenv("MONGO_NAME", "svc-wallet-test")

	db, err := mongodb.Connect()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	repo := wallet.NewRepository(db)
	// ApplyBalanceChange and GetWallet never touch userClient (only CreateWallet does),
	// so nil is safe here — passing a real *user.Client would need a live svc-user gRPC
	// connection this test has no reason to depend on.
	service := wallet.NewService(repo, nil)

	ctx := context.Background()

	seedWallet := &wallet.Wallet{
		OwnerPhone: "01000000000",
		Currency:   "EGP",
		Balance:    100,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := repo.Create(ctx, seedWallet); err != nil {
		t.Fatalf("failed to seed wallet: %v", err)
	}
	walletID := seedWallet.ID.Hex()

	t.Cleanup(func() {
		_, _ = db.Collection("wallets").DeleteOne(ctx, bson.M{"_id": seedWallet.ID})
	})

	const numGoroutines = 10
	withdrawAmount := int64(100) // only one of these can succeed

	results := make(chan error, numGoroutines)
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := service.ApplyBalanceChange(ctx, walletID, "withdraw", withdrawAmount)
			results <- err
		}()
	}

	wg.Wait()
	close(results)

	var successCount, insufficientCount, otherErrCount int
	for err := range results {
		switch {
		case err == nil:
			successCount++
		case err.Error() == "insufficient funds":
			insufficientCount++
		default:
			otherErrCount++
			t.Logf("unexpected error: %v", err)
		}
	}

	if successCount != 1 {
		t.Errorf("expected exactly 1 successful withdrawal, got %d", successCount)
	}
	if insufficientCount != numGoroutines-1 {
		t.Errorf("expected %d insufficient-funds rejections, got %d", numGoroutines-1, insufficientCount)
	}
	if otherErrCount != 0 {
		t.Errorf("got %d unexpected errors", otherErrCount)
	}

	finalWallet, err := service.GetWallet(ctx, walletID)
	if err != nil {
		t.Fatalf("failed to fetch final wallet state: %v", err)
	}
	if finalWallet.Balance != 0 {
		t.Errorf("expected final balance 0, got %d", finalWallet.Balance)
	}
}