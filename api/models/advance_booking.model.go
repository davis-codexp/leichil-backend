package models

import "leichil-backend/api/utils"

type AdvanceBooking struct {
	ContentID    *int    `json:"content_id,omitempty"`
	Price        *int    `json:"price,omitempty"`
	StartDate    *string `json:"start_date,omitempty"`
	EndDate      *string `json:"end_date,omitempty"`
	PremiereDate *string `json:"premiere_date,omitempty"`
	CreatedOn    *string `json:"created_on,omitempty"`
	CreatedBy    *int    `json:"created_by,omitempty"`
	ContentName  *string `json:"content_name,omitempty"`
	ContentPrice *int    `json:"content_price,omitempty"`
	Thumbnail    *string `json:"thumbnail,omitempty"`
	Genre        *string `json:"genre,omitempty"`
}

func (obj AdvanceBooking) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_advance_booking_fields(&obj)
	var query string = get_advance_bookings_query
	filter_query, filter := obj.QueryBuilder(params)
	query += filter_query + " ORDER BY ab.created_on DESC"
	return GetListData(filter, fields, &obj, "", query, offset, limit)
}

func (obj AdvanceBooking) GetByContentId(content_id string) ([]any, error) {
	var fields = get_advance_booking_fields(&obj)
	var query string = get_advance_bookings_query + advance_booking_content_filter
	var filter = []any{content_id}
	return GetSingleData(filter, fields, &obj, query)
}

func (obj *AdvanceBooking) Add(token UserToken) bool {
	var fields = get_advance_booking_add_fields(obj)
	fields = append(fields, token.UserID)
	return utils.RunInsert(insert_advance_booking_query, fields)
}

func (obj *AdvanceBooking) Update(content_id string) bool {
	var fields = get_advance_booking_add_fields(obj)
	fields = append(fields, content_id)
	return utils.RunUpdate(update_advance_booking_query, fields)
}

func (obj AdvanceBooking) Delete(content_id string) bool {
	return utils.RunUpdate(delete_advance_booking_query, []any{content_id})
}

func (obj AdvanceBooking) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["content"]) > 0 {
		query += advance_booking_content_filter
		filter = append(filter, params["content"])
	}
	if len(params["active"]) > 0 {
		query += advance_booking_active_filter
	}
	return query, filter
}
