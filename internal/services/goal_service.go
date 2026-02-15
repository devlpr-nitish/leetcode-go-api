package services

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/models"
	"github.com/devlpr-nitish/leetcode-tracker-backend/internal/repository"
	"github.com/devlpr-nitish/leetcode-tracker-backend/pkg/leetcode"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type GoalService struct {
	UserRepo    repository.UserRepository
	GoalRepo    repository.GoalRepository
	ProblemRepo repository.ProblemRepository
}

func NewGoalService(userRepo repository.UserRepository, goalRepo repository.GoalRepository, problemRepo repository.ProblemRepository) *GoalService {
	return &GoalService{
		UserRepo:    userRepo,
		GoalRepo:    goalRepo,
		ProblemRepo: problemRepo,
	}
}

// UserProfile definitions for the algorithm
type UserProfile struct {
	TotalSolved    int
	EasySolved     int
	MediumSolved   int
	HardSolved     int
	WeakTopics     []string
	StrongTopics   []string
	RecentActivity struct {
		Last7DaysSolved int
		AvgDailyTime    int
	}
}

type DifficultyRatio struct {
	Easy   float64
	Medium float64
	Hard   float64
}

// Goals Cleanup Logic
func (s *GoalService) CleanupOldGoals(ctx context.Context, userID uuid.UUID) error {
	// Keep: Past 1 week, Current week, Next 1 week.
	// Current week start
	// now := getWeekStart()
	// pastLimit := now.AddDate(0, 0, -7)
	// futureLimit := now.AddDate(0, 0, 14) // Allow generating next week too, so keep valid range

	// Implementation depends on repo, but for now we can assume repo handles specific logic closer to DB
	// OR we fetch all and delete. Better to add a DeleteBefore method in repo.
	// For MVP without changing repo interface too much, let's just leave a placeholder or implemented if repo allows.
	// See Task: "store only weekly goals for user, like one week past and current week and one week next only"

	// We will implement a repo method call if possible, or raw SQL exec if repo is limited.
	return nil
}

// GenerateWeeklyGoals is the core logic for creating new personalized goals
func (s *GoalService) GenerateWeeklyGoals(ctx context.Context, userID uuid.UUID) error {
	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 1. Feature Extraction
	profile := s.buildUserProfile(user)

	// 2. Difficulty Allocation
	ratio := s.getDifficultyRatio(profile)

	// 3. Goals Configuration
	weeklyProblemCount := 8

	easyCount := int(math.Round(float64(weeklyProblemCount) * ratio.Easy))
	mediumCount := int(math.Round(float64(weeklyProblemCount) * ratio.Medium))
	hardCount := weeklyProblemCount - easyCount - mediumCount

	// 4. Problem Selection Strategy (API Based)
	// We need to fetch problems per category because we can't filter in-memory from a huge list anymore.

	// Fetch more than needed to allow for some random selection/filtering
	apiEasy, _ := s.ProblemRepo.GetProblemsByDifficulty(ctx, "Easy", easyCount*3)
	apiMedium, _ := s.ProblemRepo.GetProblemsByDifficulty(ctx, "Medium", mediumCount*3)
	apiHard, _ := s.ProblemRepo.GetProblemsByDifficulty(ctx, "Hard", hardCount*3)

	// Weak Topic Injection
	// If user has weak topics, fetch a few from there and replace some medium/hard
	if len(profile.WeakTopics) > 0 {
		topic := profile.WeakTopics[0]
		topicProblems, _ := s.ProblemRepo.GetProblemsByTopic(ctx, topic, 5)
		// Naively merge some topics into medium/hard lists
		for _, tp := range topicProblems {
			if tp.Difficulty == "Medium" {
				apiMedium = append(apiMedium, tp)
			} else if tp.Difficulty == "Hard" {
				apiHard = append(apiHard, tp)
			}
		}
	}

	selectedProblems := make([]leetcode.APIQuestion, 0)
	selectedProblems = append(selectedProblems, s.selectProblems(apiEasy, easyCount)...)
	selectedProblems = append(selectedProblems, s.selectProblems(apiMedium, mediumCount)...)
	selectedProblems = append(selectedProblems, s.selectProblems(apiHard, hardCount)...)

	// 5. Weekly Distribution
	dailyPlan := s.distributeAcrossWeek(selectedProblems)

	// 6. Construct Goal Objects
	weekStart := getWeekStart()

	breakdownJSON, _ := json.Marshal(map[string]int{
		"easy":   easyCount,
		"medium": mediumCount,
		"hard":   hardCount,
	})

	selectedProblemsJSON, _ := json.Marshal(dailyPlan)
	focusTopicsJSON, _ := json.Marshal(profile.WeakTopics)

	weeklyGoal := models.WeeklyGoal{
		UserID:              user.ID,
		WeekStartDate:       weekStart,
		GoalType:            "GENERATED",
		DifficultyBreakdown: datatypes.JSON(breakdownJSON),
		SelectedProblems:    datatypes.JSON(selectedProblemsJSON),
		FocusTopics:         datatypes.JSON(focusTopicsJSON),
		Status:              "PENDING",
		CreatedAt:           time.Now(),
	}

	return s.GoalRepo.CreateBatch(ctx, []models.WeeklyGoal{weeklyGoal})
}

