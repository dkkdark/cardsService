package repository

type Ping struct {
	Result int `gorm:"column:result"`
}

type Account struct {
	ID       string `json:"id" gorm:"column:account_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   string `json:"user_id" gorm:"column:fk_user_id"`
}

type Cards struct {
	ID          string       `json:"card_id" gorm:"column:card_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreateTime  int64        `json:"create_time"`
	IsActive    bool         `json:"is_active"`
	UserID      string       `json:"user_id" gorm:"column:fk_user_id"`
	PaymentID   string       `json:"payment_id" gorm:"column:fk_payment_id"`
	Cost        float32      `json:"cost"`
	IsAgreement bool         `json:"is_agreement"`
	Prepayment  bool         `json:"is_prepayment" gorm:"column:is_prepayment"`
	TagsList    []*Tags      `json:"tags_list"`
	BookDates   []*BookDates `json:"book_date_list"`
}

type Tags struct {
	ID     string `json:"id" gorm:"column:tag_id"`
	Name   string `json:"name" gorm:"column:tags_name"`
	CardID string `json:"card_id" gorm:"column:fk_card_id"`
}

type BookDates struct {
	ID     string `json:"id" gorm:"column:id"`
	Date   string `json:"date" gorm:"column:possible_date"`
	CardID string `json:"card_id" gorm:"column:fk_card_id"`
	UserId string `json:"user_id" gorm:"column:user_book_id"`
}

type User struct {
	ID             string `json:"id" gorm:"column:user_id"`
	Username       string `json:"username"`
	Image          string `json:"image"`
	IsCreator      bool   `json:"is_creator"`
	Specialization string `json:"specialization" gorm:"column:fk_specialization_id"`
	AddInf         string `json:"add_inf" gorm:"column:fk_add_inf_id"`
	RoleName       string `json:"role_name" gorm:"column:rolename"`
}

type Specialization struct {
	ID          string `json:"id" gorm:"column:spec_id"`
	Name        string `json:"name" gorm:"column:spec_name"`
	Description string `json:"description" gorm:"column:spec_description"`
}

type AdditionalInfo struct {
	AddInfID    string `json:"id" gorm:"column:add_inf_id"`
	Description string `json:"description" gorm:"column:add_inf_descr"`
	Country     string `json:"country"`
	City        string `json:"city"`
	TypeOfWork  string `json:"type_of_work"`
}
