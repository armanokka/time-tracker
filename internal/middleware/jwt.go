package middleware

import (
	"fmt"
	"github.com/armanokka/test_task_Effective_mobile/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"time"
)

func (m Manager) validateJWTToken(c *gin.Context, token string) error {
	ctx, span := m.tracer.Start(c, "Manager.validateJWTToken")
	defer span.End()
	c.Set(utils.UserCtxKey, ctx)

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return []byte(m.cfg.JWTSecretKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("unable to extract claims")
	}

	expiresAt := time.Unix(int64(claims["expires_at"].(float64)), 0)
	if expiresAt.Sub(time.Now()) < 0 {
		return fmt.Errorf("token expired")
	}

	userID := int64(claims["id"].(float64))
	user, err := m.authUC.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	c.Set("user", user)
	return nil
}
