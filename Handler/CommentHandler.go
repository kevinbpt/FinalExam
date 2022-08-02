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

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var tempId, _ = strconv.Atoi(id)
	switch r.Method {
	case http.MethodPost:
		CreateComment(w, r)
	case http.MethodPut:
		UpdateComment(w, r, tempId)
	case http.MethodDelete:
		DeleteComment(w, r, tempId)
	case http.MethodGet:
		GetComment(w, r)
	}
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, user_id, photo_id, message, created_at, updated_at FROM Comment"
	var data model.DataResponses[model.CommentResponse]
	searched, err2 := db.Query(query)
	if err2 != nil {
		panic(err2)
	}

	for searched.Next() {
		var comment model.CommentResponse
		err := searched.Scan(&comment.Id, &comment.UserId, &comment.PhotoId, &comment.Message, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			panic(err)
		}

		//get USER
		query2 := "SELECT id, username, email FROM [User] WHERE Id = @Id"
		searched2, err2 := db.Query(query2,
			sql.Named("id", comment.UserId))
		if err2 != nil {
			panic(err2)
		}
		for searched2.Next() {
			var user model.ResponseUser
			err3 := searched2.Scan(&user.Id, &user.Username, &user.Email)
			if err3 != nil {
				panic(err3)
			}
			comment.User = user
		}
		//get PHOTO
		query3 := "SELECT id, title, caption, photo_url, user_id FROM Photo WHERE id = @id"
		searched3, err3 := db.Query(query3,
			sql.Named("id", comment.PhotoId))
		if err3 != nil {
			panic(err3)
		}
		for searched3.Next() {
			var photo model.Photo
			err3 := searched3.Scan(&photo.Id, &photo.Title, &photo.Caption, &photo.PhotoUrl, &photo.UserId)
			if err3 != nil {
				panic(err3)
			}
			comment.Photo = photo
		}

		data.Data = append(data.Data, comment)
	}

	data.Status = 200
	data.Message = "Successful"
	json.NewEncoder(w).Encode(data)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment = &model.Comment{}
	var response model.DataResponse[model.Comment]
	var timenow = time.Now()
	if err := json.NewDecoder(r.Body).Decode(comment); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *comment) {
			return
		}
		query := "INSERT INTO Comment (user_id, photo_id, message, created_at, updated_at) VALUES(@user_id, @photo_id, @message, @created_at, @updated_at)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err := db.ExecContext(ctx, query,
			sql.Named("user_id", comment.UserId),
			sql.Named("photo_id", comment.PhotoId),
			sql.Named("message", comment.Message),
			sql.Named("created_at", timenow),
			sql.Named("updated_at", timenow))
		if err != nil {
			log.Fatal(err)
		}
		response.Status = 200
		response.Message = "Comment added successfully"
		comment.CreatedAt = timenow
		response.Data = *comment
		json.NewEncoder(w).Encode(response)
	}

}

func UpdateComment(w http.ResponseWriter, r *http.Request, id int) {
	var comment = &model.Comment{}
	var response model.DataResponse[model.Comment]
	var timenow = time.Now()
	if err := json.NewDecoder(r.Body).Decode(comment); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *comment) {
			return
		}
		query := "UPDATE Comment SET message = @message, updated_at = @updated_at WHERE id = @id"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		_, err := db.ExecContext(ctx, query,
			sql.Named("message", comment.Message),
			sql.Named("updated_at", timenow),
			sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}

		response.Status = 200
		response.Message = "Comment updated successfully"
		comment.UpdatedAt = timenow
		response.Data = *comment
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request, id int) {
	if DeleteValidation(w, r, "Comment", id) {
		return
	}

	ctx, cancelfunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelfunc()

	var response model.DataResponse[model.Message]

	_, err := db.ExecContext(ctx, "DELETE FROM Comment WHERE id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
	}
	response.Status = 200
	response.Data.Message = "Your comment has been successfully deleted"
	json.NewEncoder(w).Encode(response)
}
