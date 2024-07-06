package http

type LoginRequest struct {
	Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
}

type UpdateRequest struct {
	Email          string  `json:"email" validate:"omitempty,lte=60,email"`
	Password       string  `json:"password" validate:"omitempty,gte=6"`
	Name           string  `json:"name" validate:"omitempty,gte=2"`
	Surname        string  `json:"surname" validate:"omitempty,gte=2,lte=60"`
	Patronymic     *string `json:"patronymic" validate:"omitempty,gte=2,lte=60"`
	Address        string  `json:"address" validate:"omitempty,lte=100"`
	PassportNumber int     `json:"passport_number" validate:"omitempty,len=4"`
	PassportSeries int     `json:"passport_series" validate:"omitempty,len=4"`
}
