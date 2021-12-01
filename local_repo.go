package blog

import "errors"

type LocalRepository struct {
	DB []Article
}

// List all articles in the database.
func (repo *LocalRepository) ListArticles() ([]Article, error) {

	articles := make([]Article, 0, len(repo.DB))
	articles = append(articles, repo.DB...)

	return articles, nil
}

// Get article by id.
func (repo *LocalRepository) GetArticleById(id int) (Article, error) {
	for _, article := range repo.DB {
		if article.Id == id {
			return article, nil
		}
	}
	return *new(Article), errors.New("Article not found")
}

// Add new article to the database.
func (repo *LocalRepository) PostArticle(a Article) (Article, error) {
	repo.DB = append(repo.DB, a)
	return a, nil
}

// Delete article by id.
func (repo *LocalRepository) DeleteArticleById(id int) (Article, error) {
	for index, article := range repo.DB {
		if article.Id == id {
			repo.DB = append(repo.DB[:index], repo.DB[index+1:]...)
			return article, nil
		}
	}
	return *new(Article), errors.New("Article not found")
}
