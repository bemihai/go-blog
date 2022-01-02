package api

import (
	repo "blog/repo"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var author = repo.Author{
	Id:    "b4a4de9e-2f52-4cf1-8907-3d828d403126",
	Name:  "test",
	Email: "test@email.com",
}
var article = repo.Article{
	Title:  "test",
	Body:   "test",
	Author: author,
}
var expectedArticleId = "b4a4de9e-2f52-4cf1-8907-3d828d403127"
var expectedAuthorId = "b4a4de9e-2f52-4cf1-8907-3d828d403126"

func toJson(v interface{}) io.Reader {
	json, _ := json.Marshal(v)
	return bytes.NewReader(json)
}

func TestListArticles(t *testing.T) {

	t.Run("can get all articles", func(t *testing.T) {
		r := &MockService{
			GetAuthorsByIdsFunc: func(ids []string) ([]repo.Author, error) {
				return []repo.Author{author}, nil
			},
			ListArticlesFunc: func() ([]repo.Article, error) {
				return []repo.Article{article}, nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		res := httptest.NewRecorder()
		h.ListArticles(res, req)

		var a []repo.Article
		json.Unmarshal(res.Body.Bytes(), &a) // nolint: errcheck

		require.Equal(t, res.Code, http.StatusOK)
		require.Len(t, a, 1)
		require.Equal(t, a[0].Title, article.Title)
		require.Equal(t, a[0].Author.Name, article.Author.Name)
	})

	t.Run("return 503 if get articles fails", func(t *testing.T) {
		r := &MockService{
			ListArticlesFunc: func() ([]repo.Article, error) {
				return []repo.Article{}, errors.New("couldn't fetch articles")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		res := httptest.NewRecorder()
		h.ListArticles(res, req)

		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})

	t.Run("return 503 if get authors fails", func(t *testing.T) {
		r := &MockService{
			GetAuthorsByIdsFunc: func(ids []string) ([]repo.Author, error) {
				return []repo.Author{}, errors.New("couldn't fetch authors")
			},
			ListArticlesFunc: func() ([]repo.Article, error) {
				return []repo.Article{article}, nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		res := httptest.NewRecorder()
		h.ListArticles(res, req)

		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}

func TestGetArticleById(t *testing.T) {

	t.Run("can get article by id", func(t *testing.T) {
		r := &MockService{
			GetAuthorByIdFunc: func(id string) (repo.Author, error) {
				require.Equal(t, id, expectedAuthorId)
				return author, nil
			},
			GetArticleByIdFunc: func(id string) (repo.Article, error) {
				require.Equal(t, id, expectedArticleId)
				return article, nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.GetArticleById(res, req)
		var a repo.Article
		json.Unmarshal(res.Body.Bytes(), &a) // nolint: errcheck

		require.Equal(t, res.Code, http.StatusOK)
		require.Equal(t, a.Title, article.Title)
		require.Equal(t, a.Author.Name, article.Author.Name)
	})

	t.Run("return 400 when id is invalid uuid", func(t *testing.T) {
		r := &MockService{}
		h := BlogServer{r}

		req := httptest.NewRequest(http.MethodGet, "/articles/id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "id"})
		res := httptest.NewRecorder()

		h.GetArticleById(res, req)
		require.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("return 404 when article not found", func(t *testing.T) {
		r := &MockService{
			GetArticleByIdFunc: func(id string) (repo.Article, error) {
				require.Equal(t, id, expectedArticleId)
				return repo.Article{}, repo.ErrArticleNotFound
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.GetArticleById(res, req)
		require.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("return 503 when get article fails", func(t *testing.T) {
		r := &MockService{
			GetArticleByIdFunc: func(id string) (repo.Article, error) {
				require.Equal(t, id, expectedArticleId)
				return repo.Article{}, errors.New("couldn't fetch article")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.GetArticleById(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})

	t.Run("return 503 when get author fails", func(t *testing.T) {
		r := &MockService{
			GetArticleByIdFunc: func(id string) (repo.Article, error) {
				require.Equal(t, id, expectedArticleId)
				return article, nil
			},
			GetAuthorByIdFunc: func(id string) (repo.Author, error) {
				require.Equal(t, id, expectedAuthorId)
				return repo.Author{}, errors.New("couldn't fetch author")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.GetArticleById(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}

func TestAddArticle(t *testing.T) {

	t.Run("can add valid article", func(t *testing.T) {

		r := &MockService{
			AddAuthorFunc: func(a repo.Author) (string, error) {
				return author.Id, nil
			},
			AddArticleFunc: func(a repo.Article) (string, error) {
				return expectedArticleId, nil
			},
			GetAuthorByNameAndEmailFunc: func(name, email string) (repo.Author, error) {
				return author, nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodPost, "/articles", toJson(article))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		var id string
		json.Unmarshal(res.Body.Bytes(), &id) // nolint: errcheck

		require.Equal(t, res.Code, http.StatusOK)
		require.Equal(t, id, expectedArticleId)
		require.Len(t, r.Articles, 1)
		require.Equal(t, r.Articles[0], article)
	})

	t.Run("returns 400 if body is not valid article JSON", func(t *testing.T) {

		h := BlogServer{nil}
		req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader("invalid json"))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		require.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("returns 503 if add article fails", func(t *testing.T) {

		r := &MockService{
			AddAuthorFunc: func(a repo.Author) (string, error) {
				return author.Id, nil
			},
			AddArticleFunc: func(a repo.Article) (string, error) {
				return "", errors.New("couldn't add new article")
			},
			GetAuthorByNameAndEmailFunc: func(name, email string) (repo.Author, error) {
				return author, nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodPost, "/articles", toJson(article))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})

	t.Run("returns 503 if get author by name and email fails", func(t *testing.T) {

		r := &MockService{
			GetAuthorByNameAndEmailFunc: func(name, email string) (repo.Author, error) {
				return repo.Author{}, errors.New("couldn't get author by name and email")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodPost, "/articles", toJson(article))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})

	t.Run("returns 503 if add author fails", func(t *testing.T) {

		r := &MockService{
			AddAuthorFunc: func(a repo.Author) (string, error) {
				return "", errors.New("couldn't add new author")
			},
			GetAuthorByNameAndEmailFunc: func(name, email string) (repo.Author, error) {
				return repo.Author{}, repo.ErrAuthorNotFound
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodPost, "/articles", toJson(article))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}

func TestDeleteArticleById(t *testing.T) {

	t.Run("return 400 when id is invalid uuid", func(t *testing.T) {
		r := &MockService{}
		h := BlogServer{r}

		req := httptest.NewRequest(http.MethodDelete, "/articles/id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "id"})
		res := httptest.NewRecorder()

		h.DeleteArticleById(res, req)
		require.Equal(t, res.Code, http.StatusBadRequest)
	})

	t.Run("can delete article by id", func(t *testing.T) {
		r := &MockService{
			DeleteArticleByIdFunc: func(id string) error {
				require.Equal(t, id, expectedArticleId)
				return nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.DeleteArticleById(res, req)
		require.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("return 404 if article not found", func(t *testing.T) {
		r := &MockService{
			DeleteArticleByIdFunc: func(id string) error {
				require.Equal(t, id, expectedArticleId)
				return repo.ErrArticleNotFound
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.DeleteArticleById(res, req)
		require.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("return 503 if service fails", func(t *testing.T) {
		r := &MockService{
			DeleteArticleByIdFunc: func(id string) error {
				require.Equal(t, id, expectedArticleId)
				return errors.New("service fails")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles/%s", expectedArticleId), nil)
		req = mux.SetURLVars(req, map[string]string{"id": expectedArticleId})
		res := httptest.NewRecorder()

		h.DeleteArticleById(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}

func TestDeleteAuthorByNameAndEmail(t *testing.T) {

	t.Run("can delete article by name and email", func(t *testing.T) {
		r := &MockService{
			DeleteAuthorByNameAndEmailFunc: func(name string, email string) error {
				require.Equal(t, name, author.Name)
				require.Equal(t, email, author.Email)
				return nil
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles?name=%s&email=%s", author.Name, author.Email), nil)
		res := httptest.NewRecorder()

		h.DeleteAuthorByNameAndEmail(res, req)
		require.Equal(t, res.Code, http.StatusOK)
	})

	t.Run("return 404 if article not found", func(t *testing.T) {
		r := &MockService{
			DeleteAuthorByNameAndEmailFunc: func(name string, email string) error {
				require.Equal(t, name, author.Name)
				require.Equal(t, email, author.Email)
				return repo.ErrAuthorNotFound
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles?name=%s&email=%s", author.Name, author.Email), nil)
		res := httptest.NewRecorder()

		h.DeleteAuthorByNameAndEmail(res, req)
		require.Equal(t, res.Code, http.StatusNotFound)
	})

	t.Run("return 503 if service fails", func(t *testing.T) {
		r := &MockService{
			DeleteAuthorByNameAndEmailFunc: func(name string, email string) error {
				require.Equal(t, name, author.Name)
				require.Equal(t, email, author.Email)
				return errors.New("service fails")
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/articles?name=%s&email=%s", author.Name, author.Email), nil)
		res := httptest.NewRecorder()

		h.DeleteAuthorByNameAndEmail(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}

func TestMethodNotAllowed(t *testing.T) {
	t.Run("returns 405 if method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/articles", nil)
		res := httptest.NewRecorder()
		MethodNotAllowed(res, req)
		require.Equal(t, res.Code, http.StatusMethodNotAllowed)
	})
}
