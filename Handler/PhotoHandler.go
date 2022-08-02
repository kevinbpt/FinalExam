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
)

func PhotoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var tempId, _ = strconv.Atoi(id)
	switch r.Method {
	case http.MethodPost:
		CreatePhoto(w, r)
	case http.MethodPut:
		UpdatePhoto(w, r, tempId)
	case http.MethodDelete:
		DeletePhoto(w, r, tempId)
	case http.MethodGet:
		GetPhoto(w, r)
	}
}

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, caption, photo_url, user_id, created_at, updated_at FROM Photo"
	var data model.DataResponses[model.PhotoResponse]
	searched, err2 := db.Query(query)
	if err2 != nil {
		panic(err2)
	}

	for searched.Next() {
		var photo model.PhotoResponse
		err := searched.Scan(&photo.Id, &photo.Title, &photo.Caption, &photo.PhotoUrl, &photo.UserId, &photo.CreatedAt, &photo.UpdatedAt)
		if err != nil {
			panic(err)
		}
		query2 := "SELECT username, email FROM [User] WHERE Id = @Id"
		searched2, err2 := db.Query(query2,
			sql.Named("id", photo.UserId))
		if err2 != nil {
			panic(err2)
		}

		for searched2.Next() {
			var user model.ResponseUser
			err3 := searched2.Scan(&user.Username, &user.Email)
			if err3 != nil {
				panic(err3)
			}
			photo.User = user
		}
		data.Data = append(data.Data, photo)
	}

	data.Status = 200
	data.Message = "Successful"
	json.NewEncoder(w).Encode(data)
}

func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	var photo = &model.Photo{}
	var response model.DataResponse[model.Photo]
	var timenow = time.Now()
	if err := json.NewDecoder(r.Body).Decode(photo); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *photo) {
			return
		}
		query := "INSERT INTO Photo (title, caption, photo_url, user_id, created_at, updated_at) VALUES(@title, @caption, @photo_url, @user_id, @created_at, @updated_at)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err := db.ExecContext(ctx, query,
			sql.Named("title", photo.Title),
			sql.Named("caption", photo.Caption),
			sql.Named("photo_url", photo.PhotoUrl),
			sql.Named("user_id", photo.UserId),
			sql.Named("created_at", timenow),
			sql.Named("updated_at", timenow))
		if err != nil {
			log.Fatal(err)
		}
		response.Status = 200
		response.Message = "Photo added successfully"
		photo.CreatedAt = timenow
		response.Data = *photo
		json.NewEncoder(w).Encode(response)
	}

}

func UpdatePhoto(w http.ResponseWriter, r *http.Request, id int) {
	var photo = &model.Photo{}
	var response model.DataResponse[model.Photo]
	var timenow = time.Now()
	if err := json.NewDecoder(r.Body).Decode(photo); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *photo) {
			return
		}
		query := "UPDATE Photo SET title = @title, caption = @caption, photo_url = @photo_url, updated_at = @updated_at WHERE id = @id"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		_, err := db.ExecContext(ctx, query,
			sql.Named("title", photo.Title),
			sql.Named("caption", photo.Caption),
			sql.Named("photo_url", photo.PhotoUrl),
			sql.Named("updated_at", timenow),
			sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}

		response.Status = 200
		response.Message = "Social Media updated successfully"
		photo.UpdatedAt = timenow
		response.Data = *photo
		json.NewEncoder(w).Encode(response)
	}
}

func DeletePhoto(w http.ResponseWriter, r *http.Request, id int) {
	if DeleteValidation(w, r, "Photo", id) {
		return
	}

	ctx, cancelfunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelfunc()

	var response model.DataResponse[model.Message]

	_, err := db.ExecContext(ctx, "DELETE FROM Photo WHERE id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
	}
	response.Status = 200
	response.Data.Message = "Your photo has been successfully deleted"
	json.NewEncoder(w).Encode(response)
}
