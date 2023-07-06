package server

func (s *ServiceImpl) setupRoutes() {
	s.server.POST("/add-user", s.AddUserHandler)
	s.server.POST("/login", s.LoginHandler)
	s.server.GET("/cards", s.CardsHandler)
	s.server.GET("/cards-by-id/:id", s.GetCardsByIDHandler)
	s.server.GET("/cards-by-token", s.CardsByTokenHandler)
	s.server.GET("/booked-cards", s.CardsByBookDate)
	s.server.GET("/cards-was-booked", s.CardsByUsersBooked)
	s.server.GET("/users", s.UsersHandler)

	s.server.GET("/user/:id", s.GetUserByIDHandler)
	s.server.GET("/user", s.GetUserByTokenHandler)
	s.server.PATCH("/spec", s.GetSpecializationByIDHandler)
	s.server.PATCH("/add-inf", s.GetAddInfByIDHandler)

	s.server.POST("/update-spec", s.UpdateSpecHandler)
	s.server.POST("/update-add-inf", s.UpdateAddInfHandler)
	s.server.POST("/update-creator-status", s.UpdateCreatorStatusHandler)
	s.server.POST("/update-book-date-user", s.UpdateBookDateUserHandler)
	s.server.POST("/update-card", s.UpdateCardsHandler)

	s.server.POST("/upload-image", s.UploadImageHandler)
	s.server.GET("/get-image", s.GetImage)

	s.server.POST("/send-push", s.SendMessageHandler)
	s.server.POST("/save-fcm-token", s.SaveFMCToken)
}
