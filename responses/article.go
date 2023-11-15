package responses

import (
	"go-echo-api/model"
	"go-echo-api/services"
	"go-echo-api/utils"
	"time"

	"github.com/labstack/echo/v4"
)

type ArticleResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type SingleArticleResponse struct {
	Article *ArticleResponse `json:"article"`
}

type ArticleListResponse struct {
	Articles      []*ArticleResponse `json:"articles"`
	ArticlesCount int64              `json:"articlesCount"`
}

func NewArticleResponse(c echo.Context, a *model.Article) *SingleArticleResponse {
	ar := new(ArticleResponse)
	ar.TagList = make([]string, 0)
	ar.Slug = a.Slug
	ar.Title = a.Title
	ar.Description = a.Description
	ar.Body = a.Body
	ar.CreatedAt = a.CreatedAt
	ar.UpdatedAt = a.UpdatedAt
	for _, t := range a.Tags {
		ar.TagList = append(ar.TagList, t.Tag)
	}
	for _, u := range a.Favorites {
		if u.ID == utils.UserIDFromToken(c) {
			ar.Favorited = true
		}
	}
	ar.FavoritesCount = len(a.Favorites)
	ar.Author.Username = a.Author.Username
	ar.Author.Image = a.Author.Image
	ar.Author.Bio = a.Author.Bio
	ar.Author.Following = a.Author.FollowedBy(utils.UserIDFromToken(c))
	return &SingleArticleResponse{ar}
}

func NewArticleListResponse(us *services.UserService, userID uint, articles []model.Article, count int64) *ArticleListResponse {
	r := new(ArticleListResponse)
	r.Articles = make([]*ArticleResponse, 0)
	for _, a := range articles {
		ar := new(ArticleResponse)
		ar.TagList = make([]string, 0)
		ar.Slug = a.Slug
		ar.Title = a.Title
		ar.Description = a.Description
		ar.Body = a.Body
		ar.CreatedAt = a.CreatedAt
		ar.UpdatedAt = a.UpdatedAt
		for _, t := range a.Tags {
			ar.TagList = append(ar.TagList, t.Tag)
		}
		for _, u := range a.Favorites {
			if u.ID == userID {
				ar.Favorited = true
			}
		}
		ar.FavoritesCount = len(a.Favorites)
		ar.Author.Username = a.Author.Username
		ar.Author.Image = a.Author.Image
		ar.Author.Bio = a.Author.Bio
		ar.Author.Following, _ = us.IsFollower(a.AuthorID, userID)
		r.Articles = append(r.Articles, ar)
	}
	r.ArticlesCount = count
	return r
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type SingleCommentResponse struct {
	Comment *CommentResponse `json:"comment"`
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
}

func NewCommentResponse(c echo.Context, cm *model.Comment) *SingleCommentResponse {
	comment := new(CommentResponse)
	comment.ID = cm.ID
	comment.Body = cm.Body
	comment.CreatedAt = cm.CreatedAt
	comment.UpdatedAt = cm.UpdatedAt
	comment.Author.Username = cm.User.Username
	comment.Author.Image = cm.User.Image
	comment.Author.Bio = cm.User.Bio
	comment.Author.Following = cm.User.FollowedBy(utils.UserIDFromToken(c))
	return &SingleCommentResponse{comment}
}

func NewCommentListResponse(c echo.Context, comments []model.Comment) *CommentListResponse {
	r := new(CommentListResponse)
	cr := CommentResponse{}
	r.Comments = make([]CommentResponse, 0)
	for _, i := range comments {
		cr.ID = i.ID
		cr.Body = i.Body
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.Author.Username = i.User.Username
		cr.Author.Image = i.User.Image
		cr.Author.Bio = i.User.Bio
		cr.Author.Following = i.User.FollowedBy(utils.UserIDFromToken(c))

		r.Comments = append(r.Comments, cr)
	}
	return r
}

type TagListResponse struct {
	Tags []string `json:"tags"`
}

func NewTagListResponse(tags []model.Tag) *TagListResponse {
	r := new(TagListResponse)
	for _, t := range tags {
		r.Tags = append(r.Tags, t.Tag)
	}
	return r
}
