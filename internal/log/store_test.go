package log

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var (
	write = []byte("test data")
	width = uint64(len(write)) + recWidth
)

func TestStoreAppendRead(t *testing.T) {

	file, err := ioutil.TempFile("", "store_append_read_test")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	s, err := newStore(file)
	require.NoError(t, err)

	testAppend(t, s)
	testRead(t, s)
	testReadAt(t, s)
	s, err = newStore(file)
	require.NoError(t, err)
	testRead(t, s)
}

func testAppend(t *testing.T, s *store) {

	t.Helper()
	for i := uint64(1); i < 4; i++ {
		n, pos, err := s.Append(write)
		require.NoError(t, err)
		require.Equal(t, int(pos)+n, int(width*i))
	}

}

func testRead(t *testing.T, s *store) {
	t.Helper()
	pos := uint64(0) // Declare and initialize pos variable
	for i := uint64(1); i < 4; i++ {
		readData, err := s.Read(uint64(i) * width)
		require.NoError(t, err)
		require.Equal(t, readData, write)
		pos += uint64(width) // Increment pos inside the loop
	}
}

func testReadAt(t *testing.T, s *store) {
	t.Helper()
	for i, offset := uint64(1), int64(0); i < 4; i++ {
		readData := make([]byte, len(write))
		n, err := s.ReadAt(readData, offset)
		require.NoError(t, err)
		require.Equal(t, recWidth, n)
		offset += int64(n) // Convert n to int64
		size := encoding.Uint64(readData)
		readData = make([]byte, size)
		n, err = s.ReadAt(readData, offset)
		require.NoError(t, err)
		require.Equal(t, write, readData)
		offset += int64(n) // Convert n to int64
	}
}
