package disk

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearDir(t *testing.T, dir string) {
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	for i := range entries {
		err = os.Remove(filepath.Join(dir, entries[i].Name()))
		require.NoError(t, err)
	}
}

func TestDisk(t *testing.T) {
	dir := t.TempDir()
	const maxSegmentSize = 1000
	testData := "testdata"

	disk := NewDisk(dir, maxSegmentSize, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	t.Run("write segment overflow", func(t *testing.T) {
		t.Cleanup(func() { clearDir(t, dir) })
		iterationsNumber := maxSegmentSize/len(testData) + 2
		for range iterationsNumber {
			err := disk.Write([]byte(testData))
			require.NoError(t, err)
		}

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)

		assert.Len(t, entries, 2)
	})

	t.Run("read", func(t *testing.T) {
		t.Cleanup(func() { clearDir(t, dir) })

		createLogFile := func(index int) {
			file, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("testfile%d.log", index)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			require.NoError(t, err)

			_, err = file.Write([]byte(testData))
			require.NoError(t, err)
		}

		createLogFile(1)
		createLogFile(2)
		createLogFile(3)

		data, err := disk.Read()
		require.NoError(t, err)

		assert.Equal(t, testData+testData+testData, string(data))
	})
}
