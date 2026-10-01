package models

import (
	"database/sql"
	"leichil-backend/api/utils"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

type Content struct {
	ID            *int       `json:"id"`
	Name          *string    `form:"name" json:"name,omitempty"`
	Genre         *string    `form:"genre" json:"genre,omitempty"`
	Runtime       *int       `form:"runtime" json:"runtime,omitempty"`
	ReleaseYear   *int       `form:"release_year" json:"release_year,omitempty"`
	Synopsis      *string    `form:"synopsis" json:"synopsis,omitempty"`
	Price         *int       `form:"price" json:"price,omitempty"`
	ProducerID    *int       `form:"producer_id" json:"producer_id,omitempty"`
	Thumbnail     *string    `form:"thumbnail" json:"thumbnail,omitempty"`
	Trailer       *string    `form:"trailer" json:"trailer,omitempty"`
	Video         *string    `form:"video" json:"video,omitempty"`
	Rating        *float32   `json:"rating,omitempty"`
	Category      *string    `json:"category,omitempty"`
	ContentType   *string    `json:"content_type,omitempty"`
	CreatedOn     *string    `json:"created_on,omitempty"`
	CreatedBy     *int       `json:"added_by,omitempty"`
	CreatedByName *string    `json:"added_by_name,omitempty"`
	ProducerName  *string    `json:"producer_name,omitempty"`
	Cast          *[]Cast    `json:"cast,omitempty"`
	ActorList     *[]int     `form:"actor_list" json:"actor_list,omitempty"`
	Seasons       *int       `json:"seasons,omitempty"`
	Episodes      *[]Episode `json:"episodes,omitempty"`
}

type HomeData struct {
	Carousel         *[]Content                 `json:"carousel"`
	ContinueWatching *[]CustomerContentActivity `json:"continue_watching"`
	Segments         *[]Segment                 `json:"segments"`
}

type Segment struct {
	Title       *string    `json:"title"`
	List        *[]Content `json:"list"`
	ContentType string     `json:"content_type"`
	ListType    string     `json:"list_type"`
}

type Cast struct {
	ContentID *int    `form:"content_id" json:"content_id,omitempty"`
	ActorID   *int    `form:"actor_id" json:"actor_id,omitempty"`
	ActorName *string `json:"actor_name"`
	Image     *string `json:"image"`
}

type Episode struct {
	ID            *int     `json:"id"`
	Season        *int     `form:"season" json:"season"`
	EpisodeNumber *int     `form:"episode_number" json:"episode_number"`
	Name          *string  `form:"name" json:"name"`
	ReleaseYear   *int     `form:"release_year" json:"release_year"`
	Synopsis      *string  `form:"synopsis" json:"synopsis"`
	Video         *string  `form:"video" json:"video"`
	Rating        *float32 `form:"rating" json:"rating"`
	Runtime       *int     `form:"runtime" json:"runtime"`
	Thumbnail     *string  `form:"thumbnail" json:"thumbnail"`
	CreatedOn     *string  `form:"created_on" json:"created_on"`
	CreatedBy     *int     `form:"created_by" json:"created_by"`
	ContentID     *int     `form:"content_id" json:"content_id"`
	CreatedByName *string  `form:"created_by_name" json:"created_by_name"`
	ContentName   *string  `form:"content_name" json:"content_name"`
}

type ContentRating struct {
	ContentID    *int    `json:"content_id,omitempty"`
	CustomerID   *int    `json:"customer_id,omitempty"`
	Rating       *int    `json:"rating,omitempty"`
	Review       *string `json:"review,omitempty"`
	CreatedOn    *string `json:"created_on,omitempty"`
	ContentName  *string `json:"content_name,omitempty"`
	CustomerName *string `json:"customer_name,omitempty"`
	Mobile       *string `json:"mobile,omitempty"`
}

type Rating struct {
	Rating *float32 `json:"rating"`
}

/*----------------------------- CONTENT ---------------------------------*/
func (obj Content) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_content_fields(&obj)
	var query string = get_contents_query
	var count_query string = count_contents_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query + content_order_filter
	count_query += filter_query
	return GetListData(filter, fields, &obj, count_query, query, offset, limit)
}

