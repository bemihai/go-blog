package main

import (
	repo "blog/repo"
	db "blog/repo/postgres"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func toJson(v interface{}) io.Reader {
	json, _ := json.Marshal(v)
	return bytes.NewReader(json)
}

func TestAddArticle(t *testing.T) {

	author := repo.Author{
		Id:    "b4a4de9e-2f52-4cf1-8907-3d828d403126",
		Name:  "test",
		Email: "test@email.com",
	}
	article := repo.Article{
		Title:  "test",
		Body:   "test",
		Author: author,
	}
	expectedArticleId := "b4a4de9e-2f52-4cf1-8907-3d828d403127"

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
		json.Unmarshal(res.Body.Bytes(), &id)
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
				return repo.Author{}, db.ErrAuthorNotFound
			},
		}

		h := BlogServer{r}
		req := httptest.NewRequest(http.MethodPost, "/articles", toJson(article))
		res := httptest.NewRecorder()

		h.AddArticle(res, req)
		require.Equal(t, res.Code, http.StatusServiceUnavailable)
	})
}
