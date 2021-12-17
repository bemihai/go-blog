package main

import (
	repo "blog/repo"
	db "blog/repo/postgres"
	"encoding/json"
	"errors"
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
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}

	// get author ids as a slice
	ids := make([]string, 0)
	for _, a := range articles {
		ids = append(ids, a.Author.Id)
	}

	// get authors by ids
	authors, err := h.Repository.GetAuthorsByIds(ids)
	if err != nil {
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}

	// define a map author_id: author
	author_map := make(map[string]repo.Author, len(authors))
	for _, a := range authors {
		author_map[a.Id] = a
	}

	// for each article, fill in the author
	for i := range articles {
		articles[i].Author = author_map[articles[i].Author.Id]
	}

	data, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (h *Handler) GetArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Bad request: id is not a valid uuid.", http.StatusBadRequest)
		return
	}

	article, err := h.Repository.GetArticleById(id.String())
	if err != nil {
		if errors.Is(err, db.ErrArticleNotFound) {
			http.Error(w, "Article not found.", http.StatusNotFound)
		}
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func (h *Handler) AddArticle(w http.ResponseWriter, r *http.Request) {

	var article repo.Article

	// decode the request body into an Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Bad request: body is not correct.", http.StatusBadRequest)
		return
	}

	// TO DO: check if author already exists: GetAuthorByNameAndEmail
	// add Author field in blog.authors table
	article.Author.Id, err = h.Repository.AddAuthor(article.Author)
	if err != nil {
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
	}

	// add Article in blog.articles table
	a, err := h.Repository.AddArticle(article)
	if err != nil {
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
	}

	data, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *Handler) DeleteArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Bad request: id is not a valid uuid.", http.StatusBadRequest)
		return
	}

	err = h.Repository.DeleteArticleById(id.String())
	if err != nil {
		if errors.Is(err, db.ErrArticleNotFound) {
			http.Error(w, "Article not found.", http.StatusNotFound)
		}
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteAuthorByNameAndEmail(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	email := r.FormValue("email")

	err := h.Repository.DeleteAuthorByNameAndEmail(name, email)
	if err != nil {
		if errors.Is(err, db.ErrArticleNotFound) {
			http.Error(w, "Article not found.", http.StatusNotFound)
		}
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

// Handles not allowed requests on existing endpoints.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed."))
}
