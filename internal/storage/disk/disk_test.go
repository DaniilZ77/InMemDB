package disk

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	disk := NewDisk(dir, maxSegmentSize, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	t.Run("write segment overflow", func(t *testing.T) {
		t.Cleanup(func() { clearDir(t, dir) })
		testData := "testdata"
		iterationsNumber := maxSegmentSize/len(testData) + 2
		for range iterationsNumber {
			err := disk.WriteSegment([]byte(testData))
			require.NoError(t, err)
		}

		entries, err := os.ReadDir(dir)
		require.NoError(t, err)
		assert.Len(t, entries, 2)
	})

	t.Run("read", func(t *testing.T) {
		t.Cleanup(func() { clearDir(t, dir) })
		testData1 := "lorem ipsum 1"
		testData2 := "lorem ipsum 2"
		testData3 := "lorem ipsum 3"
		createLogFile := func(data string) {
			file, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("testfile%d.log", time.Now().UnixMilli())), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			require.NoError(t, err)

			_, err = file.Write([]byte(data))
			require.NoError(t, err)
		}

		createLogFile(testData1)
		createLogFile(testData2)
		createLogFile(testData3)

		data, err := disk.ReadSegments()
		require.NoError(t, err)
		assert.Equal(t, testData1+testData2+testData3, string(data))
	})
}
