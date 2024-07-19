package main

import "io"

// BytesFromReader returns the bytes from a reader interface into a byte slice.
// The bytes are processed in numChunks chunks, where the chunks are chunkSize large.
// If numChunks is negative, the reader is processed up until the end of the file
func BytesFromReader(reader io.Reader, chunkSize int, numChunks int) []byte {
	result := make([]byte, 0)
	chunk := make([]byte, chunkSize)
	count := 0
	var err error
	var n int
	for err != io.EOF {
		if numChunks >= 0 && count >= numChunks {
			break
		}
		n, err = reader.Read(chunk)
		result = append(result, chunk[:n]...)
		count++
	}

	return result
}

// BytesToString accepts a byte slice and returns a string
func BytesToString(bytes []byte) string {
	result := ""
	for _, v := range bytes {
		result += string(v)
	}
	return result
}
