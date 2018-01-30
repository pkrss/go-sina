package kline

type K_MACDV struct {
	Index int64   `json:"i"`
	Time  int64   `json:"t"`
	High  float64 `json:"h"`
	Open  float64 `json:"o"`
	Low   float64 `json:"l"`
	Close float64 `json:"c"`
	Vol   int64   `json:"v"`
}