package models

import (
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2/log"
)

type User struct {
	ID        *int    `json:"id,omitempty"`
	Name      *string `form:"name" json:"name,omitempty"`
	Mobile    *string `form:"mobile" json:"mobile,omitempty"`
	Email     *string `form:"email" json:"email,omitempty"`
	Password  *string `json:"password,omitempty"`
	Role      *string `form:"role" json:"role,omitempty"`
	CreatedOn *string `json:"created_on,omitempty"`
}

func (user User) Get(offset, limit int, params map[string]string) (CountResult, []any, error) {
	var fields = get_user_fields(&user)
	var query string = get_users_query
	var count_query string = count_users_query
	filter_query, filter := user.QueryBuilder(params)
	query += filter_query + order_filter
	count_query += filter_query
	return GetListData(filter, fields, &user, count_query, query, offset, limit)
}

func (user User) GetById(id string) ([]any, error) {
	var fields = get_user_fields(&user)
	var query string = get_users_query + id_filter
	var filter = []any{id}
	return GetSingleData(filter, fields, &user, query)
}

func (user *User) Add() bool {
	var fields = get_user_add_fields(user)
	password, hash, err := utils.GeneratePassword()
	if err != nil {
		log.Error("Error while hashing password: ", err.Error())
		return false
	}
	log.Info(password)
	fields = append(fields, string(hash))
	return utils.RunInsert(insert_user_query, fields)
}

func (user *User) Update(id string) bool {
	var fields = get_user_add_fields(user)
	fields = append(fields, id)
	return utils.RunUpdate(update_user_query, fields)
}

func (user User) Delete(id string) bool {
	return utils.RunUpdate(delete_user_query, []any{id})
}

func (user User) QueryBuilder(params map[string]string) (string, []any) {
	var query string
	var filter []any
	if len(params["name"]) > 0 {
		query += name_filter
		filter = append(filter, "%"+params["name"]+"%")
	}
	if len(params["mobile"]) > 0 {
		query += mobile_filter
		filter = append(filter, params["mobile"])
	}
	if len(params["role"]) > 0 {
		query += user_role_filter
		filter = append(filter, params["role"])
	}
	if len(params["email"]) > 0 {
		query += email_filter
		filter = append(filter, params["email"])
	}
	return query, filter
}
