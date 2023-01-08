package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func (s *ServiceImpl) AddUser(params *AddUserParams) error {
	err := s.db.Exec("SELECT role_create(?, ?, ?, ?)", params.UserName, params.Password, params.RoleName, params.Email).Error
	if err != nil {
		return fmt.Errorf("error during AddUser, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) CheckUser(params *CheckUserParams) (string, string, error) {
	account := &Account{}
	user := &User{}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Exec("SET role ?", gorm.Expr("postgres")).Error
		if err != nil {
			return err
		}
		err = s.db.Raw("SELECT * FROM account WHERE email = ? AND password = ?",
			params.Email, params.Password).First(account).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrNotFound
			}
			return err
		}

		err = s.db.Raw("SELECT * FROM users WHERE user_id = ?", account.UserID).First(user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrNotFound
			}
			return err
		}
		return nil
	})
	return account.UserID, user.Username, err
}

func (s *ServiceImpl) GetUserById(id string) (*User, error) {
	user := &User{}
	err := s.db.Raw("SELECT * FROM users WHERE user_id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *ServiceImpl) GetCards() ([]*Cards, error) {
	cards := make([]*Cards, 0)
	tags := make([]*Tags, 0)
	bookdate := make([]*BookDates, 0)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Raw("SELECT * FROM cards as cr JOIN payment as pmt on cr.fk_payment_id = pmt.payment_id WHERE is_active = true ORDER BY create_time DESC LIMIT 100").Find(&cards).Error
		if err != nil {
			return err
		}

		for _, c := range cards {
			err := tx.Raw("SELECT * from tags where fk_card_id = ?", c.ID).Find(&tags).Error
			if err != nil {
				return err
			}
			err = tx.Raw("SELECT * from book_date where fk_card_id = ?", c.ID).Find(&bookdate).Error
			if err != nil {
				return err
			}
			c.TagsList = tags
			c.BookDates = bookdate
		}
		return nil
	})
	return cards, err
}

func (s *ServiceImpl) GetCardsByUserId(id string) ([]*Cards, error) {
	cards := make([]*Cards, 0)
	tags := make([]*Tags, 0)
	bookdate := make([]*BookDates, 0)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Raw("SELECT * FROM cards as cr JOIN payment as pmt on cr.fk_payment_id = pmt.payment_id WHERE fk_user_id = ? ORDER BY create_time DESC", id).Find(&cards).Error
		if err != nil {
			return err
		}

		for _, c := range cards {
			err := tx.Raw("SELECT * from tags where fk_card_id = ?", c.ID).Find(&tags).Error
			if err != nil {
				return err
			}
			err = tx.Raw("SELECT * from book_date where fk_card_id = ?", c.ID).Find(&bookdate).Error
			if err != nil {
				return err
			}
			c.TagsList = tags
			c.BookDates = bookdate
		}
		return nil
	})
	return cards, err
}

func (s *ServiceImpl) GetBookedCards(id string) ([]*Cards, error) {
	cards := make([]*Cards, 0)
	tags := make([]*Tags, 0)
	bookdate := make([]*BookDates, 0)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Raw("SELECT * FROM book_date WHERE user_book_id = ?", id).Find(&bookdate).Error
		if err != nil {
			return err
		}

		for _, b := range bookdate {
			err := tx.Raw("SELECT * from cards where card_id = ?", b.CardID).Find(&cards).Error
			if err != nil {
				return err
			}
		}
		for _, c := range cards {
			err := tx.Raw("SELECT * from tags where fk_card_id = ?", c.ID).Find(&tags).Error
			if err != nil {
				return err
			}
			err = tx.Raw("SELECT * from book_date where fk_card_id = ?", c.ID).Find(&bookdate).Error
			if err != nil {
				return err
			}
			c.TagsList = tags
			c.BookDates = bookdate
		}
		return nil
	})
	return cards, err
}

