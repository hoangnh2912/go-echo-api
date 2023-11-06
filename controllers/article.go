package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"go-echo-api/middleware"
	"go-echo-api/model"
	"go-echo-api/requests"
	"go-echo-api/responses"
	"go-echo-api/services"
	"go-echo-api/utils"

	"github.com/labstack/echo/v4"
)

type ArticleController struct {
	articleService *services.ArticleService
	userService    *services.UserService
}

func NewArticleController(sg *echo.Group, articleService *services.ArticleService, userService *services.UserService) *ArticleController {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	ac := &ArticleController{
		articleService,
		userService,
	}
	articles := sg.Group("/articles", jwtMiddleware, middleware.JWTWithConfig(
		middleware.JWTConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().Method == "GET" && c.Path() != "/api/articles/feed" {
					return true
				}
				return false
			},
			SigningKey: utils.JWTSecret,
		},
	))
	articles.POST("", ac.CreateArticle)
	articles.GET("/feed", ac.Feed)
	articles.PUT("/:slug", ac.UpdateArticle)
	articles.DELETE("/:slug", ac.DeleteArticle)
	articles.POST("/:slug/comments", ac.AddComment)
	articles.DELETE("/:slug/comments/:id", ac.DeleteComment)
	articles.POST("/:slug/favorite", ac.Favorite)
	articles.DELETE("/:slug/favorite", ac.Unfavorite)
	articles.GET("", ac.Articles)
	articles.GET("/:slug", ac.GetArticle)
	articles.GET("/:slug/comments", ac.GetComments)

	tags := sg.Group("/tags")
	tags.GET("", ac.Tags)
	return ac
}

