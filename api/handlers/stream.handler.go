package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"leichil-backend/api/models"
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func app_stream_handler(router fiber.Router) {
	router.Post("/", get_signed_url)
}

type StreamURL struct {
	URL *string `json:"url,omitempty"`
}

var url_required_fields = []string{"URL"}

func get_signed_url(ctx *fiber.Ctx) error {
	var body = new(StreamURL)
	flag, status, message := utils.ParseAndValidate(body, url_required_fields, ctx)
	if flag {
		var key = os.Getenv("SIGNED_URL_KEY")
		var video_url = models.SanitizeData(body.URL)
		u, err := url.Parse(video_url)
		if err == nil {
			var allowed_path = path.Base(path.Dir(u.Path))
			allowed_path = "/" + allowed_path
			signed_url, expires, err := sign_url(video_url, key, allowed_path, 5*60*60)
			if err != nil {
				log.Error("Unable to get signed URL: ", err.Error())
				flag = false
			} else {
				var response = map[string]any{
					"success":    true,
					"url":        signed_url,
					"expires_on": expires,
				}
				return ctx.Status(fiber.StatusOK).JSON(response)
			}
		} else {
			log.Error("Unable to parse video URL: ", err.Error())
			flag = false
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func sign_url(video_url, key, allowed_path string, expiration int64) (string, int64, error) {
	parsed, err := url.Parse(video_url)
	if err != nil {
		log.Error("Error in parsing video url: ", err.Error())
		return "", 0, err
	}
	var expires int64 = time.Now().Unix() + expiration
	var params url.Values = url.Values{}
	if allowed_path != "" {
		params.Set("token_path", allowed_path)
	}
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signing_parts []string
	for _, k := range keys {
		signing_parts = append(signing_parts, k+"="+params.Get(k))
	}
	var signing_data string = strings.Join(signing_parts, "&")
	var signed_path = parsed.Path
	if allowed_path != "" {
		signed_path = allowed_path
	}
	var message = fmt.Sprintf("%s%d%s", signed_path, expires, signing_data)
	var hmac_data = hmac.New(sha256.New, []byte(key))
	hmac_data.Write([]byte(message))
	var digest = hmac_data.Sum(nil)
	var token = base64.RawURLEncoding.EncodeToString(digest)
	token = "HS256-" + token
	var base = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	var query = fmt.Sprintf("bcdn_token=%s&expires=%d", token, expires)
	if allowed_path != "" {
		query += "&token_path=" + url.QueryEscape(allowed_path)
	}
	return fmt.Sprintf("%s/%s%s", base, query, parsed.Path), expires, nil
}
