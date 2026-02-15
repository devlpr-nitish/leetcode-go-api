package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/clients/leetcode"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
	"gorm.io/datatypes"
)

type UserService struct {
	UserRepo       repository.UserRepository
	LeetCodeClient *leetcode.Client
}

func NewUserService(userRepo repository.UserRepository, client *leetcode.Client) *UserService {
	return &UserService{
		UserRepo:       userRepo,
		LeetCodeClient: client,
	}
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.UserRepo.GetByUsername(ctx, username)
}

func (s *UserService) SyncUser(ctx context.Context, username string) (*models.User, error) {
	// 1. Fetch all data from LeetCode concurrently
	data := s.LeetCodeClient.FetchAllUserData(username)

	// Check for critical errors (e.g. if everything failed)
	// For now, we proceed even if some parts failed, but we should log errors.
	if len(data.Errors) > 0 {
		// In a real app, log these errors.
		// fmt.Println("Errors fetching data:", data.Errors)
		// If profile is missing, we probably can't do much.
		if data.Profile == nil && data.Stats == nil {
			return nil, fmt.Errorf("failed to fetch essential user data: %v", data.Errors)
		}
	}

	// 2. Find or Create User in DB
	user, err := s.UserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// If strictly syncing, maybe we shouldn't create if it doesn't exist?
		// But the requirement says "if user request for enpoint to sync leetcode current data then i will fetch all data... and update my db"
		// It implies the user might already exist. If not, we can create.
		user = &models.User{
			Username:  username,
			CreatedAt: time.Now(),
		}
	}

	// 3. Update User fields
	// Basic Profile
	if data.Profile != nil {
		user.Name = data.Profile.Name
		user.Avatar = data.Profile.Avatar
		user.About = data.Profile.About
		// Birthday parsing if needed, assumed string or handle formatting
		// user.Birthday = ...
		user.Country = data.Profile.Country
		user.School = data.Profile.School
		user.Company = data.Profile.Company
		user.GitHub = data.Profile.GitHub
		user.Twitter = data.Profile.Twitter
		user.LinkedIn = data.Profile.LinkedIn

		websitesJSON, _ := json.Marshal(data.Profile.Website)
		user.Websites = datatypes.JSON(websitesJSON)
	}

	// Stats
	if data.Stats != nil {
		user.Ranking = data.Stats.Ranking
		user.Reputation = data.Stats.Reputation
		user.ContributionPoint = data.Stats.ContributionPoint
		user.TotalSolved = data.Stats.TotalSolved
		user.EasySolved = data.Stats.EasySolved
		user.MediumSolved = data.Stats.MediumSolved
		user.HardSolved = data.Stats.HardSolved
	}

	// Skills
	if data.Skills != nil {
		skillTagsJSON, _ := json.Marshal(data.Skills.Data.MatchedUser.TagProblemCounts)
		user.SkillTags = datatypes.JSON(skillTagsJSON)
	}

	// Contest
	if data.Contest != nil {
		user.ContestRating = data.Contest.ContestRating
		user.ContestGlobalRanking = data.Contest.ContestGlobalRanking
		user.ContestTopPercentage = data.Contest.ContestTopPercentage
		user.TotalParticipants = data.Contest.TotalParticipants
		user.ContestAttended = data.Contest.ContestAttended

		badgesJSON, _ := json.Marshal(data.Contest.Badges)
		user.Badges = datatypes.JSON(badgesJSON)
	}

	// Calendar
	if data.Calendar != nil {
		user.Streak = data.Calendar.Streak
		user.TotalActiveDays = data.Calendar.TotalActiveDays

		activeYearsJSON, _ := json.Marshal(data.Calendar.ActiveYears)
		user.ActiveYears = datatypes.JSON(activeYearsJSON)

		// The user warned: "WARNING: This can be large". But requested to store it.
		// Assuming we just store the raw JSON string or object provided by API.
		// The API model has SubmissionCalendar as string (JSON string).
		// We can store it directly if it's already a JSON string, or marshal it.
		// In client/models.go I defined it as string.
		// If it's a JSON string, we can cast it to datatypes.JSON IF it is valid JSON.
		if data.Calendar.SubmissionCalendar != "" {
			user.SubmissionCalendar = datatypes.JSON(data.Calendar.SubmissionCalendar)
		}
	}

	// Calculate App-Specific Score (Placeholder logic)
	// TotalScore logic wasn't specified, so I'll leave it or basic sum
	user.TotalScore = user.TotalSolved * 10
	user.ScoreRank = "Beginner" // Placeholder
	if user.TotalSolved > 100 {
		user.ScoreRank = "Intermediate"
	}
	if user.TotalSolved > 500 {
		user.ScoreRank = "Advanced"
	}

	user.UpdatedAt = time.Now()

	// 4. Save to DB
	if user.ID == (models.User{}.ID) { // Check if new user (empty UUID or nil-like)
		// Wait, UUID is struct.
		// Actually, if we fetched by username and it was nil, we created a new struct.
		// Its ID will be empty.
		// However, standard flow might be: User should Register first.
		// But Sync might be called for existing user.
		// If user doesn't exist in DB, Create.
		if err := s.UserRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		if err := s.UserRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
