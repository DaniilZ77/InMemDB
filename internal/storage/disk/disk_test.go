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

func createLogFiles(t *testing.T, dir string, filesNumber int, data []string) []string {
	if len(data) == 0 {
		data = []string{"testdata"}
	}
	files := []string{}
	for i := range filesNumber {
		time.Sleep(5 * time.Millisecond)
		files = append(files, fmt.Sprintf("testfile%d.log", time.Now().UnixMilli()))
		err := os.WriteFile(filepath.Join(dir, files[i]), []byte(data[i%len(data)]), 0666)
		require.NoError(t, err)
	}

	return files
}

func TestWriteSegment_WithOverflow(t *testing.T) {
	dir := t.TempDir()
	const maxSegmentSize = 1000
	disk := NewDisk(dir, maxSegmentSize, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	testData := "testdata"
	iterationsNumber := maxSegmentSize/len(testData) + 2

	for range iterationsNumber {
		err := disk.WriteSegment([]byte(testData))
		require.NoError(t, err)
	}

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 2)
}

func TestReadSegments(t *testing.T) {
	dir := t.TempDir()
	disk := NewDisk(dir, 1000, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	testData1 := "lorem ipsum 1"
	testData2 := "lorem ipsum 2"
	testData3 := "lorem ipsum 3"
	_ = createLogFiles(t, dir, 3, []string{testData1, testData2, testData3})

	data, err := disk.ReadSegments()
	require.NoError(t, err)
	assert.Equal(t, testData1+testData2+testData3, string(data))
}

func TestNextSegment(t *testing.T) {
	dir := t.TempDir()
	disk := NewDisk(dir, 1000, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	segments := createLogFiles(t, dir, 3, nil)

	nextSegment, err := disk.NextSegment("")
	require.NoError(t, err)
	assert.Equal(t, segments[0], nextSegment)

	nextSegment, err = disk.NextSegment(segments[0])
	require.NoError(t, err)
	assert.Equal(t, segments[1], nextSegment)

	nextSegment, err = disk.NextSegment(segments[1])
	require.NoError(t, err)
	assert.Equal(t, "", nextSegment)

	nextSegment, err = disk.NextSegment(segments[2])
	require.NoError(t, err)
	assert.Equal(t, "", nextSegment)
}

func TestSegments_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	disk := NewDisk(dir, 1000, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	tests := []struct {
		name        string
		expectedErr error
		call        func() (string, error)
	}{
		{
			name: "next segment empty dir",
			call: func() (string, error) {
				return disk.NextSegment("")
			},
		},
		{
			name: "last segment empty dir",
			call: func() (string, error) {
				return disk.LastSegment()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.call()
			assert.NoError(t, err)
			assert.Empty(t, res)
		})
	}
}

func TestLastSegment(t *testing.T) {
	dir := t.TempDir()
	disk := NewDisk(dir, 1000, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	segments := createLogFiles(t, dir, 3, nil)

	segment, err := disk.LastSegment()
	require.NoError(t, err)
	assert.Equal(t, segments[2], segment)
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	disk := NewDisk(dir, 1000, slog.New(slog.NewJSONHandler(io.Discard, nil)))

	testdata := "lorem ipsum"
	err := disk.WriteFile("newfile.log", []byte(testdata))
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "/newfile.log"))
	require.NoError(t, err)

	assert.Equal(t, testdata, string(data))
}
