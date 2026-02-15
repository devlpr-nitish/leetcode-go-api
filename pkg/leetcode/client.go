package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type UserProfile struct {
	Username       string  `json:"username"`
	RealName       string  `json:"realName"`
	Ranking        int     `json:"ranking"`
	Reputation     int     `json:"reputation"`
	TotalSolved    int     `json:"totalSolved"`
	EasySolved     int     `json:"easySolved"`
	MediumSolved   int     `json:"mediumSolved"`
	HardSolved     int     `json:"hardSolved"`
	AcceptanceRate float64 `json:"acceptanceRate"`
	Contribution   int     `json:"contributionPoints"`
	GlobalRanking  int     `json:"globalRanking"`

	// Complex fields like submission stats can be added here
}

// Simplified GraphQL response structure
type graphQLResponse struct {
	Data struct {
		MatchedUser struct {
			Username string `json:"username"`
			Profile  struct {
				RealName   string `json:"realName"`
				Reputation int    `json:"reputation"`
				Ranking    int    `json:"ranking"`
			} `json:"profile"`
			SubmitStats struct {
				AcSubmissionNum []struct {
					Difficulty  string `json:"difficulty"`
					Count       int    `json:"count"`
					Submissions int    `json:"submissions"`
				} `json:"acSubmissionNum"`
			} `json:"submitStats"`
		} `json:"matchedUser"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) GetUserProfile(username string) (*UserProfile, error) {
	query := `
	query getUserProfile($username: String!) {
		matchedUser(username: $username) {
			username
			profile {
				realName
				reputation
				ranking
			}
			submitStats {
				acSubmissionNum {
					difficulty
					count
					submissions
				}
			}
		}
	}
	`

	reqBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": map[string]string{"username": username},
	})

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("leetcode api returned status: %d", resp.StatusCode)
	}

	var parsedResp graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return nil, err
	}

	if len(parsedResp.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", parsedResp.Errors[0].Message)
	}

	user := parsedResp.Data.MatchedUser
	profile := &UserProfile{
		Username:      user.Username,
		RealName:      user.Profile.RealName,
		Reputation:    user.Profile.Reputation,
		GlobalRanking: user.Profile.Ranking,
	}

	for _, stat := range user.SubmitStats.AcSubmissionNum {
		switch stat.Difficulty {
		case "All":
			profile.TotalSolved = stat.Count
		case "Easy":
			profile.EasySolved = stat.Count
		case "Medium":
			profile.MediumSolved = stat.Count
		case "Hard":
			profile.HardSolved = stat.Count
		}
	}

	return profile, nil
}

// External API Response Structures
type ExternalProblemResponse struct {
	TotalQuestions         int           `json:"totalQuestions"`
	Count                  int           `json:"count"`
	ProblemsetQuestionList []APIQuestion `json:"problemsetQuestionList"`
}

type APIQuestion struct {
	AcRate             float64    `json:"acRate"`
	Difficulty         string     `json:"difficulty"`
	FreqBar            float64    `json:"freqBar"`
	QuestionFrontendId string     `json:"questionFrontendId"`
	Title              string     `json:"title"`
	TitleSlug          string     `json:"titleSlug"`
	TopicTags          []TopicTag `json:"topicTags"`
}

type TopicTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	ID   string `json:"id"`
}

// GetProblems fetches problems from the external API
func (c *Client) GetProblems(limit int, skip int, tags []string, difficulty string) ([]APIQuestion, error) {
	// Constuct URL: https://leetcode-api-v8xt.onrender.com/problems
	baseURL := "https://leetcode-api-v8xt.onrender.com/problems"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	if skip > 0 {
		q.Add("skip", fmt.Sprintf("%d", skip))
	}
	if difficulty != "" {
		q.Add("difficulty", difficulty)
	}
	if len(tags) > 0 {
		// API format: tags=tag1+tag2
		// We'll join them space separated which encodes to + usually, or custom generic join
		// The API doc says "tags=tag1+tag2". Query encoding usually deals with this.
		// Let's loop and add multiple 'tags' keys? No, commonly it's comma or space.
		// User provided example: tags=tag1+tag2
		// This likely means a single query param "tags" with value "tag1 tag2" (space encoded to +)
		q.Add("tags", fmt.Sprintf("%s", joinTags(tags)))
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external api error: %d", resp.StatusCode)
	}

	var parsedResp ExternalProblemResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return nil, err
	}

	return parsedResp.ProblemsetQuestionList, nil
}

func joinTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	// Simple join with space, URL encoding will handle the rest
	res := tags[0]
	for i := 1; i < len(tags); i++ {
		res += " " + tags[i]
	}
	return res
}
