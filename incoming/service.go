package incoming

import "time"

type IncomingService struct {
    Repo *IncomingRepository
}

func NewService(repo *IncomingRepository) *IncomingService {
    return &IncomingService{
        Repo: repo,
    }
}

func (is *IncomingService) GeneratePath() (string, error) {
    // Generate path with expiry time  time.Minutes * 30
    is.Repo.Set("foobar", "value", time.Duration(time.Minutes * 30))
}

