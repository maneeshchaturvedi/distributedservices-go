package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var encoding = binary.BigEndian

const recWidth = 8

type store struct {
	file   *os.File
	mutex  sync.Mutex
	writer *bufio.Writer
	size   uint64
}

func newStore(file *os.File) (*store, error) {
	fi, err := os.Stat(file.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	return &store{
		file:   file,
		mutex:  sync.Mutex{},
		writer: bufio.NewWriter(file),
		size:   size,
	}, nil
}

func (s *store) Append(data []byte) (n int, pos uint64, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Get the current position
	pos = s.size
	if err := binary.Write(s.writer, encoding, uint64(len(data))); err != nil {
		return 0, 0, err
	}

	// Write the data
	w, err := s.writer.Write(data)
	if err != nil {
		return 0, 0, err
	}

	w += recWidth
	// Update the size
	s.size += uint64(w)

	return n, pos, nil
}

func (s *store) Read(position uint64) ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.writer.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, recWidth)
	if _, err := s.file.ReadAt(size, int64(position)); err != nil {
		return nil, err
	}

	data := make([]byte, encoding.Uint64(size))

	// Read the data
	if _, err := s.file.ReadAt(data, int64(position+recWidth)); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *store) ReadAt(data []byte, offset int64) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.writer.Flush(); err != nil {
		return 0, err
	}

	return s.file.ReadAt(data, offset)
}

func (s *store) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.writer.Flush(); err != nil {
		return err
	}

	if err := s.file.Close(); err != nil {
		return err
	}

	return nil
}
