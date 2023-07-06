package repository

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/jinzhu/gorm"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"path/filepath"
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

func (s *ServiceImpl) GetUserById(id string) (*Freelancer, error) {
	user := &User{}
	freelancer := &Freelancer{}

	err := s.db.Raw("SELECT * FROM users WHERE user_id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	image := &Image{}
	if user.Image != "" {
		image, _ = s.GetImageData(user.Image)
	}
	spec, _ := s.GetSpecializationById(user.Specialization)
	addInf, _ := s.GetAddInfById(user.AddInf)

	freelancer.ID = user.ID
	freelancer.Username = user.Username
	freelancer.RoleName = user.RoleName
	freelancer.IsCreator = user.IsCreator
	freelancer.Image = image
	freelancer.Specialization = spec
	freelancer.AddInf = addInf

	return freelancer, nil
}

func (s *ServiceImpl) GetImageData(path string) (img *Image, err error) {
	image := &Image{}
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file into a byte array
	data := make([]byte, fileInfo.Size())
	if _, err := file.Read(data); err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(data)
	image.Filename = filepath.Base(image.Filename)
	image.Content = data
	image.Type = contentType

	return image, nil
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

func (s *ServiceImpl) GetBookedCardsByUser(id string) ([]*Cards, error) {
	cards := make([]*Cards, 0)
	tags := make([]*Tags, 0)
	bookdate := make([]*BookDates, 0)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Raw("SELECT DISTINCT ON (fk_card_id) * FROM book_date WHERE user_book_id = ?", id).Find(&bookdate).Error
		if err != nil {
			return err
		}

		for _, b := range bookdate {
			card := &Cards{}
			err := tx.Raw("SELECT * FROM cards as cr JOIN payment as pmt on cr.fk_payment_id = pmt.payment_id where card_id = ?", b.CardID).First(&card).Error
			if err != nil {
				return err
			}
			err = tx.Raw("SELECT * from tags where fk_card_id = ?", card.ID).Find(&tags).Error
			if err != nil {
				return err
			}
			err = tx.Raw("SELECT * from book_date where fk_card_id = ?", card.ID).Find(&bookdate).Error
			if err != nil {
				return err
			}
			card.TagsList = tags
			card.BookDates = bookdate

			cards = append(cards, card)
		}
		return nil
	})
	return cards, err
}

func (s *ServiceImpl) GetUsersBookedCards(id string) ([]*BookedInfo, error) {
	bookInfo := make([]*BookedInfo, 0)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		err := tx.Raw("SELECT card_id, title, possible_date, user_book_id, username FROM cards JOIN book_date on fk_card_id = card_id JOIN users on user_id = user_book_id where fk_user_id = ? and user_book_id is not null", id).Find(&bookInfo).Error
		if err != nil {
			return err
		}
		return nil
	})
	return bookInfo, err
}

func (s *ServiceImpl) GetUsers() ([]*Freelancer, error) {
	users := make([]*User, 0)
	freelancers := make([]*Freelancer, 0)
	err := s.db.Raw("SELECT * FROM users WHERE is_creator = true ORDER BY username").Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	for _, u := range users {
		spec, _ := s.GetSpecializationById(u.Specialization)
		addInf, _ := s.GetAddInfById(u.AddInf)
		image, _ := s.GetImageData(u.Image)

		freelancer := &Freelancer{}
		freelancer.ID = u.ID
		freelancer.Username = u.Username
		freelancer.RoleName = u.RoleName
		freelancer.IsCreator = u.IsCreator
		freelancer.Image = image
		freelancer.Specialization = spec
		freelancer.AddInf = addInf

		freelancers = append(freelancers, freelancer)
	}

	return freelancers, nil
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
			return fmt.Errorf("error during UpdateAddInf, err: %w", err)
		}
		return nil
	})
	return err
}

