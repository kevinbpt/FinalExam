package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	model "finalexam/Model"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var tempId, _ = strconv.Atoi(id)
	switch r.Method {
	case http.MethodPost:
		CreateUser(w, r)
	case http.MethodPut:
		UpdateUser(w, r, tempId)
	case http.MethodDelete:
		DeleteUser(w, r, tempId)
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var auth = &model.Auth{}
	if err := json.NewDecoder(r.Body).Decode(auth); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		var user model.Auth
		var token model.DataResponse[model.Token]

		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		query := "SELECT password FROM [User] WHERE email= @email"
		err := db.QueryRowContext(ctx, query, sql.Named("email", auth.Email)).Scan(&user.Password)
		if err != nil {
			token.Status = 400
			token.Message = err.Error()
		}

		err2 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
		if err2 != nil {
			token.Status = 400
			token.Message = err2.Error()
		}

		validToken, err3 := GenerateJWT(auth.Email)
		if err3 != nil {
			token.Status = 401
			token.Message = err3.Error()
		}
		token.Status = 200
		token.Message = "Success"
		token.Data.Token = validToken
		json.NewEncoder(w).Encode(token)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user = &model.CreateUser{}
	var response model.DataResponse[model.CreateUser]
	var timenow = time.Now()

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if CreateUserValidation(w, r, *user) {
			return
		}
		hashed, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if errHash != nil {
			log.Fatal(errHash)
		}
		query := "INSERT INTO dbo.[User] (username, password, email, age, created_at, updated_at) VALUES(@username, @password, @email, @age, @created_at, @updated_at)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err := db.ExecContext(ctx, query,
			sql.Named("username", user.Username),
			sql.Named("password", string(hashed)),
			sql.Named("email", user.Email),
			sql.Named("age", user.Age),
			sql.Named("created_at", timenow),
			sql.Named("updated_at", timenow))
		if err != nil {
			log.Fatal(err)
		}
		response.Status = 200
		response.Message = "User added successfully"
		response.Data = *user
		response.Data.Password = string(hashed)
		json.NewEncoder(w).Encode(response)
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var user = &model.UpdateUser{}
	var response model.DataResponse[model.UpdateUser]
	var timenow = time.Now()
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if UpdateUserValidation(w, r, *user, id) {
			return
		}
		query := "UPDATE dbo.[User] SET username = @username, email = @email, updated_at = @updated_at WHERE id = @id"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		_, err := db.ExecContext(ctx, query,
			sql.Named("username", user.Username),
			sql.Named("email", user.Email),
			sql.Named("updated_at", timenow),
			sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}

		response.Status = 200
		response.Message = "User updated successfully"
		response.Data = *user
		response.Data.UpdatedAt = timenow
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if DeleteValidation(w, r, "User", id) {
		return
	}

	ctx, cancelfunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelfunc()

	var response model.DataResponse[model.Message]

	_, err := db.ExecContext(ctx, "DELETE FROM dbo.[User] WHERE id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
	}
	response.Status = 200
	response.Data.Message = "Your account has been successfully deleted"
	json.NewEncoder(w).Encode(response)
}
