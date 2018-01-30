package minute

type MINUTEDATA struct {
	Index     int64   `json:"i"`
	Time      int64   `json:"t"`
	Price     float64 `json:"p"`
	Avg       float64 `json:"a"`
	Vol       int64   `json:"v"`
	UpDn      float64 `json:"updn"`
	UpDnPer   float64 `json:"updnPer"`
	Positions int64   `json:"po"`
}
