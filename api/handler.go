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

// BlogServer is responsible to answer to http request.
type BlogServer struct {
	Service repo.BlogService
}

func (h *BlogServer) ListArticles(w http.ResponseWriter, r *http.Request) {

	// get articles
	articles, err := h.Service.ListArticles()
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
	authors, err := h.Service.GetAuthorsByIds(ids)
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

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

func (h *BlogServer) GetArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Bad request: id is not a valid uuid.", http.StatusBadRequest)
		return
	}

	q := id.String()
	article, err := h.Service.GetArticleById(q)
	if err != nil {
		if errors.Is(err, db.ErrArticleNotFound) {
			http.Error(w, "Article not found.", http.StatusNotFound)
			return
		}
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}

	// get article's author
	article.Author, err = h.Service.GetAuthorById(article.Author.Id)
	if err != nil {
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}

	data, err := json.Marshal(article)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

func (h *BlogServer) AddArticle(w http.ResponseWriter, r *http.Request) {

	var article repo.Article

	// decode the request body into an Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Bad request: body is not correct.", http.StatusBadRequest)
		return
	}

	// add Author field in blog.authors table if not already exists
	author, err := h.Service.GetAuthorByNameAndEmail(article.Author.Name, article.Author.Email)
	if err != nil {
		if errors.Is(err, db.ErrAuthorNotFound) {
			article.Author.Id, err = h.Service.AddAuthor(article.Author)
			if err != nil {
				http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
				return
			}
		}
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	} else {
		article.Author.Id = author.Id
	}

	// add Article in blog.articles table
	a, err := h.Service.AddArticle(article)
	if err != nil {
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}

	data, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}

func (h *BlogServer) DeleteArticleById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Bad request: id is not a valid uuid.", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteArticleById(id.String())
	if err != nil {
		if errors.Is(err, db.ErrArticleNotFound) {
			http.Error(w, "Article not found.", http.StatusNotFound)
			return
		}
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}
}

func (h *BlogServer) DeleteAuthorByNameAndEmail(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	email := r.FormValue("email")

	err := h.Service.DeleteAuthorByNameAndEmail(name, email)
	if err != nil {
		if errors.Is(err, db.ErrAuthorNotFound) {
			http.Error(w, "Author not found.", http.StatusNotFound)
			return
		}
		http.Error(w, "Service unavailable.", http.StatusServiceUnavailable)
		return
	}
}

// MethodNotAllowed handles not allowed requests on existing endpoints.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("Method not allowed."))
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
