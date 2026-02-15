# LeetCode Tracker Backend

Backend service for the LeetCode Tracker application, built with Go (Echo), Gorm, and PostgreSQL.

## Features

### User Comparison with Caching
Compare two LeetCode users using AI analysis. The system caches results to optimize performance and reduce API costs.

- **Endpoint**: `POST /api/v1/compare`
- **Caching Logic**:
  - Comparison results are stored in the `user_comparisons` table.
  - **Cache Duration**: 10 hours.
  - **Mechanism**:
    1. Checks for a valid existing comparison between the two users (order-independent) created within the last 10 hours.
    2. If found, returns the cached result.
    3. If not found or expired, fetches a new analysis from the AI provider (Gemini).
    4. Saves the new result to the database for future requests.

## Tech Stack
- **Language**: Go
- **Framework**: Echo
- **Database**: PostgreSQL
- **ORM**: Gorm
- **AI Integration**: Google GenAI (Gemini)

## Setup

1. **Clone the repository**
2. **Configure Environment Variables**
   Create a `.env` file with:
   ```env
   PORT=8080
   DATABASE_URL=postgres://user:password@localhost:5432/dbname
   GEMINI_API_KEY=your_api_key
   LEETCODE_API_ENDPOINT=https://leetcode.com/graphql
   ```
3. **Run the Server**
   ```bash
   go run cmd/server/main.go
   ```

## Development
- **Database Migration**: Automatically handled by Gorm on startup.
- **Linting**: Standard Go tools.
