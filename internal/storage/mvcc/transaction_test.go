package mvcc

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	engine := NewMockEngine(t)
	wal := NewMockWal(t)
	coordinator, err := NewCoordinator(engine, wal, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	var beginTxID int64
	engine.EXPECT().Get(mock.Anything, "age").Return("20", true).Once()
	engine.EXPECT().SetMany(mock.Anything, mock.Anything).Return().Once()
	engine.EXPECT().ExistsBetween(mock.MatchedBy(func(id int64) bool {
		beginTxID = id
		return true
	}), mock.MatchedBy(func(id int64) bool {
		return beginTxID < id
	}), mock.Anything).Return(false)
	wal.EXPECT().Save(mock.Anything, mock.Anything).Return(true)

	tx := coordinator.BeginTransaction()
	err = tx.Set("name", "Daniil")
	assert.NoError(t, err)
	value, ok := tx.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "Daniil", value)
	value, ok = tx.Get("age")
	assert.True(t, ok)
	assert.Equal(t, "20", value)
	err = tx.Del("name")
	assert.NoError(t, err)
	_, ok = tx.Get("name")
	assert.False(t, ok)
	err = tx.Commit()
	assert.NoError(t, err)
}

func TestTransactionInterfered(t *testing.T) {
	engine := NewMockEngine(t)
	wal := NewMockWal(t)
	coordinator, err := NewCoordinator(engine, wal, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	engine.EXPECT().ExistsBetween(mock.Anything, mock.Anything, mock.Anything).Return(true).Once()
	wal.EXPECT().Save(mock.Anything, mock.Anything).Return(true).Once()

	tx := coordinator.BeginTransaction()
	err = tx.Set("name", "Daniil")
	assert.NoError(t, err)
	err = tx.Commit()
	assert.ErrorIs(t, err, ErrTransactionInterfered)
}

func TestTransactionRollback(t *testing.T) {
	engine := NewMockEngine(t)
	wal := NewMockWal(t)
	coordinator, err := NewCoordinator(engine, wal, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	tx := coordinator.BeginTransaction()
	err = tx.Rollback()
	assert.NoError(t, err)
}

func TestTransactionTwiceEnd(t *testing.T) {
	engine := NewMockEngine(t)
	wal := NewMockWal(t)
	coordinator, err := NewCoordinator(engine, wal, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	wal.EXPECT().Save(mock.Anything, mock.Anything).Return(true).Once()
	engine.EXPECT().SetMany(mock.Anything, mock.Anything).Return().Once()

	tx := coordinator.BeginTransaction()
	err = tx.Commit()
	assert.NoError(t, err)
	err = tx.Rollback()
	assert.ErrorIs(t, err, ErrTransactionFinished)
}
