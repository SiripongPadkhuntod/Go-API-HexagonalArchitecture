package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"hexagonalarchitecture/internal/core/domain"
	"hexagonalarchitecture/internal/core/port"
)

func TestCreateUsesRepositoryAndOutboundAdapterInterfaces(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockOutboundAPIClient{}
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
	if !outbound.doCalled {
		t.Fatal("expected outbound API Do to be called")
	}
	if outbound.request.Method != "POST" {
		t.Fatalf("expected outbound method POST, got %s", outbound.request.Method)
	}
	if outbound.request.Path != "/users/events" {
		t.Fatalf("expected outbound path /users/events, got %s", outbound.request.Path)
	}
	if user.Email != "jane@example.com" {
		t.Fatalf("expected normalized email, got %s", user.Email)
	}

	var event userCreatedEvent
	if err := json.Unmarshal(outbound.request.Body, &event); err != nil {
		t.Fatalf("unmarshal outbound event: %v", err)
	}
	if event.ID != user.ID || event.Name != user.Name || event.Email != user.Email {
		t.Fatalf("outbound event does not match created user: %+v", event)
	}
}

func TestCreateInvalidInputDoesNotCallDependencies(t *testing.T) {
	repo := &mockUserRepository{}
	outbound := &mockOutboundAPIClient{}
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
	if outbound.doCalled {
		t.Fatal("outbound API should not be called for invalid input")
	}
}

func TestCreateReturnsRepositoryErrorWithoutCallingOutbound(t *testing.T) {
	repoErr := errors.New("repository failed")
	repo := &mockUserRepository{createErr: repoErr}
	outbound := &mockOutboundAPIClient{}
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
	if outbound.doCalled {
		t.Fatal("outbound API should not be called when repository fails")
	}
}

func TestUpdateUsesRepositoryInterfaceOnly(t *testing.T) {
	existingUser := domain.NewUser("usr_123", "Jane Doe", "jane@example.com")
	repo := &mockUserRepository{findByIDUser: existingUser}
	outbound := &mockOutboundAPIClient{}
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
	if outbound.doCalled {
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

type mockOutboundAPIClient struct {
	doCalled bool
	request  port.OutboundAPIRequest
	err      error
}

func (c *mockOutboundAPIClient) Do(ctx context.Context, request port.OutboundAPIRequest) (port.OutboundAPIResponse, error) {
	c.doCalled = true
	c.request = request
	if c.err != nil {
		return port.OutboundAPIResponse{}, c.err
	}

	return port.OutboundAPIResponse{StatusCode: 204}, nil
}