func (obj HomeData) GetHomeList(token CustomerToken) (HomeData, error) {
	var result HomeData
	var segments []Segment
	var content Content
	var cca CustomerContentActivity
	var continue_list []CustomerContentActivity
	var first_list []Content
	var second_list []Content
	var third_list []Content
	var first_title string
	var second_title string
	var third_title string
	var fields = []any{
		&content.ID, &content.Name, &content.Thumbnail,
		&content.ReleaseYear, &content.Genre, &content.Rating,
	}

	var first_section_query string = get_home_contents_query + content_order_filter
	_, first_result, err := GetListData([]any{}, fields, &content, "", first_section_query, 0, 5)
	if err != nil {
		log.Error("Unable to get First List: ", err.Error())
	} else {
		first_list = utils.FormatListResult[Content](first_result)
	}

	var second_section_query string = get_trending_contents_query + cca_order_filter
	_, second_result, err := GetListData([]any{}, fields, &content, "", second_section_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Second List: ", err.Error())
	} else {
		second_list = utils.FormatListResult[Content](second_result)
	}

	var third_section_query string = get_home_contents_query + content_type_filter + content_order_filter
	_, third_result, err := GetListData([]any{"SERIES"}, fields, &content, "", third_section_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Third List: ", err.Error())
	} else {
		third_list = utils.FormatListResult[Content](third_result)
	}

	var continue_query string = get_continue_watching_query + cca_customer_filter + cca_order_filter
	var continue_fields = get_continue_watching_fields(&cca)
	_, continue_result, err := GetListData([]any{token.CustomerID}, continue_fields, &cca, "", continue_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Continue List: ", err.Error())
	} else {
		continue_list = utils.FormatListResult[CustomerContentActivity](continue_result)
	}

	first_title = "New Releases"
	var first_segment = Segment{
		Title:       &first_title,
		List:        &first_list,
		ContentType: "",
		ListType:    "new",
	}
	segments = append(segments, first_segment)

	second_title = "Trending"
	var second_segment = Segment{
		Title:       &second_title,
		List:        &second_list,
		ContentType: "",
		ListType:    "trending",
	}
	segments = append(segments, second_segment)

	third_title = "Latest Shows"
	var third_segment = Segment{
		Title:       &third_title,
		List:        &third_list,
		ContentType: "SERIES",
		ListType:    "new",
	}
	segments = append(segments, third_segment)
	result.Segments = &segments
	result.Carousel = &first_list
	result.ContinueWatching = &continue_list
	return result, nil
}

func (obj Content) GetTrendingList(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = []any{
		&obj.ID, &obj.Name, &obj.Thumbnail,
		&obj.ReleaseYear, &obj.Genre, &obj.Rating,
	}
	var query string = get_trending_contents_query
	var filter = []any{}
	if len(params["content_type"]) > 0 {
		query += cca_content_type_filter
		filter = append(filter, params["content_type"])
	}
	query += cca_order_filter
	return GetListData(filter, fields, &obj, "", query, offset, limit)
}

func (obj Content) GetNewList(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = []any{
		&obj.ID, &obj.Name, &obj.Thumbnail,
		&obj.ReleaseYear, &obj.Genre, &obj.Rating,
	}
	var query string = get_home_contents_query
	var filter = []any{}
	if len(params["content_type"]) > 0 {
		query += content_type_filter
		filter = append(filter, params["content_type"])
	}
	query += content_order_filter
	return GetListData(filter, fields, &obj, "", query, offset, limit)
}

func (obj HomeData) GetCustomList(token CustomerToken, content_type string) (HomeData, error) {
	var result HomeData
	var segments []Segment
	var content Content
	var cca CustomerContentActivity
	var continue_list []CustomerContentActivity
	var new_list []Content
	var trending_list []Content
	var first_title string
	var second_title string
	var fields = []any{
		&content.ID, &content.Name, &content.Thumbnail,
		&content.ReleaseYear, &content.Genre, &content.Rating,
	}

	var new_releases_query string = get_home_contents_query + content_type_filter + content_order_filter
	_, new_releases_result, err := GetListData([]any{strings.ToUpper(content_type)}, fields, &content, "", new_releases_query, 0, 5)
	if err != nil {
		log.Error("Unable to get New Releases List: ", err.Error())
	} else {
		new_list = utils.FormatListResult[Content](new_releases_result)
	}

	var trending_query string = get_trending_contents_query + cca_content_type_filter + cca_order_filter
	_, trending_result, err := GetListData([]any{strings.ToUpper(content_type)}, fields, &content, "", trending_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Trending List: ", err.Error())
	} else {
		trending_list = utils.FormatListResult[Content](trending_result)
	}

	var continue_query string = get_continue_watching_query + cca_content_type_filter + cca_customer_filter + cca_order_filter
	var continue_fields = get_continue_watching_fields(&cca)
	_, continue_result, err := GetListData([]any{strings.ToUpper(content_type), token.CustomerID}, continue_fields, &cca, "", continue_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Continue List: ", err.Error())
	} else {
		continue_list = utils.FormatListResult[CustomerContentActivity](continue_result)
	}

	first_title = "New Releases"
	var first_segment = Segment{
		Title:       &first_title,
		List:        &new_list,
		ContentType: content_type,
		ListType:    "new",
	}
	segments = append(segments, first_segment)

	second_title = "Trending"
	var second_segment = Segment{
		Title:       &second_title,
		List:        &trending_list,
		ContentType: content_type,
		ListType:    "trending",
	}
	segments = append(segments, second_segment)

	result.Carousel = &new_list
	result.Segments = &segments
	result.ContinueWatching = &continue_list
	return result, nil
}

