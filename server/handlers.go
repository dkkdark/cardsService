package server

import (
	"cardsService/helpers"
	"cardsService/repository"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
	err = s.repository.AddUser(&repository.AddUserParams{
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

	userId, role, err := s.repository.CheckUser(&repository.CheckUserParams{
		Email:    req.Email,
		Password: helpers.GetMD5Hash(req.Password),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.JSON(http.StatusNotFound, &ErrorResponse{ErrorMessage: "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	token, err := s.tokenService.GetToken(userId, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &LoginResponse{Token: token})
}

func (s *ServiceImpl) GetUserByIDHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}
	id := c.Param("id")

	user, err := s.repository.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &user)
}

func (s *ServiceImpl) GetUserByTokenHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	user, err := s.repository.GetUserById(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &user)
}

func (s *ServiceImpl) CardsHandler(c echo.Context) error {
	// check authorization
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.repository.GetCards()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) CardsByIdHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.repository.GetCardsByUserId(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) CardsByBookDate(c echo.Context) error {
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	cards, err := s.repository.GetBookedCards(token.UserId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &cards)
}

func (s *ServiceImpl) UsersHandler(c echo.Context) error {
	// check authorization
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	users, err := s.repository.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}
	return c.JSON(http.StatusOK, &users)
}

func (s *ServiceImpl) GetSpecializationByIDHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &IDRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	spec, err := s.repository.GetSpecializationById(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &spec)
}

func (s *ServiceImpl) GetAddInfByIDHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &IDRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	addInf, err := s.repository.GetAddInfById(req.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &addInf)
}

func (s *ServiceImpl) UpdateSpecHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateSpecRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.repository.UpdateSpec(token.Role, &repository.UpdateSpecialization{
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
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateAddInfRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.repository.UpdateAddInf(token.Role, &repository.UpdateAddInf{
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
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateCreatorStatusRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	err = s.repository.UpdateCreatorStatus(&repository.UpdateCreatorStatusParams{
		UserID:   req.UserID,
		UserName: req.UserName,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) UpdateBookDateUserHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	_, err := s.tokenService.ParseToken(authToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, &EmptyResponse{})
	}

	req := &UpdateBookDateUserRequest{}
	err = c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{ErrorMessage: err.Error()})
	}

	fmt.Println(req.UserID)
	fmt.Println(req.BookID)
	err = s.repository.UpdateBookDatesUser(&repository.UpdateBookDateUserParams{
		UserID: req.UserID,
		BookID: req.BookID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}

func (s *ServiceImpl) UpdateCardsHandler(c echo.Context) error {
	authToken := s.getAuthToken(c)
	token, err := s.tokenService.ParseToken(authToken)
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
	for _, t := range req.Tags {
		tags = append(tags, t.Name)
	}
	for _, d := range req.BookDates {
		bookdates = append(bookdates, d.Date)
	}

	err = s.repository.AddCard(token.Role, &repository.UpdateCard{
		CardID:      req.CardID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		Cost:        req.Cost,
		Tags:        tags,
		BookDates:   bookdates,
		IsAgreement: req.IsAgreement,
		Prepayment:  req.Prepayment,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{ErrorMessage: err.Error()})
	}

	return c.JSON(http.StatusOK, &EmptyResponse{})
}
