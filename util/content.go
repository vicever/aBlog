package util

func WordCount(content string) int64 {
	// todo: article words count
	return int64(len(content))
}

func ReadingTimeCount(length int64) int64 {
	// todo : article reading time
	return length / 5
}
