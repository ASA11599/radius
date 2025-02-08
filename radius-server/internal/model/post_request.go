package model

type PostRequest struct {
	Location Location `json:"location"`
	Content string `json:"content"`
	Duration int64 `json:"duration"`
}

func (pr *PostRequest) Valid() bool {
	return pr.Location.Valid() && (len(pr.Content) > 0) && ((pr.Duration < 3600) && (pr.Duration > 0))
}
