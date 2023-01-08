package server

type AddUserRequest struct {
	MasterPassword string `json:"master_password"`
	UserName       string `json:"user_name"`
	Password       string `json:"password"`
	RoleName       string `json:"role_name"`
	Email          string `json:"email"`
}

type EmptyResponse struct{}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type IDRequest struct {
	ID string `json:"id"`
}

type UpdateSpecRequest struct {
	UserID          string `json:"id"`
	SpecID          string `json:"spec_id"`
	SpecName        string `json:"name"`
	SpecDescription string `json:"description"`
}

type UpdateAddInfRequest struct {
	UserID      string `json:"id"`
	AddInfID    string `json:"add_inf_id"`
	Description string `json:"description"`
	Country     string `json:"country"`
	City        string `json:"city"`
	TypeOfWork  string `json:"type_of_work"`
}

type UpdateCreatorStatusRequest struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type UpdateBookDateUserRequest struct {
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
}

type TagsList struct {
	Name string
}

type BookDates struct {
	Date string `json:"possible_date"`
}

type UpdateCardsRequest struct {
	CardID      string       `json:"card_id"`
	UserID      string       `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	IsActive    bool         `json:"is_active"`
	Cost        float32      `json:"cost"`
	Tags        []*TagsList  `json:"tags_list"`
	BookDates   []*BookDates `json:"book_date_list"`
	IsAgreement bool         `json:"is_agreement"`
	Prepayment  bool         `json:"is_prepayment"`
}