// GetArticle godoc
// @Summary Get an article
// @Description Get an article. Auth not required
// @ID get-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article to get"
// @Success 200 {object} responses.SingleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug} [get]
func (ac *ArticleController) GetArticle(c echo.Context) error {
	slug := c.Param("slug")
	a, err := ac.articleService.GetBySlug(slug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	return c.JSON(http.StatusOK, "test")
}

// Articles godoc
// @Summary Get recent articles globally
// @Description Get most recent articles globally. Use query parameters to filter results. Auth is optional
// @ID get-articles
// @Tags article
// @Accept  json
// @Produce  json
// @Param tag query string false "Filter by tag"
// @Param author query string false "Filter by author (username)"
// @Param favorited query string false "Filter by favorites of a user (username)"
// @Param limit query integer false "Limit number of articles returned (default is 20)"
// @Param offset query integer false "Offset/skip number of articles (default is 0)"
// @Success 200 {object} responses.ArticleListResponse
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles [get]
func (ac *ArticleController) Articles(c echo.Context) error {
	var (
		articles []model.Article
		count    int
	)

	tag := c.QueryParam("tag")
	author := c.QueryParam("author")
	favoritedBy := c.QueryParam("favorited")

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}

	if tag != "" {
		articles, count, err = ac.articleService.ListByTag(tag, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if author != "" {
		articles, count, err = ac.articleService.ListByAuthor(author, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if favoritedBy != "" {
		articles, count, err = ac.articleService.ListByWhoFavorited(favoritedBy, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		articles, count, err = ac.articleService.List(offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}

	return c.JSON(http.StatusOK, responses.NewArticleListResponse(ac.userService, utils.UserIDFromToken(c), articles, count))
}

// Feed godoc
// @Summary Get recent articles from users you follow
// @Description Get most recent articles from users you follow. Use query parameters to limit. Auth is required
// @ID feed
// @Tags article
// @Accept  json
// @Produce  json
// @Param limit query integer false "Limit number of articles returned (default is 20)"
// @Param offset query integer false "Offset/skip number of articles (default is 0)"
// @Success 200 {object} responses.ArticleListResponse
// @Failure 401 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/feed [get]
func (ac *ArticleController) Feed(c echo.Context) error {
	var (
		articles []model.Article
		count    int
	)

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}

	articles, count, err = ac.articleService.ListFeed(utils.UserIDFromToken(c), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, responses.NewArticleListResponse(ac.userService, utils.UserIDFromToken(c), articles, count))
}

// CreateArticle godoc
// @Summary Create an article
// @Description Create an article. Auth is require
// @ID create-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param article body requests.ArticleCreateRequest true "Article to create"
// @Success 201 {object} responses.SingleArticleResponse
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles [post]
func (ac *ArticleController) CreateArticle(c echo.Context) error {
	var a model.Article

	req := &requests.ArticleCreateRequest{}
	if err := req.Bind(c, &a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	a.AuthorID = utils.UserIDFromToken(c)

	err := ac.articleService.CreateArticle(&a)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, responses.NewArticleResponse(c, &a))
}

// UpdateArticle godoc
// @Summary Update an article
// @Description Update an article. Auth is required
// @ID update-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article to update"
// @Param article body requests.ArticleCreateRequest true "Article to update"
// @Success 200 {object} responses.SingleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug} [put]
func (ac *ArticleController) UpdateArticle(c echo.Context) error {
	slug := c.Param("slug")

	a, err := ac.articleService.GetUserArticleBySlug(utils.UserIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	req := &requests.ArticleUpdateRequest{}
	req.Populate(a)

	if err := req.Bind(c, a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = ac.articleService.UpdateArticle(a, req.Article.Tags); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, responses.NewArticleResponse(c, a))
}

// DeleteArticle godoc
// @Summary Delete an article
// @Description Delete an article. Auth is required
// @ID delete-article
// @Tags article
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug} [delete]
func (ac *ArticleController) DeleteArticle(c echo.Context) error {
	slug := c.Param("slug")

	a, err := ac.articleService.GetUserArticleBySlug(utils.UserIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	err = ac.articleService.DeleteArticle(a)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

// AddComment godoc
// @Summary Create a comment for an article
// @Description Create a comment for an article. Auth is required
// @ID add-comment
// @Tags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to create a comment for"
// @Param comment body requests.CreateCommentRequest true "Comment you want to create"
// @Success 201 {object} responses.SingleCommentResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/comments [post]
func (ac *ArticleController) AddComment(c echo.Context) error {
	slug := c.Param("slug")

	a, err := ac.articleService.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	var cm model.Comment

	req := &requests.CreateCommentRequest{}
	if err := req.Bind(c, &cm); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	if err = ac.articleService.AddComment(a, &cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, responses.NewCommentResponse(c, &cm))
}

// GetComments godoc
// @Summary Get the comments for an article
// @Description Get the comments for an article. Auth is optional
// @ID get-comments
// @Tags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to get comments for"
// @Success 200 {object} responses.CommentListResponse
// @Failure 422 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/comments [get]
func (ac *ArticleController) GetComments(c echo.Context) error {
	slug := c.Param("slug")

	cm, err := ac.articleService.GetCommentsBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, responses.NewCommentListResponse(c, cm))
}

// DeleteComment godoc
// @Summary Delete a comment for an article
// @Description Delete a comment for an article. Auth is required
// @ID delete-comments
// @Tags comment
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to delete a comment for"
// @Param id path integer true "ID of the comment you want to delete"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/comments/{id} [delete]
func (ac *ArticleController) DeleteComment(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	cm, err := ac.articleService.GetCommentByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if cm == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if cm.UserID != utils.UserIDFromToken(c) {
		return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}

	if err := ac.articleService.DeleteComment(cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

// Favorite godoc
// @Summary Favorite an article
// @Description Favorite an article. Auth is required
// @ID favorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to favorite"
// @Success 200 {object} responses.SingleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [post]
func (ac *ArticleController) Favorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := ac.articleService.GetBySlug(slug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if err := ac.articleService.AddFavorite(a, utils.UserIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, responses.NewArticleResponse(c, a))
}

// Unfavorite godoc
// @Summary Unfavorite an article
// @Description Unfavorite an article. Auth is required
// @ID unfavorite
// @Tags favorite
// @Accept  json
// @Produce  json
// @Param slug path string true "Slug of the article that you want to unfavorite"
// @Success 200 {object} responses.SingleArticleResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /articles/{slug}/favorite [delete]
func (ac *ArticleController) Unfavorite(c echo.Context) error {
	slug := c.Param("slug")

	a, err := ac.articleService.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if err := ac.articleService.RemoveFavorite(a, utils.UserIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, responses.NewArticleResponse(c, a))
}

// Tags godoc
// @Summary Get tags
// @Description Get tags
// @ID tags
// @Tags tag
// @Accept  json
// @Produce  json
// @Success 201 {object} responses.TagListResponse
// @Failure 400 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tags [get]
func (ac *ArticleController) Tags(c echo.Context) error {
	tags, err := ac.articleService.ListTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, responses.NewTagListResponse(tags))
}
