package http

import (
	"context"
	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository"
	"effective-mobile-task/pkg/logger"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "effective-mobile-task/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service interface {
	AddSong(ctx context.Context, song *models.Song) (uint64, error)
	UpdateSong(ctx context.Context, song *models.Song) error
	DeleteSong(ctx context.Context, id uint64) error
	GetSongLyrics(ctx context.Context, id uint64, offset uint64) (string, error)
	GetSongs(ctx context.Context, creds repository.Creds, offset uint64, limit uint64) ([]*models.Song, error)
}

type Client interface {
	GetSongInfo(group, song string) (*models.SongInfo, error)
}

type Handler struct {
	client  Client
	service Service
	logger  logger.Logger
}

func NewHandler(router *gin.Engine, client Client, service Service, logger logger.Logger) {
	handler := &Handler{
		client:  client,
		service: service,
		logger:  logger,
	}

	router.GET("/library", handler.GetSongs)
	router.GET("/song/:id/lyrics", handler.GetSongLyrics)
	router.POST("/song", handler.AddSong)
	router.PATCH("/song/:id", handler.UpdateSong)
	router.DELETE("/song/:id", handler.DeleteSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// AddSong godoc
// @Summary AddSong
// @Description Add new song to the library
// @Accept json
// @Produce json
// @Param input body models.AddSongInput true "Song data"
// @Success 200
// @Failure 400 {string} map[string]string
// @Failure 404 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /song [post]
func (h *Handler) AddSong(c *gin.Context) {
	var inp models.AddSongInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to read song from JSON", zap.String("err", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read song from JSON",
		})

		return
	}

	song := models.Song{
		Group: inp.Group,
		Song:  inp.Song,
	}

	//info, err := h.client.GetSongInfo(song.Group, song.Song)
	//if err != nil {
	//	h.logger.Error(c.Request.Context(), "Failed to get song info", zap.String("err", err.Error()))
	//
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": err.Error(),
	//	})
	//
	//	return
	//}

	//song.ReleaseDate = info.ReleaseDate
	//song.Text = info.Text
	//song.Link = info.Link

	parsedTime, err := time.Parse("02.01.2006", "16.07.2006")
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to parse time format", zap.String("expected", "02.01.2006"))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	song.ReleaseDate = parsedTime
	song.Text = "text1\n\ntext2\n\ntext3"
	song.Link = "link"

	id, err := h.service.AddSong(c.Request.Context(), &song)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to add song", zap.String("err", err.Error()))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add song",
		})

		return
	}

	h.logger.Info(c.Request.Context(), "Song added", zap.Uint64("id", id))

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

// GetSongLyrics godoc
// @Summary GetSongLyrics
// @Description List song lyrics by verses with paginating
// @Produce json
// @Param id path uint64 true "Song id"
// @Param page query uint64 true "Page number"
// @Success 200
// @Failure 400 {string} map[string]string
// @Failure 404 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /song/{id}/lyrics [get]
func (h *Handler) GetSongLyrics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to parse song id", zap.String("err", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Song ID",
		})

		return
	}

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page == 0 {
		h.logger.Error(c.Request.Context(), "Failed to get song lyrics")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid page value",
		})

		return
	}

	offset := page - 1

	verses, err := h.service.GetSongLyrics(c.Request.Context(), id, offset)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to get song lyrics", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verses": verses,
	})
}

// UpdateSong godoc
// @Summary UpdateSong
// @Description Update song in the library by ID
// @Accept json
// @Produce json
// @Param id path uint64 true "Song id"
// @Param input body models.UpdateSongInput true "Song data"
// @Success 200
// @Failure 400 {string} map[string]string
// @Failure 404 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /song/{id} [patch]
func (h *Handler) UpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to parse song id", zap.String("err", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Song ID",
		})

		return
	}

	var inp models.UpdateSongInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to read song from JSON", zap.String("err", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read song from JSON",
		})

		return
	}

	parsedTime, err := time.Parse("02.01.2006", inp.ReleaseDate)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to parse time format", zap.String("expected", "02.01.2006"))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	song := models.Song{
		Id:          id,
		Group:       inp.Group,
		Song:        inp.Song,
		ReleaseDate: parsedTime,
		Text:        inp.Text,
		Link:        inp.Link,
	}

	if err := h.service.UpdateSong(c.Request.Context(), &song); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to update song", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	h.logger.Info(c.Request.Context(), "Song updated", zap.Uint64("id", id))

	c.Status(http.StatusOK)
}

// DeleteSong godoc
// @Summary DeleteSong
// @Description Delete song from the library by ID
// @Produce json
// @Param id path uint64 true "Song id"
// @Success 200
// @Failure 400 {string} map[string]string
// @Failure 404 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /song/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to parse song id", zap.String("err", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Song ID",
		})

		return
	}

	if err := h.service.DeleteSong(c.Request.Context(), id); err != nil {
		h.logger.Error(c.Request.Context(), "Failed to delete song", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete song",
		})

		return
	}

	h.logger.Info(c.Request.Context(), "Song deleted", zap.Uint64("id", id))

	c.Status(http.StatusOK)
}

// GetSongs godoc
// @Summary GetSongs
// @Description List all songs in the library. Songs could be filtered
// @Produce json
// @Param group query string false "Song group"
// @Param song query string false "Song name"
// @Param release_date query string false "Song release date"
// @Param page query uint64 true "Page number"
// @Param limit query uint64 true "Limit of songs"
// @Success 201
// @Failure 404 {string} map[string]string
// @Failure 500 {string} map[string]string
// @Router /library [get]
func (h *Handler) GetSongs(c *gin.Context) {
	creds := repository.Creds{}
	if group := c.Query("group"); group != "" {
		creds["group"] = group
	}
	if song := c.Query("song"); song != "" {
		creds["song"] = song
	}
	if date := c.Query("release_date"); date != "" {
		creds["release_data"] = date
	}

	page, err := strconv.ParseUint(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page == 0 {
		h.logger.Error(c.Request.Context(), "Failed to get song lyrics")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid page value",
		})

		return
	}

	limit, err := strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit == 0 {
		h.logger.Error(c.Request.Context(), "Failed to get song lyrics")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid limit value",
		})

		return
	}

	offset := (page - 1) * limit

	songs, err := h.service.GetSongs(c.Request.Context(), creds, offset, limit)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Failed to list songs", zap.String("err", err.Error()))

		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list songs",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
	})
}
