package handler

import (
	"database/sql"
	"encoding/json"
	model "finalexam/Model"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func CreateUserValidation(w http.ResponseWriter, r *http.Request, user model.CreateUser) bool {
	//field validation
	var errorResponses model.DataResponses[*model.IError]
	errorResponses.Status = 401
	errorResponses.Message = "Bad Request"
	errVal := Validator.Struct(user)
	if errVal != nil {
		for _, err := range errVal.(validator.ValidationErrors) {
			var el model.IError
			el.Field = err.Field()
			el.Message = err.Error()
			el.Value = err.Param()
			errorResponses.Data = append(errorResponses.Data, &el)
		}
		json.NewEncoder(w).Encode(errorResponses)
		return true
	}

	// database unique validation
	searchedUsers := []model.CreateUser{}
	searchQuery := "SELECT username, email FROM [User] WHERE username = @username OR email = @email"
	searched, err2 := db.Query(searchQuery, sql.Named("username", user.Username), sql.Named("email", user.Email))
	if err2 != nil {
		panic(err2)
	}

	for searched.Next() {
		var searchUser model.CreateUser
		err := searched.Scan(&searchUser.Username, &searchUser.Email)
		if err != nil {
			panic(err)
		} else {
			var el model.IError
			if searchUser.Username == user.Username {
				el.Field = "Username"
				el.Message = "Username is taken"
				el.Value = searchUser.Username
				errorResponses.Data = append(errorResponses.Data, &el)
			} else if searchUser.Email == user.Email {
				el.Field = "Email"
				el.Message = "Email is taken"
				el.Value = searchUser.Email
				errorResponses.Data = append(errorResponses.Data, &el)
			}
		}
		searchedUsers = append(searchedUsers, searchUser)
	}
	if len(searchedUsers) > 0 {
		json.NewEncoder(w).Encode(errorResponses)
		return true
	}
	//end of validation

	return false
}

func UpdateUserValidation(w http.ResponseWriter, r *http.Request, user model.UpdateUser, id int) bool {
	//field validation
	var errorResponses model.DataResponses[*model.IError]
	errorResponses.Status = 401
	errorResponses.Message = "Bad Request"
	errVal := Validator.Struct(user)
	if errVal != nil {
		for _, err := range errVal.(validator.ValidationErrors) {
			var el model.IError
			el.Field = err.Field()
			el.Message = err.Error()
			el.Value = err.Param()
			errorResponses.Data = append(errorResponses.Data, &el)
		}
		json.NewEncoder(w).Encode(errorResponses)
		return true
	}

	// database unique validation
	searchedUsers := []model.UpdateUser{}
	searchQuery := "SELECT id, username, email FROM [User] WHERE username = @username OR email = @email OR id = @id"
	searched, err2 := db.Query(searchQuery,
		sql.Named("username", user.Username),
		sql.Named("email", user.Email),
		sql.Named("id", id))
	if err2 != nil {
		panic(err2)
	}

	for searched.Next() {
		var searchUser model.UpdateUser
		err := searched.Scan(&searchUser.Id, &searchUser.Username, &searchUser.Email)
		if err != nil {
			panic(err)
		} else {
			var el model.IError
			var updateId bool
			if searchUser.Id != id {
				if searchUser.Username == user.Username {
					el.Field = "Username"
					el.Message = "Username is taken"
					el.Value = searchUser.Username
					errorResponses.Data = append(errorResponses.Data, &el)
				} else if searchUser.Email == user.Email {
					el.Field = "Email"
					el.Message = "Email is taken"
					el.Value = searchUser.Email
					errorResponses.Data = append(errorResponses.Data, &el)
				}
			} else {
				updateId = true
			}

			if !updateId {
				var errorIdResponses model.DataResponse[model.Message]
				errorIdResponses.Status = 401
				errorIdResponses.Message = "Bad Request"
				errorIdResponses.Data.Message = "ID not found"
				json.NewEncoder(w).Encode(errorIdResponses)
				return true
			}
		}
		searchedUsers = append(searchedUsers, searchUser)
	}
	if len(searchedUsers) > 0 {
		json.NewEncoder(w).Encode(errorResponses)
		return true
	}
	//end of validation

	return false
}

func ValidateFields(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	var errorResponses model.DataResponses[*model.IError]
	errorResponses.Status = 401
	errorResponses.Message = "Bad Request"

	var validateData = model.IValidate{CustomStruct: data}

	errVal := Validator.Struct(validateData)
	if errVal != nil {
		for _, err := range errVal.(validator.ValidationErrors) {
			var el model.IError
			el.Field = err.Field()
			el.Message = err.Error()
			el.Value = err.Param()
			errorResponses.Data = append(errorResponses.Data, &el)
		}
		json.NewEncoder(w).Encode(errorResponses)
		return true
	}
	return false
}

func DeleteValidation(w http.ResponseWriter, r *http.Request, action string, id int) bool {
	var query string

	switch action {
	case "User":
		query = "SELECT id FROM [User] WHERE Id = @Id"
	case "SocialMedia":
		query = "SELECT id FROM SocialMedia WHERE Id = @Id"
	case "Photo":
		query = "SELECT id FROM Photo WHERE Id = @Id"
	case "Comment":
		query = "SELECT id FROM Comment WHERE Id = @Id"
	}

	searched, err2 := db.Query(query,
		sql.Named("id", id))
	if err2 != nil {
		panic(err2)
	}

	if !searched.Next() {
		var errorIdResponses model.DataResponse[model.Message]
		errorIdResponses.Status = 401
		errorIdResponses.Message = "Bad Request"
		errorIdResponses.Data.Message = "Delete fail : ID not found"
		json.NewEncoder(w).Encode(errorIdResponses)
		return true
	}
	return false
}
