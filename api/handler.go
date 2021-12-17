package main

import (
	repo "blog/repo"
	db "blog/repo/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler is responsible to answer to http request.
type Handler struct {
	Repository repo.Repository
}

func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {

	// get articles
	articles, err := h.Repository.ListArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// for each article, get its author
	for i := range articles {
		auth, err := h.Repository.GetAuthorById(articles[i].Author.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		articles[i].Author.Name = auth.Name
		articles[i].Author.Email = auth.Email
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
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.Repository.GetArticleById(id)
	var dberr db.DatabaseError
	if err != nil {
		if errors.As(err, &dberr) {
			http.Error(w, dberr.Message, http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func (h *Handler) AddArticle(w http.ResponseWriter, r *http.Request) {

	var article repo.Article

	// decode the request body into an Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add Author field in blog.authors table
	article.Author.Id, err = h.Repository.AddAuthor(article.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	// add Article in blog.articles table
	a, err := h.Repository.AddArticle(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
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
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repository.DeleteArticleById(id)
	var dberr db.DatabaseError
	if err != nil {
		if errors.As(err, &dberr) {
			http.Error(w, dberr.Message, http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

func (h *Handler) DeleteAuthorByNameAndEmail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	email := r.FormValue("email")

	fmt.Printf("name: %s, email: %s ", name, email)

	err := h.Repository.DeleteAuthorByNameAndEmail(name, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// Handles not allowed requests on existing endpoints.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}
