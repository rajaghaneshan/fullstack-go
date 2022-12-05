package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rajaghaneshan/go-fullstack/api/auth"
	"github.com/rajaghaneshan/go-fullstack/api/models"
	"github.com/rajaghaneshan/go-fullstack/api/responses"
	"github.com/rajaghaneshan/go-fullstack/api/utils/formaterror"
)

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	post.Prepare()
	err = post.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formaterror := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formaterror)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}

	postReceived, err := post.FindPostById(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post := models.Post{}
	//check if post exists
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	//check if post belongs to the user
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	//read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//process request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//request user is same as ppost author
	if uid != postUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	postUpdate.Prepare()
	err = postUpdate.Validate()
	if err != nil{
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.ID = post.ID
	postUpdated, err := postUpdate.UpdateAPost(server.DB)
	if err != nil{
		formaterror := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formaterror)
		return
	}

	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil{
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	post:= models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil{
		responses.ERROR(w, http.StatusNotFound, errors.New("post not found"))
		return
	}

	if uid != post.AuthorID{
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil{
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w,http.StatusNoContent, "")
}