package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
	"google.golang.org/genai"
)

type ComparisonService struct {
	client *genai.Client
	repo   *repository.ComparisonRepository
}

func NewComparisonService(client *genai.Client) *ComparisonService {
	return &ComparisonService{
		client: client,
		repo:   repository.NewComparisonRepository(),
	}
}

type ComparisonResult struct {
	Winner       string   `json:"winner"`
	WinnerPoints []string `json:"winner_points"` // 4 points on why the winner won
	LoserPoints  []string `json:"loser_points"`  // 4 points on what the loser needs to improve
}

func (s *ComparisonService) CompareUsers(ctx context.Context, user1Name, user1Data, user2Name, user2Data string) (*ComparisonResult, error) {
	// 1. Check Cache
	if existing, err := s.repo.GetValidComparison(user1Name, user2Name, 10*time.Hour); err == nil && existing != nil {
		var cachedResult ComparisonResult
		if err := json.Unmarshal([]byte(existing.Result), &cachedResult); err == nil {
			log.Printf("Returning cached comparison for %s vs %s", user1Name, user2Name)
			return &cachedResult, nil
		}
	}

	prompt := fmt.Sprintf(`
You are an expert judge of coding profiles, specifically focusing on LeetCode data. Compare the following two users: "%s" and "%s" based on their provided data.
Determine the winner based on their achievements, skills, and overall profile quality.
Provide exactly 4 distinct points explaining why the winner won. Use their actual name in the points.
Provide exactly 4 distinct points explaining what the loser needs to improve to reach the winner's level. Use their actual name in the points.

Return the result in strictly valid JSON format with the following keys:
- winner: "%s" or "%s" (or "tie" if indistinguishable)
- winner_points: ["point 1", "point 2", "point 3", "point 4"]
- loser_points: ["point 1", "point 2", "point 3", "point 4"]

%s Data:
%s

%s Data:
%s
`, user1Name, user2Name, user1Name, user2Name, user1Name, user1Data, user2Name, user2Data)

	resp, err := s.client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	jsonStr := resp.Text()
	// Clean up potential markdown code blocks if any (though MIME type should handle it)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")

	var result ComparisonResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("Failed to unmarshal JSON: %s", jsonStr)
		return nil, fmt.Errorf("failed to parse comparison result: %w", err)
	}

	// Save to Cache
	if err := s.repo.SaveComparison(user1Name, user2Name, result); err != nil {
		log.Printf("Failed to save comparison to cache: %v", err)
	}

	return &result, nil
}
