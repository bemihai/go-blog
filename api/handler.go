package main

import (
	"blog"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler is responsible to answer to http request.
type Handler struct {
	Repository ArticlesRepository
}

// Articles repository interface.
type ArticlesRepository interface {
	ListArticles() ([]blog.Article, error)
	GetArticleById(id int) (blog.Article, error)
	PostArticle(a blog.Article) (blog.Article, error)
	DeleteArticleById(id int) (blog.Article, error)
	DeleteAuthorByNameAndEmail(name string, email string) (blog.Author, error)
}

func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {

	articles, err := h.Repository.ListArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (h *Handler) GetArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.Repository.GetArticleById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (h *Handler) PostArticle(w http.ResponseWriter, r *http.Request) {

	var article blog.Article

	// decode the request body into the Article struct
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a, err := h.Repository.PostArticle(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	data, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (h *Handler) DeleteArticleById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.Repository.DeleteArticleById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (h *Handler) DeleteAuthorByNameAndEmail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	email := r.FormValue("email")

	fmt.Printf("name: %s, email: %s ", name, email)

	article, err := h.Repository.DeleteAuthorByNameAndEmail(name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

// Handles not allowed requests on existing endpoints.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}
