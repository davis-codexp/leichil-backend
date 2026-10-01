package models

import "leichil-backend/api/utils"

type Producer struct {
	ID            *int    `json:"id"`
	Name          *string `form:"name" json:"name"`
	Mobile        *string `form:"mobile" json:"mobile"`
	Email         *string `form:"email" json:"email"`
	Address       *string `form:"address" json:"address"`
	CreatedOn     *string `json:"added_on"`
	CreatedBy     *string `json:"added_by"`
	CreatedByName *string `json:"added_by_name"`
}

func (producer Producer) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_producer_fields(&producer)
	var query string = get_producers_query
	var count_query string = count_producers_query
	filter_query, filter := producer.QueryBuilder(params)
	query += filter_query + order_producer_filter
	count_query += filter_query
	return GetListData(filter, fields, &producer, count_query, query, offset, limit)
}

func (producer Producer) GetById(id string) ([]any, error) {
	var fields = get_producer_fields(&producer)
	var query string = get_producers_query + producer_id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &producer, query)
}

func (producer *Producer) Add(token UserToken) bool {
	var fields = get_producer_add_fields(producer)
	fields = append(fields, token.UserID)
	return utils.RunInsert(insert_producer_query, fields)
}

func (producer *Producer) Update(id string) bool {
	var fields = get_producer_add_fields(producer)
	fields = append(fields, id)
	return utils.RunUpdate(update_producer_query, fields)
}

func (producer Producer) Delete(id string) bool {
	return utils.RunUpdate(delete_producer_query, []any{id})
}

func (producer Producer) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += producer_name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["mobile"]) > 0 {
		query += producer_mobile_filter
		filter = append(filter, params["mobile"])
	}
	if len(params["email"]) > 0 {
		query += producer_email_filter
		filter = append(filter, params["email"])
	}
	return query, filter
}
