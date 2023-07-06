package server

import (
	repository2 "cardsService/internal/repository"
	"cardsService/pkg/helpers"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Cards handlers

func (s *ServiceImpl) CardsHandler(c echo.Context) error {
	// check authorization
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.Repository.GetCards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) CardsByTokenHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.Repository.GetCardsByUserId(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) GetCardsByIDHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}
	id := c.Param("id")

	cards, err := s.Repository.GetCardsByUserId(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) CardsByBookDate(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.Repository.GetBookedCardsByUser(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) CardsByUsersBooked(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	bookInfo, err := s.Repository.GetUsersBookedCards(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &bookInfo)
}

func (s *ServiceImpl) UpdateCardsHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateCardsRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	tags := make([]string, 0)
	bookdates := make([]string, 0)
	bookdatesUserId := make([]string, 0)
	for _, t := range req.Tags {
		tags = append(tags, t.Name)
	}
	for _, d := range req.BookDates {
		bookdates = append(bookdates, d.Date)
		bookdatesUserId = append(bookdatesUserId, d.UserId)
	}

	err = s.Repository.AddCard(token.Role, &repository2.UpdateCard{
		CardID:          req.CardID,
		UserID:          req.UserID,
		Title:           req.Title,
		Description:     req.Description,
		IsActive:        req.IsActive,
		Cost:            req.Cost,
		Tags:            tags,
		BookDates:       bookdates,
		BookDatesUserId: bookdatesUserId,
		IsAgreement:     req.IsAgreement,
		Prepayment:      req.Prepayment,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

// Users handlers

func (s *ServiceImpl) AddUserHandler(c echo.Context) error {
	req := &AddUserRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	if req.MasterPassword != s.masterPassword {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	//save to db (password md5)
	err = s.Repository.AddUser(&repository2.AddUserParams{
		UserName: req.UserName,
		Password: helpers.GetMD5Hash(req.Password),
		RoleName: req.RoleName,
		Email:    req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) LoginHandler(c echo.Context) error {
	req := &LoginRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	userId, role, err := s.Repository.CheckUser(&repository2.CheckUserParams{
		Email:    req.Email,
		Password: helpers.GetMD5Hash(req.Password),
	})
	if err != nil {
		if errors.Is(err, repository2.ErrNotFound) {
			return c.JSON(http.StatusNotFound, &ErrorResponse{ErrorMessage: "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	token, err := s.TokenService.GetToken(userId, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &LoginResponse{Token: token})
}

func (s *ServiceImpl) GetUserByIDHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}
	id := c.Param("id")

	user, err := s.Repository.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &user)
}

func (s *ServiceImpl) GetUserByTokenHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	user, err := s.Repository.GetUserById(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &user)
}

func (s *ServiceImpl) UsersHandler(c echo.Context) error {
	// check authorization
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	users, err := s.Repository.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &users)
}

// User info handlers

func (s *ServiceImpl) GetSpecializationByIDHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &IDRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	spec, err := s.Repository.GetSpecializationById(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &spec)
}

func (s *ServiceImpl) GetAddInfByIDHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &IDRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	addInf, err := s.Repository.GetAddInfById(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &addInf)
}

func (s *ServiceImpl) UpdateSpecHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateSpecRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UpdateSpec(token.Role, &repository2.UpdateSpecialization{
		UserID:          req.UserID,
		SpecID:          req.SpecID,
		SpecName:        req.SpecName,
		SpecDescription: req.SpecDescription,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) UpdateAddInfHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateAddInfRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UpdateAddInf(token.Role, &repository2.UpdateAddInf{
		UserID:      req.UserID,
		AddInfID:    req.AddInfID,
		Description: req.Description,
		Country:     req.Country,
		City:        req.City,
		TypeOfWork:  req.TypeOfWork,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) UpdateCreatorStatusHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateCreatorStatusRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UpdateCreatorStatus(&repository2.UpdateCreatorStatusParams{
		UserID:   req.UserID,
		UserName: req.UserName,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) UpdateBookDateUserHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateBookDateUserRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UpdateBookDatesUser(&repository2.UpdateBookDateUserParams{
		UserID: req.UserID,
		BookID: req.BookID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

// Image handlers

func (s *ServiceImpl) UploadImageHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	file, err := c.FormFile("pic")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	defer src.Close()

	path := "images/" + file.Filename

	dst, err := os.Create(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UploadImage(&repository2.UploadImageParams{
		ID:   token.UserId,
		Path: path,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) GetImage(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	token, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	path, err := s.Repository.GetImage(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	if path == nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: "User doesn't have a profile image"})
	}

	file, err := os.Open(path.Path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	data := make([]byte, fileInfo.Size())
	if _, err := file.Read(data); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	contentType := http.DetectContentType(data)
	return c.JSON(http.StatusOK, &ImageResponse{
		Filename: filepath.Base(path.Path),
		Content:  data,
		Type:     contentType,
	})
}

// Messenger handlers

func (s *ServiceImpl) SaveFMCToken(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateFCMToken{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.Repository.UpdateFCMToken(&repository2.UpdateFCMTokenParams{
		UserID: req.UserID,
		Token:  req.Token,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) SendMessageHandler(c echo.Context) error {
	authToken := s.GetAuthToken(c)
	_, err := s.TokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &MessageRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	fmt.Println(req)

	err = s.Repository.SendPush(&repository2.MessageStruct{
		ID:             req.ID,
		Message:        req.Message,
		SenderUsername: req.SenderUsername,
		To:             req.To,
	})
	if err != nil {
		fmt.Println("err", err)
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}
