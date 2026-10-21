package handlers

import (
	"fmt"
	"leichil-backend/api/models"
	"leichil-backend/api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var content_required_fields = []string{
	"Name", "Genre", "Runtime", "ReleaseYear", "Synopsis", "Price", "IsPublished",
	"ProducerID", "Thumbnail", "Video", "Category", "ContentType", "ActorList",
}

var episode_required_fields = []string{
	"Name", "Season", "EpisodeNumber", "Runtime", "ReleaseYear", "Synopsis",
	"Thumbnail", "Video",
}

var content_rating_required_fields = []string{"Rating", "Review"}

func content_handler(router fiber.Router) {
	router.Get("/", get_contents)
	router.Get("/:id", get_content_by_id)
	router.Post("/", add_content)
	router.Put("/:id", update_content)
	router.Delete("/:id", delete_content)
	router.Get("/:id/episodes", get_episodes)
	router.Get("/:id/episodes/:episode_id", get_episode_by_id)
	router.Post("/:id/episodes", add_episode)
	router.Put("/:id/episodes/:episode_id", update_episode)
	router.Delete("/:id/episodes/:episode_id", delete_episode)
}

func app_content_handler(router fiber.Router) {
	router.Get("/", get_contents)
	router.Get("/:id", get_content_by_id)
	router.Get("/home/list", get_home_contents)
	router.Get("/continue/list", get_continue_watching_list)
	router.Get("/all/list", get_all_list)
	router.Get("/:id/episodes", get_episodes)
	router.Get("/:id/episodes/:episode_id", get_episode_by_id)
	router.Get("/:id/ratings/list", get_content_ratings)
	router.Post("/:id/ratings", add_content_rating)
	router.Put("/:id/ratings", update_content_rating)
	router.Delete("/:id/ratings/", delete_content_rating)
}

/*----------------------CONTENT API ROUTES-----------------*/
func get_contents(ctx *fiber.Ctx) error {
	var content models.Content
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var published int = utils.GetIntQueryParams(ctx.Query("published"), 1)
	var params = map[string]string{
		"name":      ctx.Query("name"),
		"genre":     ctx.Query("genre"),
		"type":      ctx.Query("type"),
		"published": fmt.Sprintf("%d", published),
	}
	var key = fmt.Sprintf("content-%s", ctx.OriginalURL())
	cache_data, err := utils.GetData(key)
	if err != nil {
		_, result, err := content.Get(offset, limit, params)
		if len(result) > 0 {
			utils.SetListData[models.Content](key, result, 0)
		}
		return utils.GetArrayResponse(ctx, 0, result, "Unable to get Contents", err)
	}
	return utils.GetCachedArrayResponse(ctx, 0, cache_data)
}

func get_home_contents(ctx *fiber.Ctx) error {
	var content models.HomeData
	var err error = nil
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var params = strings.ToUpper(strings.TrimSpace(ctx.Query("content-type")))
	if len(params) == 0 {
		content, err = content.GetHomeList(token)
	} else {
		content, err = content.GetCustomList(token, params)
	}
	if err != nil {
		log.Error("Error in getting App Home Contents: ", err.Error())
		return utils.GetMessageResponse("", "Unable to get app home contents", fiber.StatusInternalServerError, false, ctx)
	}
	return utils.GetSingleResponse(ctx, []any{content}, "Unable to get Home Data", err)
}

