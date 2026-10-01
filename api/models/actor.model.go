package models

import "leichil-backend/api/utils"

type Actor struct {
	ID            *int    `json:"id"`
	Name          *string `form:"name" json:"name"`
	Image         *string `form:"image" json:"image"`
	Gender        *string `form:"gender" json:"gender,omitempty"`
	CreatedOn     *string `json:"created_on"`
	CreatedBy     *int    `json:"created_by,omitempty"`
	CreatedByName *string `json:"created_by_name"`
}

func (obj Actor) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_actor_fields(&obj)
	var query string = get_actors_query
	var count_query string = count_actors_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query + order_actors_filter
	count_query += filter_query
	return GetListData(filter, fields, &obj, count_query, query, offset, limit)
}

func (obj Actor) GetById(id string) ([]any, error) {
	var fields = get_actor_fields(&obj)
	var query string = get_actors_query + actor_id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &obj, query)
}

func (obj *Actor) Add(token UserToken) bool {
	var fields = get_actor_add_fields(obj)
	fields = append(fields, token.UserID)
	return utils.RunInsert(insert_actor_query, fields)
}

func (obj *Actor) Update(id string) bool {
	var fields = get_actor_add_fields(obj)
	fields = append(fields, id)
	return utils.RunUpdate(update_actor_query, fields)
}

func (obj Actor) Delete(id string) bool {
	return utils.RunUpdate(delete_actor_query, []any{id})
}

func (obj Actor) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += actor_name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["gender"]) > 0 {
		query += actor_gender_filter
		filter = append(filter, params["gender"])
	}
	return query, filter
}
