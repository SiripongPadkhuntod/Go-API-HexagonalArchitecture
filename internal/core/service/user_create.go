package service

import (
	"context"
	"fmt"
	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/usecase"
	"strings"
)

func (s *appService) Create(ctx context.Context, input usecase.CreateUserInput) (domain.User, error) {
	if err := validateUser(input.Name, input.Email); err != nil {
		return domain.User{}, err
	}

	user := domain.NewUser(s.ids.NewID(), input.Name, input.Email, s.clock.Now())
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	if err := s.publisher.PublishUserCreated(ctx, createdUser); err != nil { // บรรทัดนี้คือการส่ง event ไปยัง outbound adapter เพื่อนำไปประมวลผลต่อ
		s.logger.Error("failed to publish user created event", "user_id", createdUser.ID, "error", err)
	}

	return createdUser, nil
}

func validateUser(name, email string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if !emailPattern.MatchString(strings.TrimSpace(email)) {
		return fmt.Errorf("%w: email is invalid", domain.ErrInvalidInput)
	}

	return nil
}
