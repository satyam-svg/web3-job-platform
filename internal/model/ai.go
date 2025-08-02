package model

type JobRecommendation struct {
	Title         string `json:"title"`
	Company       string `json:"company"`
	MatchingScore int    `json:"matching_score"`
	Reasoning     string `json:"reasoning"`
	Recommended   bool   `json:"recommended"`
}

type JobRecommendationResponse struct {
	Recommendations []JobRecommendation `json:"recommendations"`
}
