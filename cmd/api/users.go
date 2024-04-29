package main

// func (app *application) createUser(c echo.Context) error {

// 	var User struct {
// 		Name     string `json:"name" `
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 		Phone    string `json:"phone"`
// 		Role     string `json:"role"`
// 		IsActive bool   `json:"isActive"`
// 	}

// 	if err := c.Bind(&User); err != nil {
// 		return err
// 	}

// 	user.PutTime()

// 	err := user.SetPassword()

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

// 	defer cancel()

// 	// check if email already exists
// 	var result mongodb.User

// 	err = app.models.Users.FindOne(ctx, bson.M{"email": user.Email}).Decode(&result)

// 	if err == nil {
// 		return c.JSON(http.StatusBadRequest, "email already exists")
// 	}

// 	_, err = app.models.Users.InsertOne(ctx, user)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err)
// 	}

// 	return c.JSON(http.StatusOK, user)

// }
