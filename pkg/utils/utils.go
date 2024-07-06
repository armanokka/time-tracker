package utils

import (
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/armanokka/test_task_Effective_mobile/pkg/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const UserCtxKey = "ctx"
const defaultSearchQueryLimit = 1

type Response struct {
	Ok bool `json:"ok"`
}

var validate = validator.New()

func LogResponseError(c *gin.Context, log logger.Logger, err error) {
	log.Errorf("ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		requestid.Get(c), c.ClientIP(), err.Error())
}

// CreateSessionCookie configure jwt cookie
func CreateSessionCookie(cfg *config.CookieConfig, session string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Name,
		Value: session,
		Path:  "/",
		// Domain: "/",
		// Expires:    time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.MaxAge,
		Secure:     cfg.Secure,
		HttpOnly:   cfg.HttpOnly,
		SameSite:   0,
	}
}

func ReadRequest(c *gin.Context, request interface{}) error {
	if err := c.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(c, request)
}

type UsersQueryResponse struct {
	Users      []*models.User `json:"users"`
	Count      int            `json:"count"`
	Page       int            `json:"page"`
	TotalCount int            `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}

type UsersQuery struct {
	MinID             int    `form:"min_id"`
	MaxID             int    `form:"max_id"`
	Email             string `form:"email"`
	Name              string `form:"name"`
	Surname           string `form:"surname"`
	Patronymic        string `form:"patronymic"`
	Address           string `form:"address"`
	MaxPassportNumber int    `form:"max_passport_number"`
	MinPassportNumber int    `form:"min_passport_number"`
	MaxPassportSeries int    `form:"max_passport_series"`
	MinPassportSeries int    `form:"min_passport_series"`
	Limit             int    `form:"limit"`
	Page              int    `form:"page"`
}

func (u UsersQuery) GetLimit() int {
	if u.Limit == 0 {
		return defaultSearchQueryLimit
	}
	return u.Limit
}

func (u UsersQuery) GetOffset() int {
	if u.Page > 0 {
		u.Page--
	}
	return u.Page * u.Limit
}