func (cca CustomerContentActivity) GetContinueWatching(token CustomerToken) ([]CustomerContentActivity, error) {
	var continue_list []CustomerContentActivity
	var continue_query string = get_continue_watching_query + cca_customer_filter + cca_order_filter
	var continue_fields = get_continue_watching_fields(&cca)
	_, continue_result, err := GetListData([]any{token.CustomerID}, continue_fields, &cca, "", continue_query, 0, 5)
	if err != nil {
		log.Error("Unable to get Continue List: ", err.Error())
		return continue_list, err
	}
	continue_list = utils.FormatListResult[CustomerContentActivity](continue_result)
	return continue_list, nil
}

func (obj Content) GetById(id string) (Content, error) {
	var episode Episode
	var episode_list []Episode
	var cast Cast
	var content_fields = get_content_fields(&obj)
	var cast_fields = get_cast_fields(&cast)
	var content_query string = get_contents_query + content_id_filter
	var content_filter = []any{id}
	result, err := GetSingleData(content_filter, content_fields, &obj, content_query)
	if err != nil {
		log.Error("Error in Getting Content By ID: ", err.Error())
		return obj, err
	}
	if len(result) == 0 {
		return obj, nil
	}
	obj = utils.FormatSingleResult[Content](result)
	if SanitizeData(obj.ContentType) == "SERIES" {
		var seasons int = 1
		var episode_fields = get_episode_fields(&episode)
		_, episode_result, err := GetListData(content_filter, episode_fields, &episode, "", get_episodes_query, 0, 1000)
		if err != nil {
			log.Error("Error in Getting Content Episodes: ", err.Error())
		} else {
			for _, item := range episode_result {
				var ep = item.(Episode)
				seasons = max(seasons, SanitizeData(ep.Season))
				episode_list = append(episode_list, ep)
			}
			obj.Seasons = &seasons
			obj.Episodes = &episode_list
		}
	}
	_, cast_result, err := GetListData(content_filter, cast_fields, &cast, "", get_casts_query, 0, 1000)
	if err != nil {
		log.Error("Error in Getting Content Cast: ", err.Error())
		return obj, err
	}
	var actors []Cast
	for _, elem := range cast_result {
		actors = append(actors, elem.(Cast))
	}
	obj.Cast = &actors
	return obj, nil
}

func (obj *Content) Add(token UserToken) bool {
	var DB *sql.DB = utils.GetDBInstance()
	var tx *sql.Tx
	var insert_content *sql.Stmt
	var insert_content_result sql.Result
	var insert_content_cast *sql.Stmt
	var insert_content_cast_result sql.Result
	var insert_cast_query string
	var cast_fields []any
	var content_id int64
	var content_fields = get_content_add_fields(obj)
	content_fields = append(content_fields, token.UserID)

	tx, err := DB.Begin()
	if err != nil {
		log.Error("Unable to get DB Instance: ", err.Error())
		return false
	}
	//INSERT CONTENT
	if err := PrepareTransaction(&insert_content, tx, insert_content_query, "Error in Preparing Insert Content"); err != nil {
		return false
	}
	if err := RunTransaction(&insert_content_result, insert_content, tx, content_fields, "Error in Insert Content Query"); err != nil {
		return false
	}
	content_id, err = insert_content_result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return false
	}

	//INSERT CASTS
	insert_cast_query, cast_fields = buildCastQuery(obj.ActorList, int(content_id))
	if err := PrepareTransaction(&insert_content_cast, tx, insert_cast_query, "Error in Preparing Insert Cast"); err != nil {
		return false
	}
	if err := RunTransaction(&insert_content_cast_result, insert_content_cast, tx, cast_fields, "Error in Insert Cast Query"); err != nil {
		return false
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Error("Error in Insert Content Transaction: ", err.Error())
		return false
	}
	return true
}

func (obj *Content) Update(id string) bool {
	var fields = get_content_add_fields(obj)
	fields = append(fields, id)
	return utils.RunUpdate(update_content_query, fields)
}

func (obj Content) Delete(id string) bool {
	return utils.RunUpdate(delete_content_query, []any{id})
}

