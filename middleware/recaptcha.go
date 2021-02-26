package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"blog/common"
	"blog/conf"
)

var requestClient *http.Client

var (
	errorMap  map[string]string
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
	Recaptcha = NewRecaptchaMiddleware(config.Server.RecaptchaSecretKey, "https://recaptcha.net")
	requestClient = &http.Client{
		Timeout: time.Second * 5,
	}
}

type RecaptchaMiddleware struct {
	secretKey string
	baseUrl   string
}

type captchaResponseBody struct {
	Success            bool      `json:"success"`
	ChallengeTimestamp time.Time `json:"challenge_ts"` // timestamp of the challenge load (ISO format yyyy-MM-dd'T'HH:mm:ssZZ)
	Hostname           string    `json:"hostname"`     // the hostname of the site where the reCAPTCHA was solved
	ErrorCodes         []string  `json:"error-codes"`  // optional
}

func NewRecaptchaMiddleware(secret string, baseUrl string) *RecaptchaMiddleware {
	return &RecaptchaMiddleware{
		secretKey: secret,
		baseUrl:   baseUrl,
	}
}

func (r *RecaptchaMiddleware) sendApiRequest(c *gin.Context, token string) (*captchaResponseBody, error) {
	requestUrl := fmt.Sprintf("%s/recaptcha/api/siteverify?secret=%s&response=%s&remoteip=%s", r.baseUrl, r.secretKey, token, c.ClientIP())
	resp, err := requestClient.Post(requestUrl, "application/json", nil)
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

func (r *RecaptchaMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := r.sendApiRequest(c, c.Request.Header.Get("captchaToken"))
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