func (s *ServiceImpl) GetUsers() ([]*User, error) {
	users := make([]*User, 0)
	err := s.db.Raw("SELECT * FROM users WHERE is_creator = true").Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return users, nil
}

func (s *ServiceImpl) GetSpecializationById(id string) (*Specialization, error) {
	spec := &Specialization{}
	err := s.db.Raw("SELECT * FROM specialization WHERE spec_id = ?", id).First(spec).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return spec, nil
}

func (s *ServiceImpl) GetAddInfById(id string) (*AdditionalInfo, error) {
	addInf := &AdditionalInfo{}
	err := s.db.Raw("SELECT * FROM additional_info WHERE add_inf_id = ?", id).First(addInf).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return addInf, nil
}

func (s *ServiceImpl) UpdateSpec(role string, params *UpdateSpecialization) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Exec("SET role ?", gorm.Expr(role)).Error
		if err != nil {
			return err
		}
		err = s.db.Exec("CALL update_specialization(?, ?, ?, ?)", params.UserID, params.SpecID,
			params.SpecName, params.SpecDescription).Error
		if err != nil {
			return fmt.Errorf("error during UpdateSpec, err: %w", err)
		}
		return nil
	})
	return err
}

func (s *ServiceImpl) UpdateAddInf(role string, params *UpdateAddInf) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Exec("SET role ?", gorm.Expr(role)).Error
		if err != nil {
			return err
		}
		err = s.db.Exec("CALL update_additional_inf(?, ?, ?, ?, ?, ?)", params.UserID, params.AddInfID,
			params.Description, params.Country, params.City, params.TypeOfWork).Error
		if err != nil {
			return fmt.Errorf("error during UpdateSpec, err: %w", err)
		}
		return nil
	})
	return err
}

func (s *ServiceImpl) UpdateCreatorStatus(params *UpdateCreatorStatusParams) error {
	err := s.db.Exec("SELECT update_role_to_creator(?, ?)", params.UserID, params.UserName).Error
	if err != nil {
		return fmt.Errorf("error during UpdateSpec, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) UpdateBookDatesUser(params *UpdateBookDateUserParams) error {
	err := s.db.Exec("SELECT update_book_date_user(?, ?)", params.UserID, params.BookID).Error
	if err != nil {
		return fmt.Errorf("error during UpdateSpec, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) AddCard(role string, params *UpdateCard) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Exec("SET role ?", gorm.Expr(role)).Error
		if err != nil {
			return err
		}
		err = s.db.Exec("CALL update_cards(?, ?, ?, ?, ?, ?, ?, ?)", params.CardID, params.UserID,
			params.Title, params.Description, params.IsActive, params.Cost, params.IsAgreement, params.Prepayment).Error
		if err != nil {
			return fmt.Errorf("error during UpdateSpec, err: %w", err)
		}

		if len(params.Tags) > 0 {
			for i, tag := range params.Tags {
				err = s.db.Exec("CALL update_tags(?, ?, ?)", params.CardID, tag, i).Error
				if err != nil {
					return fmt.Errorf("error during UpdateSpec, err: %w", err)
				}
			}
		} else {
			err = s.db.Exec("DELETE FROM tags WHERE fk_card_id = ?", params.CardID).Error
			if err != nil {
				return fmt.Errorf("error during UpdateSpec, err: %w", err)
			}
		}

		if len(params.BookDates) > 0 {
			for i, date := range params.BookDates {
				err = s.db.Exec("CALL update_book_date(?, ?, ?)", params.CardID, date, i).Error
				if err != nil {
					return fmt.Errorf("error during UpdateSpec, err: %w", err)
				}
			}
		} else {
			err = s.db.Exec("DELETE FROM book_date WHERE fk_card_id = ?", params.CardID).Error
			if err != nil {
				return fmt.Errorf("error during UpdateSpec, err: %w", err)
			}
		}

		return nil
	})
	return err
}