func (obj Content) DeleteCast(content_id, actor_id string) bool {
	return utils.RunUpdate(delete_cast_query, []any{content_id, actor_id})
}

func (obj Content) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += content_name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["genre"]) > 0 {
		query += genre_filter
		filter = append(filter, params["genre"])
	}
	if len(params["category"]) > 0 {
		query += content_category_filter
		filter = append(filter, params["category"])
	}
	if len(params["type"]) > 0 {
		query += content_type_filter
		filter = append(filter, params["type"])
	}
	if len(params["year"]) > 0 {
		query += content_year_filter
		filter = append(filter, params["year"])
	}
	return query, filter
}

func buildCastQuery(actors *[]int, content_id int) (string, []any) {
	var query string = insert_casts_query
	var fields []any
	var value_array []string
	for _, actor_id := range *actors {
		value_array = append(value_array, "(?, ?)")
		fields = append(fields, content_id, actor_id)
	}
	values := strings.Join(value_array, ",")
	query += values + on_duplicate_cast
	return query, fields
}

/*---------------------------- EPISODES -------------------------------*/
func (episode Episode) Get(offset, limit int, content_id string, params map[string]string) (CountResult, []any, error) {
	var fields = get_episode_fields(&episode)
	var query string = get_episodes_query
	var count_query string = count_episodes_query
	filter_query, filter := episode.QueryBuilder(content_id, params)
	query += filter_query + episode_order_filter
	count_query += filter_query
	return GetListData(filter, fields, &episode, count_query, query, offset, limit)
}

func (episode Episode) GetById(episode_id, content_id string) ([]any, error) {
	var fields = get_episode_fields(&episode)
	var query string = get_episodes_query + content_id + episode_id_filter
	var filter = []any{content_id, episode_id}
	return GetSingleData(filter, fields, &episode, query)
}

func (episode *Episode) Add(content_id string, token UserToken) bool {
	var fields = get_episode_add_fields(episode)
	fields = append(fields, content_id, token.UserID)
	return utils.RunInsert(insert_episode_query, fields)
}

func (episode *Episode) Update(id string) bool {
	var fields = get_episode_add_fields(episode)
	fields = append(fields, id)
	return utils.RunUpdate(update_episode_query, fields)
}

func (episode Episode) Delete(id string) bool {
	return utils.RunUpdate(delete_episode_query, []any{id})
}

func (episode Episode) QueryBuilder(content_id string, params map[string]string) (string, []any) {
	var query string
	var filter = []any{content_id}
	if len(params["name"]) > 0 {
		query += episode_name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["season"]) > 0 {
		query += episode_season_filter
		filter = append(filter, params["season"])
	}
	return query, filter
}

/*---------------------------- RATINGS & REVIEWS -------------------------------*/
func (cr ContentRating) Get(offset, limit int, content_id string, params map[string]string) (CountResult, []any, error) {
	var fields = get_content_rating_fields(&cr)
	var query string = get_content_ratings_query + rating_content_filter
	var count_query string = count_content_rating_query + rating_content_filter
	filter_query, filter := cr.QueryBuilder(content_id, params)
	query += filter_query + rating_order_filter
	count_query += filter_query
	return GetListData(filter, fields, &cr, count_query, query, offset, limit)
}

func (cr *ContentRating) Add(content_id string, token CustomerToken) bool {
	var rating Rating
	var fields = get_content_rating_add_fields(cr)
	var rating_fields = []any{&rating.Rating}
	fields = append(fields, content_id, token.CustomerID)

	if utils.RunInsert(insert_content_rating_query, fields) {
		result, err := GetSingleData([]any{content_id}, rating_fields, &rating, get_average_rating_query)
		if err != nil {
			log.Error("Error in Getting Average Rating: ", err.Error())
		} else {
			rating = utils.FormatSingleResult[Rating](result)
			if !utils.RunUpdate(update_content_average_rating_query, []any{rating.Rating, content_id}) {
				log.Error("Error in Updating Average Rating: ")
			}
		}
	}
	return true
}

func (cr *ContentRating) Update(id string, token CustomerToken) bool {
	var fields = get_content_rating_add_fields(cr)
	fields = append(fields, id, token.CustomerID)
	return utils.RunUpdate(insert_content_rating_query, fields)
}

func (cr ContentRating) Delete(id string, token CustomerToken) bool {
	return utils.RunUpdate(delete_content_rating_query, []any{id, token.CustomerID})
}

func (cr ContentRating) QueryBuilder(content_id string, params map[string]string) (string, []any) {
	var query string
	var filter = []any{content_id}
	if len(params["customer_id"]) > 0 {
		query += rating_customer_filter
		filter = append(filter, params["customer_id"])
	}
	return query, filter
}
