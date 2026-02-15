package repository

import (
	"context"
	"math/rand"

	"github.com/devlpr-nitish/leetcode-tracker-backend/pkg/leetcode"
)

type ProblemRepository interface {
	GetProblemsByTopic(ctx context.Context, topic string, limit int) ([]leetcode.APIQuestion, error)
	GetProblemsByDifficulty(ctx context.Context, difficulty string, limit int) ([]leetcode.APIQuestion, error)
}

type problemRepository struct {
	client *leetcode.Client
}

func NewProblemRepository(client *leetcode.Client) ProblemRepository {
	return &problemRepository{client: client}
}

func (r *problemRepository) GetProblemsByTopic(ctx context.Context, topic string, limit int) ([]leetcode.APIQuestion, error) {
	// Random skip to get variety
	skip := rand.Intn(50)
	return r.client.GetProblems(limit, skip, []string{topic}, "")
}

func (r *problemRepository) GetProblemsByDifficulty(ctx context.Context, difficulty string, limit int) ([]leetcode.APIQuestion, error) {
	// Random skip
	skip := rand.Intn(100)
	return r.client.GetProblems(limit, skip, nil, difficulty)
}
