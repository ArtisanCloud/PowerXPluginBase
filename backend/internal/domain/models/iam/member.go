package iam

type Member struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Position   string `json:"position"`
	Department string `json:"department"`
	Status     string `json:"status"`
	JoinedAt   string `json:"joined_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
