package models

type AdRequest struct {
	AdPlacementID string `json:"adPlacementId,omitempty"`
}

type AdObject struct {
	AdID     string  `json:"adId"`
	BidPrice float64 `json:"bidPrice"`
}

type BidResponse struct {
	Status   int       `json:"status"`
	AdObject *AdObject `json:"adObject,omitempty"`
}
