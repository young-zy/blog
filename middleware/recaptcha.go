package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/conf"
)

var requestClient *http.Client

var (
	errorMap map[string]string
	// Recaptcha the recaptcha middleware for this project
	Recaptcha *RecaptchaMiddleware
)

func init() {
	errorMap = map[string]string{
		"missing-input-secret":   "The secret parameter is missing.",
		"invalid-input-secret":   "The secret parameter is invalid or malformed.",
		"missing-input-response": "The captchaToken parameter is missing.",
		"invalid-input-response": "The captchaToken parameter is invalid or malformed.",
		"bad-request":            "The request is invalid or malformed.",
		"timeout-or-duplicate":   "The response is no longer valid: either is too old or has been used previously.",
	}
	config := conf.Config
	Recaptcha = newRecaptchaMiddleware(config.Server.RecaptchaSecretKey, "https://recaptcha.net")
	requestClient = &http.Client{
		Timeout: time.Second * 5,
	}
}

// RecaptchaMiddleware is a recaptcha authentication middleware
type RecaptchaMiddleware struct {
	secretKey string
	baseURL   string
}

type captchaResponseBody struct {
	Success            bool      `json:"success"`
	ChallengeTimestamp time.Time `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname           string    `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes         []string  `json:"error-codes"`  // optional
}

// newRecaptchaMiddleware creates a new recaptcha middleware
func newRecaptchaMiddleware(secret string, baseURL string) *RecaptchaMiddleware {
	return &RecaptchaMiddleware{
		secretKey: secret,
		baseURL:   baseURL,
	}
}

func (r *RecaptchaMiddleware) sendAPIRequest(c *gin.Context, token string) (*captchaResponseBody, error) {
	requestURL := fmt.Sprintf("%s/recaptcha/api/siteverify?secret=%s&response=%s&remoteip=%s", r.baseURL, r.secretKey, token, c.ClientIP())
	resp, err := requestClient.Post(requestURL, "application/json", nil)
	if err != nil {
		return nil, err
	}
	respStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respBody := &captchaResponseBody{}
	err = json.Unmarshal(respStr, respBody)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// Middleware returns a gin handler function that handles the recaptcha token
func (r *RecaptchaMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := r.sendAPIRequest(c, c.Request.Header.Get("captchaToken"))
		if err != nil {
			common.NewInternalError(c, err)
		}
		if !resp.Success {
			errorMessageBuilder := &strings.Builder{}
			for _, errorCode := range resp.ErrorCodes {
				errorMessageBuilder.WriteString(errorMap[errorCode])
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": errorMessageBuilder.String(),
			})
		} else {
			c.Next()
		}
	}
}
