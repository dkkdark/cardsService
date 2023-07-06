package repository

type AddUserParams struct {
	UserName string
	Password string
	RoleName string
	Email    string
}

type CheckUserParams struct {
	Email    string
	Password string
}

type MessageStruct struct {
	ID             string
	SenderUsername string
	Message        string
	To             string
}

type UpdateSpecialization struct {
	UserID          string
	SpecID          string
	SpecName        string
	SpecDescription string
}

type UpdateAddInf struct {
	UserID      string
	AddInfID    string
	Description string
	Country     string
	City        string
	TypeOfWork  string
}

type UpdateCreatorStatusParams struct {
	UserID   string
	UserName string
}

type UpdateBookDateUserParams struct {
	UserID string
	BookID string
}

type UploadImageParams struct {
	ID   string
	Path string
}

type UpdateFCMTokenParams struct {
	UserID string
	Token  string
}

type TagsList struct {
	Name string
}

type UpdateCard struct {
	CardID          string
	UserID          string
	Title           string
	Description     string
	IsActive        bool
	Cost            float32
	Tags            []string
	BookDates       []string
	BookDatesUserId []string
	IsAgreement     bool
	Prepayment      bool
}

type BookDatesStruct struct {
	ID     string
	Date   string
	CardID string
	UserId string
}
