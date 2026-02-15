package leetcode

// ProfileResponse represents the response from GET /<username>
type ProfileResponse struct {
	Name     string   `json:"name"`
	Avatar   string   `json:"avatar"`
	About    string   `json:"about"`
	Birthday string   `json:"birthday"` // API might return string or null
	Country  string   `json:"country"`
	School   string   `json:"school"`
	Company  string   `json:"company"`
	GitHub   string   `json:"gitHub"`
	Twitter  string   `json:"twitter"`
	LinkedIn string   `json:"linkedIn"`
	Website  []string `json:"website"`
}

// StatsResponse represents the response from GET /<username>/profile
// Note: The user mentioned "getUserStats", but the endpoint is /<username>/profile in their table.
// Wait, the table says:
// Stats -> getUserStats -> GET /<username>/profile
// But Profile -> getUserProfile -> GET /<username>
// This seems inverted or specific to the wrapper they are using.
// I will follow the user's table mapping for the struct names.
type StatsResponse struct {
	Ranking           int `json:"ranking"`
	Reputation        int `json:"reputation"`
	ContributionPoint int `json:"contributionPoint"`
	TotalSolved       int `json:"totalSolved"`
	EasySolved        int `json:"easySolved"`
	MediumSolved      int `json:"mediumSolved"`
	HardSolved        int `json:"hardSolved"`
}

// SkillStats represents the skill breakdown
type SkillStats struct {
	TagName        string `json:"tagName"`
	TagSlug        string `json:"tagSlug"`
	ProblemsSolved int    `json:"problemsSolved"`
}

// SkillsResponse represents the response from GET /<username>/skill
type SkillsResponse struct {
	Data struct {
		MatchedUser struct {
			TagProblemCounts struct {
				Advanced     []SkillStats `json:"advanced"`
				Intermediate []SkillStats `json:"intermediate"`
				Fundamental  []SkillStats `json:"fundamental"`
			} `json:"tagProblemCounts"`
		} `json:"matchedUser"`
	} `json:"data"`
}

// ContestResponse represents the response from GET /<username>/contest
type ContestResponse struct {
	ContestRating        float64 `json:"contestRating"`
	ContestGlobalRanking int     `json:"contestGlobalRanking"`
	ContestTopPercentage float64 `json:"contestTopPercentage"`
	TotalParticipants    int     `json:"totalParticipants"`
	ContestAttended      int     `json:"contestAttended"`
	Badges               []struct {
		Name string `json:"name"`
		Icon string `json:"icon"` // Url to icon
	} `json:"badges"`
}

// CalendarResponse represents the response from GET /<username>/calendar
type CalendarResponse struct {
	Streak             int    `json:"streak"`
	TotalActiveDays    int    `json:"totalActiveDays"`
	ActiveYears        []int  `json:"activeYears"`
	SubmissionCalendar string `json:"submissionCalendar"` // JSON string
}
