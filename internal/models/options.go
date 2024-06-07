package models

type Options struct {
	AnswerFilename         string `json:"answer_filename"`
	DeltaTimestamp         int    `json:"delta"`
	GoroutinesCount        int    `json:"goroutines_count"`
	TimestampsPerGoroutine int    `json:"timestamps_per_goroutine"`
}
