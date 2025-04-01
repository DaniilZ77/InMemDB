package replication

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/DaniilZ77/InMemDB/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRequest(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	disk := NewMockNextSegmentProvider(t)
	master, err := NewMaster(disk, dir, slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	firstSegment, secondSegment, testData := "first_segment.log", "second_segment.log", "testdata"
	data, err := common.Encode(NewRequest("first_segment.log"))
	require.NoError(t, err)

	disk.EXPECT().NextSegment(firstSegment).
		Return(secondSegment, nil).Once()

	err = os.WriteFile(dir+"/"+secondSegment, []byte(testData), 0666)
	require.NoError(t, err)

	response, err := master.HandleRequest(data)
	require.NoError(t, err)

	decodedResponse, err := common.DecodeOne[Response](response)
	require.NoError(t, err)

	assert.True(t, decodedResponse.Ok)
	assert.Equal(t, []byte(testData), decodedResponse.Segment)
	assert.Equal(t, secondSegment, decodedResponse.Filename)
}

func TestHandleRequest_Error(t *testing.T) {
	disk := NewMockNextSegmentProvider(t)
	master, err := NewMaster(disk, "testdir", slog.New(slog.NewJSONHandler(io.Discard, nil)))
	require.NoError(t, err)

	okRequest, err := common.Encode(NewRequest("segment.log"))
	require.NoError(t, err)

	tests := []struct {
		name    string
		request []byte
		mock    func()
	}{
		{
			name:    "empty request",
			request: []byte(nil),
			mock:    func() {},
		},
		{
			name:    "invalid request",
			request: []byte("invalid"),
			mock:    func() {},
		},
		{
			name:    "next segment error",
			request: okRequest,
			mock: func() {
				disk.EXPECT().NextSegment("segment.log").Return("", errors.New("failed to get next segment")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			response, err := master.HandleRequest(tt.request)
			require.NoError(t, err)

			decodedResponse, err := common.DecodeOne[Response](response)
			require.NoError(t, err)

			assert.False(t, decodedResponse.Ok)
		})
	}
}
