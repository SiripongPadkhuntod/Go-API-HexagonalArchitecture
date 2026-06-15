package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"hexagonalarchitecture/internal/core/domain"
)

func TestCreateUsesRepositoryAndOutboundAdapterInterfaces(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockUserEventPublisher{}
	users := NewUserService(repo, outbound)

	user, err := users.Create(context.Background(), CreateUserInput{
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
	if outbound.publishedUser.ID != user.ID || outbound.publishedUser.Name != user.Name || outbound.publishedUser.Email != user.Email {
		t.Fatalf("published user does not match created user: %+v", outbound.publishedUser)
	}
}

func TestCreateInvalidInputDoesNotCallDependencies(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockUserEventPublisher{}
	users := NewUserService(repo, outbound)

	_, err := users.Create(context.Background(), CreateUserInput{
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
	users := NewUserService(repo, outbound)

	_, err := users.Create(context.Background(), CreateUserInput{
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
	users := NewUserService(repo, outbound)

	user, err := users.Create(context.Background(), CreateUserInput{
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
}

func TestUpdateUsesRepositoryInterfaceOnly(t *testing.T) {
	existingUser := domain.NewUser("usr_123", "Jane Doe", "jane@example.com")
	repo := &mockUserRepository{findByIDUser: existingUser}
	outbound := &mockUserEventPublisher{}
	users := NewUserService(repo, outbound)

	updatedUser, err := users.Update(context.Background(), "usr_123", UpdateUserInput{
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
