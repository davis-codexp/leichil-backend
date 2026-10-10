package models

import (
	"leichil-backend/api/utils"

	"github.com/gofiber/fiber/v2/log"
)

type DeviceToken struct {
	DeviceToken *string `json:"device_token,omitempty"`
	CustomerID  *int    `json:"customer_id,omitempty"`
	CreatedOn   *string `json:"created_on,omitempty"`
}

func (obj DeviceToken) Get(token CustomerToken) (CountResult, []any, error) {
	var fields = []any{&obj.DeviceToken, &obj.CustomerID, &obj.CreatedOn}
	return GetListData([]any{token.CustomerID}, fields, &obj, "", get_device_tokens_query, 0, 10000)
}

func (obj *DeviceToken) Add(token CustomerToken) bool {
	var fields = []any{obj.DeviceToken, token.CustomerID}
	return utils.RunInsert(insert_device_token_query, fields)
}

func (obj DeviceToken) Delete(device_token string) bool {
	if err := utils.RunDelete(delete_device_token_query, []any{device_token}); err != nil {
		log.Error("Error in Deleting Device Token: ", err.Error())
		return false
	}
	return true
}
