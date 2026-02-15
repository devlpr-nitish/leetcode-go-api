package leetcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	baseURL := os.Getenv("LEETCODE_API_URL")
	if baseURL == "" {
		baseURL = "https://leetcode-api-v8xt.onrender.com" // Default fallback
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetUserProfile(username string) (*ProfileResponse, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, username) // Assuming endpoint is /<username>
	var resp ProfileResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	// The API might return website as a list of strings
	return &resp, nil
}

func (c *Client) GetUserStats(username string) (*StatsResponse, error) {
	url := fmt.Sprintf("%s/%s/profile", c.BaseURL, username) // Endpoint /<username>/profile
	var resp StatsResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetUserSkills(username string) (*SkillsResponse, error) {
	url := fmt.Sprintf("%s/%s/skill", c.BaseURL, username)
	var resp SkillsResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetUserContest(username string) (*ContestResponse, error) {
	url := fmt.Sprintf("%s/%s/contest", c.BaseURL, username)
	var resp ContestResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GetUserCalendar(username string) (*CalendarResponse, error) {
	url := fmt.Sprintf("%s/%s/calendar", c.BaseURL, username)
	var resp CalendarResponse
	if err := c.fetch(url, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) fetch(url string, target interface{}) error {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code %d for url %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response from %s: %w", url, err)
	}
	return nil
}

// FetchAllUserData fetches all available data for a user concurrently
type AllUserData struct {
	Profile  *ProfileResponse
	Stats    *StatsResponse
	Skills   *SkillsResponse
	Contest  *ContestResponse
	Calendar *CalendarResponse
	Errors   []error
}

func (c *Client) FetchAllUserData(username string) *AllUserData {
	var wg sync.WaitGroup
	result := &AllUserData{}
	var mu sync.Mutex

	fetchops := []func(){
		func() {
			res, err := c.GetUserProfile(username)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("profile: %w", err))
			} else {
				result.Profile = res
			}
		},
		func() {
			res, err := c.GetUserStats(username)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("stats: %w", err))
			} else {
				result.Stats = res
			}
		},
		func() {
			res, err := c.GetUserSkills(username)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("skills: %w", err))
			} else {
				result.Skills = res
			}
		},
		func() {
			res, err := c.GetUserContest(username)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("contest: %w", err))
			} else {
				result.Contest = res
			}
		},
		func() {
			res, err := c.GetUserCalendar(username)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("calendar: %w", err))
			} else {
				result.Calendar = res
			}
		},
	}

	wg.Add(len(fetchops))
	for _, op := range fetchops {
		go func(o func()) {
			defer wg.Done()
			o()
		}(op)
	}
	wg.Wait()

	return result
}
