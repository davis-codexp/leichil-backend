package models

func get_user_fields(obj *User) []any {
	return []any{
		&obj.ID, &obj.Name, &obj.Mobile, &obj.Email,
		&obj.Role, &obj.CreatedOn,
	}
}

func get_user_add_fields(obj *User) []any {
	return []any{
		obj.Name, obj.Mobile, obj.Email, obj.Role,
	}
}

func get_producer_fields(obj *Producer) []any {
	return []any{
		&obj.ID, &obj.Name, &obj.Mobile, &obj.Email,
		&obj.Address, &obj.CreatedOn, &obj.CreatedBy,
		&obj.CreatedByName,
	}
}

func get_producer_add_fields(obj *Producer) []any {
	return []any{obj.Name, obj.Mobile, obj.Email, obj.Address}
}

func get_content_fields(obj *Content) []any {
	return []any{
		&obj.ID, &obj.Name, &obj.Genre, &obj.Runtime,
		&obj.ReleaseYear, &obj.Synopsis, &obj.Price,
		&obj.ProducerID, &obj.ProducerName, &obj.Thumbnail,
		&obj.Trailer, &obj.Video, &obj.Category, &obj.Price, &obj.Rating,
		&obj.ContentType, &obj.CreatedOn, &obj.CreatedBy, &obj.CreatedByName,
	}
}

func get_cast_fields(obj *Cast) []any {
	return []any{&obj.ActorID, &obj.ActorName, &obj.Image}
}

func get_content_add_fields(obj *Content) []any {
	return []any{
		obj.Name, obj.Genre, obj.Runtime, obj.ReleaseYear, obj.Synopsis,
		obj.Price, obj.ProducerID, obj.Thumbnail, obj.Trailer, obj.Video,
		obj.Category, obj.ContentType, obj.IsPublished,
	}
}

func get_actor_fields(obj *Actor) []any {
	return []any{
		&obj.ID, &obj.Name, &obj.Image,
		&obj.Gender, &obj.CreatedBy, &obj.CreatedOn,
		&obj.CreatedByName,
	}
}

func get_actor_add_fields(obj *Actor) []any {
	return []any{obj.Name, obj.Image, obj.Gender}
}

func get_customer_fields(obj *Customer) []any {
	return []any{&obj.ID, &obj.Name, &obj.Mobile, &obj.Avatar, &obj.CreatedOn}
}

func get_customer_add_fields(obj *Customer) []any {
	return []any{obj.Name, obj.Mobile}
}

func get_subscription_plan_fields(obj *SubscriptionPlan) []any {
	return []any{
		&obj.ID, &obj.Name, &obj.Features, &obj.Price, &obj.DurationMonths,
		&obj.CreatedOn, &obj.IsActive, &obj.CreatedBy, &obj.CreatedByName,
	}
}

func get_subscription_plan_add_fields(obj *SubscriptionPlan) []any {
	return []any{
		obj.Name, obj.Features, obj.Price, obj.DurationMonths,
		obj.IsActive,
	}
}

func get_subscription_fields(obj *Subscription) []any {
	return []any{
		&obj.ID, &obj.CustomerID, &obj.PlanID, &obj.Price, &obj.ExpiresOn,
		&obj.CreatedOn, &obj.CustomerName, &obj.CustomerMobile, &obj.CustomerAvatar,
		&obj.PlanName,
	}
}

func get_subscription_add_fields(obj *Subscription) []any {
	return []any{obj.PlanID}
}

func get_episode_fields(obj *Episode) []any {
	return []any{
		&obj.ID, &obj.Season, &obj.EpisodeNumber, &obj.Name, &obj.Runtime,
		&obj.ReleaseYear, &obj.Synopsis, &obj.Thumbnail, &obj.Video,
		&obj.Rating, &obj.CreatedOn, &obj.CreatedBy, &obj.CreatedByName,
		&obj.ContentID, &obj.ContentName,
	}
}

func get_episode_add_fields(obj *Episode) []any {
	return []any{
		obj.Season, obj.EpisodeNumber, obj.Name, obj.Runtime, obj.ReleaseYear,
		obj.Synopsis, obj.Thumbnail, obj.Video,
	}
}

func get_cca_fields(obj *CustomerContentActivity) []any {
	return []any{
		&obj.ID, &obj.CustomerID, &obj.ContentID, &obj.EpisodeID,
		&obj.CreatedOn, &obj.ContentName, &obj.EpisodeNumber, &obj.EpisodeName,
		&obj.CustomerName, &obj.CustomerMobile,
	}
}

func get_continue_watching_fields(obj *CustomerContentActivity) []any {
	return []any{
		&obj.ContentID, &obj.CustomerID, &obj.CreatedOn,
		&obj.ContentName, &obj.ContentType, &obj.Thumbnail,
	}
}

func get_content_rating_fields(obj *ContentRating) []any {
	return []any{
		&obj.ContentID, &obj.CustomerID, &obj.Rating, &obj.Review,
		&obj.CreatedOn, &obj.ContentName, &obj.CustomerName, &obj.Mobile,
	}
}

func get_content_rating_add_fields(obj *ContentRating) []any {
	return []any{obj.Rating, obj.Review}
}

func get_purchase_fields(obj *Purchase) []any {
	return []any{
		&obj.CustomerID, &obj.ContentID, &obj.Price, &obj.Discount,
		&obj.NetAmount, &obj.PaymentMode, &obj.CreatedOn, &obj.CustomerName,
		&obj.CustomerMobile, &obj.ContentName, &obj.ContentType,
	}
}

func get_purchase_add_fields(obj *Purchase) []any {
	return []any{
		obj.ContentID, obj.Price, obj.Discount, obj.NetAmount,
		obj.PaymentMode, obj.CustomerID,
	}
}

func get_payment_order_fields(obj *PaymentOrder) []any {
	return []any{
		&obj.OrderID, &obj.PaymentID, &obj.Amount, &obj.Status,
	}
}

func get_payment_order_add_fields(obj *PaymentOrder) []any {
	return []any{obj.OrderID, obj.Amount}
}
