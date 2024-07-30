package http

type LoginRequest struct {
	Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
}

type UpdateUserRequest struct {
	ID         int64   `json:"-"`
	Email      *string `json:"email" validate:"omitempty,lte=60,email"`
	Password   *string `json:"password" validate:"omitempty,gte=6,lte=256"`
	Name       *string `json:"name" validate:"omitempty,gte=2,lte=60"`
	Surname    *string `json:"surname" validate:"omitempty,gte=2,lte=60"`
	Patronymic *string `json:"patronymic" validate:"omitempty,gte=2,lte=60"`
	Address    *string `json:"address" validate:"omitempty,lte=100"`
}
