package models

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Company string `json:"company"`
	JobTitle string `json:"job_title"`
}

type RecordValue struct {
	Count int 
	Data User
}