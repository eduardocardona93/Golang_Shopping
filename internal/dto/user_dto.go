package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/eduardocardona93/golang_shopping/internal/domain"
)

// CreateUserRequest is the payload to create a new User.
type CreateUserRequest struct {
	Nombre         string `json:"nombre" validate:"required,min=1,max=100"`
	Apellido       string `json:"apellido" validate:"required,min=1,max=100"`
	Telefono       string `json:"telefono" validate:"omitempty,max=20"`
	Direccion      string `json:"direccion" validate:"omitempty,max=255"`
	Identificacion string `json:"identificacion" validate:"required,min=1,max=50"`
}

// UpdateUserRequest is the payload to update an existing User.
type UpdateUserRequest struct {
	Nombre         string `json:"nombre" validate:"required,min=1,max=100"`
	Apellido       string `json:"apellido" validate:"required,min=1,max=100"`
	Telefono       string `json:"telefono" validate:"omitempty,max=20"`
	Direccion      string `json:"direccion" validate:"omitempty,max=255"`
	Identificacion string `json:"identificacion" validate:"required,min=1,max=50"`
}

// UserResponse is the API representation of a User.
type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	Nombre         string    `json:"nombre"`
	Apellido       string    `json:"apellido"`
	Telefono       string    `json:"telefono"`
	Direccion      string    `json:"direccion"`
	Identificacion string    `json:"identificacion"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r CreateUserRequest) ToDomain() *domain.User {
	return &domain.User{
		Nombre:         r.Nombre,
		Apellido:       r.Apellido,
		Telefono:       r.Telefono,
		Direccion:      r.Direccion,
		Identificacion: r.Identificacion,
	}
}

func (r UpdateUserRequest) ApplyTo(u *domain.User) {
	u.Nombre = r.Nombre
	u.Apellido = r.Apellido
	u.Telefono = r.Telefono
	u.Direccion = r.Direccion
	u.Identificacion = r.Identificacion
}

func NewUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:             u.ID,
		Nombre:         u.Nombre,
		Apellido:       u.Apellido,
		Telefono:       u.Telefono,
		Direccion:      u.Direccion,
		Identificacion: u.Identificacion,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func NewUserResponseList(users []domain.User) []UserResponse {
	out := make([]UserResponse, 0, len(users))
	for i := range users {
		out = append(out, NewUserResponse(&users[i]))
	}
	return out
}
