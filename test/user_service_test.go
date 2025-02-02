package test

/*func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.User{})
	return db
}

func TestLoginIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	service := user_service.NewUserService(db)

	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		ctx := &utils.Context{Ctx: c}
		var request user_service.LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, werr := service.Login(ctx, request)
		if werr != nil {
			c.JSON(werr.Code, gin.H{"error": werr.Err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	t.Run("New user login", func(t *testing.T) {
		loginRequest := user_service.LoginRequest{PhoneNumber: "1234567890"}
		jsonData, _ := json.Marshal(loginRequest)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var loginResponse user_service.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &loginResponse)

		assert.NotZero(t, loginResponse.UserID)
		assert.NotEmpty(t, loginResponse.AccessToken)
	})

	t.Run("Existing user login", func(t *testing.T) {
		db.Create(&models.User{PhoneNumber: "0987654321", CreatedAt: time.Now(), UpdatedAt: time.Now()})

		loginRequest := user_service.LoginRequest{PhoneNumber: "0987654321"}
		jsonData, _ := json.Marshal(loginRequest)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var loginResponse user_service.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &loginResponse)

		assert.NotZero(t, loginResponse.UserID)
		assert.NotEmpty(t, loginResponse.AccessToken)
	})

	t.Run("Invalid phone number", func(t *testing.T) {
		loginRequest := user_service.LoginRequest{PhoneNumber: "  "}
		jsonData, _ := json.Marshal(loginRequest)

		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid phone number")
	})
}*/