func get_continue_watching_list(ctx *fiber.Ctx) error {
	var cca models.CustomerContentActivity
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	result, err := cca.GetContinueWatching(token)
	if err != nil {
		log.Error("Error in getting Continue Watching List: ", err.Error())
		return utils.GetMessageResponse("", "Unable to get continue watching list", fiber.StatusInternalServerError, false, ctx)
	}
	var response = map[string]any{
		"success": true,
		"result":  result,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func get_all_list(ctx *fiber.Ctx) error {
	var content models.Content
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"content_type": ctx.Query("content_type"),
	}
	var list_type = strings.ToLower(ctx.Query("type"))
	switch list_type {
	case "trending":
		_, result, err := content.GetTrendingList(offset, limit, params)
		return utils.GetArrayResponse(ctx, 0, result, "Unable to get Trending List", err)
	case "new":
		_, result, err := content.GetNewList(offset, limit, params)
		return utils.GetArrayResponse(ctx, 0, result, "Unable to get New List", err)
	default:
		return utils.GetMessageResponse("", "Invalid Request", fiber.StatusBadRequest, false, ctx)
	}
}

func get_content_by_id(ctx *fiber.Ctx) error {
	var id = ctx.Params("id")
	var key = fmt.Sprintf("content-%s", id)
	cache_data, err := utils.GetData(key)
	if err != nil {
		var published = utils.GetIntQueryParams(ctx.Query("published"), 1)
		var content models.Content
		result, err := content.GetById(id, published)
		if err != nil {
			log.Error(err.Error())
			return utils.GetMessageResponse("", "Unable to get Content Data", fiber.StatusInternalServerError, false, ctx)
		}
		if result.ID != nil {
			utils.SetSingleData[models.Content](key, []any{result}, 0)
		}
		return utils.GetSingleResponse(ctx, []any{result}, "Unable to get Content Data", err)
	}
	return utils.GetCachedSingleResponse(ctx, cache_data)
}

func add_content(ctx *fiber.Ctx) error {
	var content = new(models.Content)
	flag, status, message := utils.ParseAndValidate(content, content_required_fields, ctx)
	if flag {
		var token = ctx.Locals("user_token").(models.UserToken)
		if content.Add(token) {
			status = fiber.StatusOK
			message = "Content Added"
			utils.DeleteKeyPattern("content*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Content"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_content(ctx *fiber.Ctx) error {
	var content = new(models.Content)
	flag, status, message := utils.ParseAndValidate(content, content_required_fields, ctx)
	if flag {
		if content.Update(ctx.Params("id")) {
			status = fiber.StatusOK
			message = "Content Updated"
			utils.DeleteKeyPattern("content*")
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Content"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_content(ctx *fiber.Ctx) error {
	var content models.Content
	var flag bool = content.Delete(ctx.Params("id"))
	if flag {
		utils.DeleteKeyPattern("content*")
	}
	return utils.GetMessageResponse("Content Deleted", "Unable to Delete Content", fiber.StatusOK, flag, ctx)
}

/*----------------------EPISODE API ROUTES-----------------*/
func get_episodes(ctx *fiber.Ctx) error {
	var episode models.Episode
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"name":   ctx.Query("name"),
		"season": ctx.Query("season"),
	}
	var content_id string = ctx.Params("id")
	count, result, err := episode.Get(offset, limit, content_id, params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Episodes", err)
}

func get_episode_by_id(ctx *fiber.Ctx) error {
	var episode models.Episode
	result, err := episode.GetById(ctx.Params("episode_id"), ctx.Params("id"))
	if err != nil {
		log.Error(err.Error())
		return utils.GetMessageResponse("", "Unable to get Episode Data", fiber.StatusInternalServerError, false, ctx)
	}
	var response = map[string]any{
		"success": true,
		"result":  result,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func add_episode(ctx *fiber.Ctx) error {
	var episode = new(models.Episode)
	flag, status, message := utils.ParseAndValidate(episode, episode_required_fields, ctx)
	if flag {
		var content_id = ctx.Params("id")
		var token = ctx.Locals("user_token").(models.UserToken)
		if episode.Add(content_id, token) {
			status = fiber.StatusOK
			message = "Content Added"
			utils.DeleteData(fmt.Sprintf("content-%s", content_id))
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Episode"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_episode(ctx *fiber.Ctx) error {
	var episode = new(models.Episode)
	flag, status, message := utils.ParseAndValidate(episode, episode_required_fields, ctx)
	if flag {
		if episode.Update(ctx.Params("episode_id")) {
			status = fiber.StatusOK
			message = "Episode Updated"
			utils.DeleteData(fmt.Sprintf("content-%s", ctx.Params("id")))
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Episode"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_episode(ctx *fiber.Ctx) error {
	var episode models.Episode
	var flag bool = episode.Delete(ctx.Params("episode_id"))
	if flag {
		utils.DeleteData(fmt.Sprintf("content-%s", ctx.Params("id")))
	}
	return utils.GetMessageResponse("Episode Deleted", "Unable to Delete Episode", fiber.StatusOK, flag, ctx)
}

/*----------------------CONTENT RATING API ROUTES-----------------*/
func get_content_ratings(ctx *fiber.Ctx) error {
	var cr models.ContentRating
	var offset int = utils.GetIntQueryParams(ctx.Query("offset"), 0)
	var limit int = utils.GetIntQueryParams(ctx.Query("limit"), 10)
	var params = map[string]string{
		"customer_id": ctx.Query("customer_id"),
	}
	count, result, err := cr.Get(offset, limit, ctx.Params("id"), params)
	var count_data int = models.SanitizeData(count.Count)
	return utils.GetArrayResponse(ctx, count_data, result, "Unable to get Content Ratings", err)
}

func add_content_rating(ctx *fiber.Ctx) error {
	var cr = new(models.ContentRating)
	flag, status, message := utils.ParseAndValidate(cr, content_rating_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		var content_id = ctx.Params("id")
		if cr.Add(content_id, token) {
			log.Info(fmt.Sprintf("content-%s", content_id))
			utils.DeleteData(fmt.Sprintf("content-%s", content_id))
			status = fiber.StatusOK
			message = "Content Rating Added"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Add Content Rating"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func update_content_rating(ctx *fiber.Ctx) error {
	var cr = new(models.ContentRating)
	flag, status, message := utils.ParseAndValidate(cr, content_rating_required_fields, ctx)
	if flag {
		var token = ctx.Locals("customer_token").(models.CustomerToken)
		var content_id = ctx.Params("id")
		if cr.Update(content_id, token) {
			utils.DeleteData(fmt.Sprintf("content-%s", content_id))
			status = fiber.StatusOK
			message = "Content Rating Updated"
		} else {
			status = fiber.StatusInternalServerError
			message = "Unable to Update Content Rating"
		}
	}
	return utils.GetMessageResponse(message, message, status, flag, ctx)
}

func delete_content_rating(ctx *fiber.Ctx) error {
	var cr models.ContentRating
	var token = ctx.Locals("customer_token").(models.CustomerToken)
	var flag bool = cr.Delete(ctx.Params("id"), token)
	return utils.GetMessageResponse("Content Rating Deleted", "Unable to Delete Content Rating", fiber.StatusOK, flag, ctx)
}
