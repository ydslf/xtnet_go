package util

func IsPow2(size uint32) bool {
	return (size & (size - 1)) == 0
}

func SizeOfPow2(size uint32) uint32 {
	if IsPow2(size) {
		return size
	}
	size = size - 1
	size = size | (size >> 1)
	size = size | (size >> 2)
	size = size | (size >> 4)
	size = size | (size >> 8)
	size = size | (size >> 16)
	return size + 1
}