func (s *GoalService) buildUserProfile(user *models.User) UserProfile {
	// Logic: Analyze TopicStats to find low accuracy topics
	weakTopics := []string{}
	// Mock: if easy solved < 10, maybe "Array" is a good start
	if user.EasySolved < 10 {
		weakTopics = append(weakTopics, "Array")
	}

	return UserProfile{
		TotalSolved:  user.EasySolved + user.MediumSolved + user.HardSolved,
		EasySolved:   user.EasySolved,
		MediumSolved: user.MediumSolved,
		HardSolved:   user.HardSolved,
		WeakTopics:   weakTopics,
	}
}

func (s *GoalService) getDifficultyRatio(user UserProfile) DifficultyRatio {
	if user.TotalSolved < 150 {
		return DifficultyRatio{Easy: 0.5, Medium: 0.4, Hard: 0.1}
	}
	if user.TotalSolved < 400 {
		return DifficultyRatio{Easy: 0.3, Medium: 0.5, Hard: 0.2}
	}
	return DifficultyRatio{Easy: 0.2, Medium: 0.5, Hard: 0.3}
}

func (s *GoalService) selectProblems(candidates []leetcode.APIQuestion, count int) []leetcode.APIQuestion {
	// Deduplicate by ID
	unique := make(map[string]leetcode.APIQuestion)
	for _, p := range candidates {
		unique[p.QuestionFrontendId] = p
	}

	uniqueList := make([]leetcode.APIQuestion, 0, len(unique))
	for _, p := range unique {
		uniqueList = append(uniqueList, p)
	}

	// Simple Scoring/Sorting
	sort.Slice(uniqueList, func(i, j int) bool {
		// Prefer higher acceptance rate for confidence, or lower for challenge?
		// Let's mix: Prefer high frequency if available, else random stir
		return uniqueList[i].AcRate > uniqueList[j].AcRate // Higher AC rate first
	})

	result := []leetcode.APIQuestion{}
	for i := 0; i < count && i < len(uniqueList); i++ {
		result = append(result, uniqueList[i])
	}
	return result
}

func (s *GoalService) distributeAcrossWeek(problems []leetcode.APIQuestion) map[string][]string {
	schedule := make(map[string][]string)

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	// Buckets
	easy := []leetcode.APIQuestion{}
	medium := []leetcode.APIQuestion{}
	hard := []leetcode.APIQuestion{}

	for _, p := range problems {
		switch strings.ToLower(p.Difficulty) {
		case "easy":
			easy = append(easy, p)
		case "medium":
			medium = append(medium, p)
		case "hard":
			hard = append(hard, p)
		}
	}

	// Monday: Easy
	if len(easy) > 0 {
		schedule["Monday"] = append(schedule["Monday"], easy[0].Title)
		easy = easy[1:]
	}

	// Tue-Thu: Medium Focus
	midDays := []string{"Tuesday", "Wednesday", "Thursday"}
	for _, day := range midDays {
		if len(medium) > 0 {
			schedule[day] = append(schedule[day], medium[0].Title)
			medium = medium[1:]
		}
	}

	// Fill rest
	remaining := append(easy, append(medium, hard...)...)
	rand.Shuffle(len(remaining), func(i, j int) { remaining[i], remaining[j] = remaining[j], remaining[i] })

	dayIdx := 0
	for _, p := range remaining {
		day := days[dayIdx%7]
		schedule[day] = append(schedule[day], p.Title)
		dayIdx++
	}

	return schedule
}

func getWeekStart() time.Time {
	now := time.Now().UTC()
	// Calculate start of the week (Monday)
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := now.AddDate(0, 0, offset)
	// Truncate to midnight
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
}

func (s *GoalService) GetUserGoals(ctx context.Context, userID uuid.UUID) ([]models.WeeklyGoal, error) {
	weekStart := getWeekStart()
	return s.GoalRepo.GetWeeklyGoals(ctx, userID, weekStart)
}
