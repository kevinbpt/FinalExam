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

func SocmedHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var tempId, _ = strconv.Atoi(id)
	switch r.Method {
	case http.MethodPost:
		CreateSocialMedia(w, r)
	case http.MethodPut:
		UpdateSocialMedia(w, r, tempId)
	case http.MethodDelete:
		DeleteSocialMedia(w, r, tempId)
	case http.MethodGet:
		GetSocialMedia(w, r)
	}
}

func GetSocialMedia(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, name, social_media_url, UserId FROM SocialMedia"
	var data model.DataResponses[model.SocialMedia]
	searched, err2 := db.Query(query)
	if err2 != nil {
		panic(err2)
	}

	for searched.Next() {
		var socmed model.SocialMedia
		err := searched.Scan(&socmed.Id, &socmed.Name, &socmed.SocialMediaUrl, &socmed.UserId)
		if err != nil {
			panic(err)
		}
		data.Data = append(data.Data, socmed)
	}

	data.Status = 200
	data.Message = "Successful"
	json.NewEncoder(w).Encode(data)
}

func CreateSocialMedia(w http.ResponseWriter, r *http.Request) {
	var socmed = &model.SocialMedia{}
	var response model.DataResponse[model.SocialMedia]

	if err := json.NewDecoder(r.Body).Decode(socmed); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *socmed) {
			return
		}
		query := "INSERT INTO SocialMedia (name, social_media_url, UserId) VALUES(@name, @social_media_url, @UserId)"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()
		_, err := db.ExecContext(ctx, query,
			sql.Named("name", socmed.Name),
			sql.Named("social_media_url", socmed.SocialMediaUrl),
			sql.Named("UserId", socmed.UserId))
		if err != nil {
			log.Fatal(err)
		}
		response.Status = 200
		response.Message = "Social Media added successfully"
		response.Data = *socmed
		json.NewEncoder(w).Encode(response)
	}

}

func UpdateSocialMedia(w http.ResponseWriter, r *http.Request, id int) {
	var socmed = &model.SocialMedia{}
	var response model.DataResponse[model.SocialMedia]
	if err := json.NewDecoder(r.Body).Decode(socmed); err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	} else {
		if ValidateFields(w, r, *socmed) {
			return
		}
		query := "UPDATE SocialMedia SET name = @name, social_media_url = @social_media_url WHERE id = @id"
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		_, err := db.ExecContext(ctx, query,
			sql.Named("name", socmed.Name),
			sql.Named("social_media_url", socmed.SocialMediaUrl),
			sql.Named("id", id))
		if err != nil {
			log.Fatal(err)
		}

		response.Status = 200
		response.Message = "Social Media updated successfully"
		response.Data = *socmed
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteSocialMedia(w http.ResponseWriter, r *http.Request, id int) {
	if DeleteValidation(w, r, "SocialMedia", id) {
		return
	}

	ctx, cancelfunc := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancelfunc()

	var response model.DataResponse[model.Message]

	_, err := db.ExecContext(ctx, "DELETE FROM SocialMedia WHERE id=@id",
		sql.Named("id", id))

	if err != nil {
		log.Fatal(err)
	}
	response.Status = 200
	response.Data.Message = "Your social media has been successfully deleted"
	json.NewEncoder(w).Encode(response)
}
