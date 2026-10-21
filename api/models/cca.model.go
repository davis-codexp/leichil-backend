package models

import "leichil-backend/api/utils"

type CustomerContentActivity struct {
	ID             *int    `json:"activity_id,omitempty"`
	CustomerID     *int    `json:"customer_id,omitempty"`
	ContentID      *int    `json:"id,omitempty"`
	ContentType    *string `json:"content_type,omitempty"`
	Thumbnail      *string `json:"content_thumbnail,omitempty"`
	EpisodeID      *int    `json:"episode_id,omitempty"`
	CreatedOn      *string `json:"created_on,omitempty"`
	CustomerName   *string `json:"customer_name,omitempty"`
	CustomerMobile *string `json:"customer_mobile,omitempty"`
	ContentName    *string `json:"content_name,omitempty"`
	EpisodeName    *string `json:"episode_name,omitempty"`
	EpisodeNumber  *int    `json:"episode_number,omitempty"`
}

func (cca CustomerContentActivity) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_cca_fields(&cca)
	var query string = get_customer_content_activities_query
	filter_query, filter := cca.QueryBuilder(params)
	query += filter_query + cca_order_filter
	return GetListData(filter, fields, &cca, "", query, offset, limit)
}

func (cca *CustomerContentActivity) Add(token CustomerToken) bool {
	var fields = []any{token.CustomerID, cca.ContentID, cca.EpisodeID}
	return utils.RunInsert(insert_content_activity_query, fields)
}

func (cca CustomerContentActivity) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["content"]) > 0 {
		query += cca_content_filter
		filter = append(filter, params["content"])
	}
	if len(params["customer"]) > 0 {
		query += cca_customer_filter
		filter = append(filter, params["customer"])
	}
	if len(params["episode"]) > 0 {
		query += cca_content_filter + cca_episode_filter
		filter = append(filter, params["content"], params["episode"])
	}
	if len(params["start_date"]) > 0 && len(params["end_date"]) > 0 {
		query += cca_date_range_filter
		filter = append(filter, params["start_date"], params["end_date"])
	}
	return query, filter
}

/*--------------------------CUSTOMER FAVOURITES-------------------------*/
type CustomerFavourite struct {
	CustomerID   *int     `json:"customer_id,omitempty"`
	ContentID    *int     `json:"content_id,omitempty"`
	CreatedOn    *string  `json:"created_on,omitempty"`
	CustomerName *string  `json:"customer_name,omitempty"`
	ContentName  *string  `json:"name,omitempty"`
	Thumbnail    *string  `json:"thumbnail,omitempty"`
	Video        *string  `json:"video,omitempty"`
	Runtime      *int     `json:"runtime,omitempty"`
	Genre        *string  `json:"genre,omitempty"`
	ReleaseYear  *int     `json:"release_year,omitempty"`
	Rating       *float32 `json:"rating,omitempty"`
}

func (fav CustomerFavourite) Get(offset, limit int, params map[string]string, token CustomerToken) (CountResult, []any, error) {
	var fields []any = get_fav_fields(&fav)
	var query string = get_customer_favourites_query
	filter_query, filter := fav.QueryBuilder(params, token)
	query += filter_query + fav_order_filter
	return GetListData(filter, fields, &fav, "", query, offset, limit)
}

func (fav *CustomerFavourite) Add(token CustomerToken) bool {
	var fields = []any{token.CustomerID, fav.ContentID}
	return utils.RunInsert(insert_customer_favourite_query, fields)
}

func (fav CustomerFavourite) Delete(content_id string, token CustomerToken) bool {
	return utils.RunUpdate(delete_customer_favourite_query, []any{token.CustomerID, content_id})
}

func (fav CustomerFavourite) QueryBuilder(params map[string]string, token CustomerToken) (string, []any) {
	var query string
	var filter = []any{token.CustomerID}
	if len(params["content"]) > 0 {
		query += fav_content_filter
		filter = append(filter, params["content"])
	}
	return query, filter
}
