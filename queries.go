package blog

const (
	listArticlesQuery = `SELECT art.id, art.title, art.body, art.posted_at, auth.name, auth.email
	FROM blog.articles art
	LEFT JOIN blog.authors auth on art.author_id = auth.id;`

	getArticleByIdQuery = `SELECT art.id, art.title, art.body, art.posted_at, auth.name, auth.email
	FROM blog.articles art 
	LEFT JOIN blog.authors auth on art.author_id = auth.id
	WHERE art.id = $1;`

	addArticleQuery = `INSERT INTO blog.articles(title, body, author_id) values ($1, $2, $3) RETURNING id;`
)
