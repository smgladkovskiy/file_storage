package file

import "math"

func Split(file []byte, parts int) [][]byte {
	totalBytes := len(file)
	if totalBytes == 0 {
		return nil
	}

	chunkSize := int(math.Ceil(float64(totalBytes) / float64(parts)))
	chunks := make([][]byte, 0, parts)
	part := 0
	from := 0
	to := chunkSize
	for part < parts {
		if part > 0 {
			to = chunkSize * (part + 1)
			from = chunkSize * (part)
		}

		chunk := make([]byte, chunkSize)

		if part == parts-1 && totalBytes-from > 0 {
			to = totalBytes
			chunk = make([]byte, totalBytes-from)
		}

		if to > totalBytes {
			chunk = nil
		} else {
			copy(chunk, file[from:to])
		}

		chunks = append(chunks, chunk)
		part++
	}

	return chunks
}
