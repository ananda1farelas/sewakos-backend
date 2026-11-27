package repository

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager adalah satu-satunya tempat
// transaksi database dibuka dan ditutup.
type TxManager struct {
	pool        *pgxpool.Pool
	maxAttempts int
}

// NewTxManager membuat TxManager baru.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool:        pool,
		maxAttempts: 3,
	}
}

// WithTx menjalankan fungsi di dalam satu transaksi.
// Jika fungsi menghasilkan error, seluruh perubahan dibatalkan.
// Jika berhasil, seluruh perubahan disimpan.
func (m *TxManager) WithTx(
	ctx context.Context,
	fn func(tx DBTX) error,
) error {
	var err error

	for attempt := 1; attempt <= m.maxAttempts; attempt++ {
		err = m.runOnce(ctx, fn)

		if !isRetryable(err) {
			return err
		}

		// Memberikan jeda acak sebelum mencoba kembali.
		wait := time.Duration(attempt*20+rand.IntN(30)) *
			time.Millisecond

		select {
		case <-ctx.Done():
			return err
		case <-time.After(wait):
		}
	}

	return err
}

// runOnce menjalankan satu kali percobaan transaksi.
func (m *TxManager) runOnce(
	ctx context.Context,
	fn func(tx DBTX) error,
) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("memulai transaction: %w", err)
	}

	// Rollback selalu dijalankan sebagai pengaman.
	// Setelah Commit berhasil, Rollback tidak mengubah data.
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// isRetryable menentukan apakah error bisa dicoba kembali.
func isRetryable(err error) bool {
	code := pgErrorCode(err)

	return code == codeSerializationFailure ||
		code == codeDeadlockDetected
}
