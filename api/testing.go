package main

import repo "blog/repo"

type MockService struct {
	ListArticlesFunc               func() ([]repo.Article, error)
	ListAuthorsFunc                func() ([]repo.Author, error)
	GetArticleByIdFunc             func(id string) (repo.Article, error)
	GetAuthorByIdFunc              func(id string) (repo.Author, error)
	GetAuthorsByIdsFunc            func(ids []string) ([]repo.Author, error)
	GetAuthorByNameAndEmailFunc    func(name string, email string) (repo.Author, error)
	AddArticleFunc                 func(a repo.Article) (string, error)
	AddAuthorFunc                  func(a repo.Author) (string, error)
	DeleteArticleByIdFunc          func(id string) error
	DeleteAuthorByIdFunc           func(id string) error
	DeleteAuthorByNameAndEmailFunc func(name string, email string) error
	Articles                       []repo.Article
	Authors                        []repo.Author
}

func (r *MockService) ListArticles() ([]repo.Article, error) {
	return r.ListArticlesFunc()
}

func (r *MockService) ListAuthors() ([]repo.Author, error) {
	return r.ListAuthorsFunc()
}

func (r *MockService) GetArticleById(id string) (repo.Article, error) {
	return r.GetArticleByIdFunc(id)
}

func (r *MockService) GetAuthorById(id string) (repo.Author, error) {
	return r.GetAuthorByIdFunc(id)
}

func (r *MockService) GetAuthorsByIds(id []string) ([]repo.Author, error) {
	return r.GetAuthorsByIdsFunc(id)
}

func (r *MockService) GetAuthorByNameAndEmail(name string, email string) (repo.Author, error) {
	return r.GetAuthorByNameAndEmailFunc(name, email)
}

func (r *MockService) AddAuthor(a repo.Author) (string, error) {
	r.Authors = append(r.Authors, a)
	return r.AddAuthorFunc(a)
}

func (r *MockService) AddArticle(a repo.Article) (string, error) {
	r.Articles = append(r.Articles, a)
	return r.AddArticleFunc(a)
}

func (r *MockService) DeleteAuthorById(id string) error {
	return r.DeleteAuthorByIdFunc(id)
}

func (r *MockService) DeleteArticleById(id string) error {
	return r.DeleteArticleByIdFunc(id)
}

func (r *MockService) DeleteAuthorByNameAndEmail(name string, email string) error {
	return r.DeleteAuthorByNameAndEmailFunc(name, email)
}