func (s *ServiceImpl) UpdateCreatorStatus(params *UpdateCreatorStatusParams) error {
	err := s.db.Exec("SELECT update_role_to_creator(?, ?)", params.UserID, params.UserName).Error
	if err != nil {
		return fmt.Errorf("error during UpdateCreatorStatus, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) UpdateBookDatesUser(params *UpdateBookDateUserParams) error {
	err := s.db.Exec("SELECT update_book_date_user(?, ?)", params.UserID, params.BookID).Error
	if err != nil {
		return fmt.Errorf("error during UpdateBookDatesUser, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) UploadImage(params *UploadImageParams) error {
	err := s.db.Exec("UPDATE users SET image = ? WHERE user_id = ?", params.Path, params.ID).Error
	if err != nil {
		return fmt.Errorf("error during UploadImage, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) UpdateFCMToken(params *UpdateFCMTokenParams) error {
	err := s.db.Exec("UPDATE users SET fcm_token = ? WHERE user_id = ?", params.Token, params.UserID).Error
	if err != nil {
		return fmt.Errorf("error during UploadImage, err: %w", err)
	}
	return nil
}

func (s *ServiceImpl) GetFCMToken(userId string) (*TokenFCMStructure, error) {
	token := &TokenFCMStructure{}
	err := s.db.Raw("SELECT fcm_token FROM users WHERE user_id = ?", userId).First(token).Error
	if err != nil {
		return nil, fmt.Errorf("error during GetFCMToken, err: %w", err)
	}
	return token, nil
}

func (s *ServiceImpl) GetImage(id string) (*PathStructure, error) {
	path := &PathStructure{}
	err := s.db.Raw("SELECT image FROM users WHERE user_id = ?", id).First(path).Error
	if err != nil {
		return nil, fmt.Errorf("error during GetImage, err: %w", err)
	}
	return path, nil
}

func (s *ServiceImpl) SendPush(params *MessageStruct) error {
	ctx := context.Background()

	token, err := s.GetFCMToken(params.To)
	if err != nil {
		return err
	}
	fmt.Println("token: ", token)

	// you should add your own json file
	serviceAccountKeyFilePath, err := filepath.Abs("./tasks-app-6ef9f-firebase-adminsdk-wm7zf-6c08b5a721.json")
	if err != nil {
		return err
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: params.SenderUsername,
			Body:  params.Message,
		},
		Token: token.FCMToken,
		Data: map[string]string{
			"id": params.ID,
		},
	}

	response, err := fcmClient.Send(ctx, message)
	if err != nil {
		return err
	}

	fmt.Println("Successfully sent message:", response)
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
			return fmt.Errorf("error during AddCard, err: %w", err)
		}

		if len(params.Tags) > 0 {
			for i, tag := range params.Tags {
				err = s.db.Exec("CALL update_tags(?, ?, ?)", params.CardID, tag, i).Error
				if err != nil {
					return fmt.Errorf("error during AddCard, err: %w", err)
				}
			}
		} else {
			err = s.db.Exec("DELETE FROM tags WHERE fk_card_id = ?", params.CardID).Error
			if err != nil {
				return fmt.Errorf("error during AddCard, err: %w", err)
			}
		}

		if len(params.BookDates) > 0 {
			for i, date := range params.BookDates {
				if params.BookDatesUserId[i] == "" {
					err = s.db.Exec("CALL update_book_date(?, ?, ?, ?)", params.CardID, date, nil, i).Error
					if err != nil {
						return fmt.Errorf("error during AddCard, err: %w", err)
					}
				} else {
					err = s.db.Exec("CALL update_book_date(?, ?, ?, ?)", params.CardID, date, params.BookDatesUserId[i], i).Error
					if err != nil {
						return fmt.Errorf("error during AddCard, err: %w", err)
					}
				}
			}
		} else {
			err = s.db.Exec("DELETE FROM book_date WHERE fk_card_id = ?", params.CardID).Error
			if err != nil {
				return fmt.Errorf("error during AddCard, err: %w", err)
			}
		}

		return nil
	})
	return err
}
