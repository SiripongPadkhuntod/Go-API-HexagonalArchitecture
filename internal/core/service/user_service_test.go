package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/core/usecase"
)

func TestCreateUsesRepositoryAndOutboundAdapterInterfaces(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockUserEventPublisher{}
	logger := &mockLogger{}
	ids := &mockIDGenerator{id: "usr_test"}
	clock := &mockClock{now: time.Date(2026, 6, 15, 1, 2, 3, 0, time.UTC)}
	users := newTestUserService(repo, outbound, logger, ids, clock)

	user, err := users.Create(context.Background(), usecase.CreateUserInput{
		Name:  "Jane Doe",
		Email: "Jane@Example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}
	if !outbound.publishCalled {
		t.Fatal("expected user event publisher to be called")
	}
	if user.Email != "jane@example.com" {
		t.Fatalf("expected normalized email, got %s", user.Email)
	}
	if user.ID != "usr_test" {
		t.Fatalf("expected injected id, got %s", user.ID)
	}
	if !user.CreatedAt.Equal(clock.now) || !user.UpdatedAt.Equal(clock.now) {
		t.Fatalf("expected injected time, got created=%s updated=%s", user.CreatedAt, user.UpdatedAt)
	}
	if outbound.publishedUser.ID != user.ID || outbound.publishedUser.Name != user.Name || outbound.publishedUser.Email != user.Email {
		t.Fatalf("published user does not match created user: %+v", outbound.publishedUser)
	}
}

func TestCreateInvalidInputDoesNotCallDependencies(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockUserEventPublisher{}
	logger := &mockLogger{}
	ids := &mockIDGenerator{id: "usr_test"}
	clock := &mockClock{now: time.Date(2026, 6, 15, 1, 2, 3, 0, time.UTC)}
	users := newTestUserService(repo, outbound, logger, ids, clock)

	_, err := users.Create(context.Background(), usecase.CreateUserInput{
		Name:  "",
		Email: "invalid-email",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if repo.createCalled {
		t.Fatal("repository should not be called for invalid input")
	}
	if outbound.publishCalled {
		t.Fatal("outbound API should not be called for invalid input")
	}
}

func TestCreateReturnsRepositoryErrorWithoutCallingOutbound(t *testing.T) {
	repoErr := errors.New("repository failed")
	repo := &mockUserRepository{createErr: repoErr}
	outbound := &mockUserEventPublisher{}
	logger := &mockLogger{}
	ids := &mockIDGenerator{id: "usr_test"}
	clock := &mockClock{now: time.Date(2026, 6, 15, 1, 2, 3, 0, time.UTC)}
	users := newTestUserService(repo, outbound, logger, ids, clock)

	_, err := users.Create(context.Background(), usecase.CreateUserInput{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}
	if outbound.publishCalled {
		t.Fatal("outbound API should not be called when repository fails")
	}
}

func TestCreateReturnsUserWhenBestEffortPublisherFails(t *testing.T) {
	publisherErr := errors.New("publisher failed")
	repo := &mockUserRepository{}
	outbound := &mockUserEventPublisher{err: publisherErr}
	logger := &mockLogger{}
	ids := &mockIDGenerator{id: "usr_test"}
	clock := &mockClock{now: time.Date(2026, 6, 15, 1, 2, 3, 0, time.UTC)}
	users := newTestUserService(repo, outbound, logger, ids, clock)

	user, err := users.Create(context.Background(), usecase.CreateUserInput{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})
	if err != nil {
		t.Fatalf("create should ignore best-effort publisher error, got %v", err)
	}
	if strings.TrimSpace(user.ID) == "" {
		t.Fatal("expected created user id")
	}
	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}
	if !outbound.publishCalled {
		t.Fatal("expected user event publisher to be called")
	}
	if !logger.errorCalled {
		t.Fatal("expected publisher error to be logged")
	}
}

func TestUpdateUsesRepositoryInterfaceOnly(t *testing.T) {
	createdAt := time.Date(2026, 6, 15, 1, 2, 3, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	existingUser := domain.NewUser("usr_123", "Jane Doe", "jane@example.com", createdAt)
	repo := &mockUserRepository{findByIDUser: existingUser}
	outbound := &mockUserEventPublisher{}
	logger := &mockLogger{}
	ids := &mockIDGenerator{id: "unused"}
	clock := &mockClock{now: updatedAt}
	users := newTestUserService(repo, outbound, logger, ids, clock)

	updatedUser, err := users.Update(context.Background(), "usr_123", usecase.UpdateUserInput{
		Name:  "Jane Smith",
		Email: "Jane.Smith@Example.com",
	})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}

	if !repo.findByIDCalled {
		t.Fatal("expected repository FindByID to be called")
	}
	if !repo.updateCalled {
		t.Fatal("expected repository Update to be called")
	}
	if outbound.publishCalled {
		t.Fatal("outbound API should not be called by update logic")
	}
	if updatedUser.Name != "Jane Smith" {
		t.Fatalf("expected updated name, got %s", updatedUser.Name)
	}
	if updatedUser.Email != "jane.smith@example.com" {
		t.Fatalf("expected normalized email, got %s", updatedUser.Email)
	}
	if !updatedUser.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected created time to remain unchanged, got %s", updatedUser.CreatedAt)
	}
	if !updatedUser.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("expected injected updated time, got %s", updatedUser.UpdatedAt)
	}
}

type mockUserRepository struct {
	createCalled   bool
	findAllCalled  bool
	findByIDCalled bool
	updateCalled   bool
	deleteCalled   bool

	createErr    error
	findByIDUser domain.User
	users        []domain.User
}

func newTestUserService(repo port.AppRepository, publisher port.UserEventPublisher, logger port.Logger, ids port.IDGenerator, clock port.Clock) port.AppService {
	return NewAppService(AppServiceDeps{
		Repo:      repo,
		Publisher: publisher,
		Logger:    logger,
		IDs:       ids,
		Clock:     clock,
	})
}

func (r *mockUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	r.createCalled = true
	if r.createErr != nil {
		return domain.User{}, r.createErr
	}

	return user, nil
}

func (r *mockUserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	r.findAllCalled = true
	return r.users, nil
}

func (r *mockUserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	r.findByIDCalled = true
	if strings.TrimSpace(r.findByIDUser.ID) == "" {
		return domain.User{}, domain.ErrUserNotFound
	}

	return r.findByIDUser, nil
}

func (r *mockUserRepository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	r.updateCalled = true
	return user, nil
}

func (r *mockUserRepository) Delete(ctx context.Context, id string) error {
	r.deleteCalled = true
	return nil
}

type mockUserEventPublisher struct {
	publishCalled bool
	publishedUser domain.User
	err           error
}

func (p *mockUserEventPublisher) PublishUserCreated(ctx context.Context, user domain.User) error {
	p.publishCalled = true
	p.publishedUser = user
	return p.err
}

type mockLogger struct {
	infoCalled  bool
	errorCalled bool
	fatalCalled bool
}

func (l *mockLogger) Info(msg string, args ...any) {
	l.infoCalled = true
}

func (l *mockLogger) Error(msg string, args ...any) {
	l.errorCalled = true
}

func (l *mockLogger) Fatal(msg string, args ...any) {
	l.fatalCalled = true
}

type mockIDGenerator struct {
	id string
}

func (g *mockIDGenerator) NewID() string {
	return g.id
}

type mockClock struct {
	now time.Time
}

func (c *mockClock) Now() time.Time {
	return c.now
}
