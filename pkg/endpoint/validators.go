package endpoint

type RecentlyPlayedQurey struct {
	Email  string `validate:"required,email"`
	Limit  int    `validate:"number"`
	Before string `validate:"omitempty,datetime=2006-01-02"`
	After  string `validate:"omitempty,datetime=2006-01-02"`
}
