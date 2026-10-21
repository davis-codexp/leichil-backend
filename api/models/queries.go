package models

var id_filter = " AND id = ?"
var name_filter = " AND name LIKE ?"
var mobile_filter = " AND mobile = ?"
var email_filter = " AND email = ?"
var genre_filter = " AND genre = ?"
var order_filter = " ORDER BY created_on DESC"

/*-------------USER QUERIES---------*/
var get_users_query = `
	SELECT id, name, mobile, email, role, created_on
	FROM users WHERE is_deleted = 0
`
var count_users_query = `SELECT COUNT(id) FROM users WHERE is_deleted = 0`
var insert_user_query = `
	INSERT INTO users(name, mobile, email, role, password) VALUES(?, ?, ?, ?, ?)
`
var update_user_query = `
	UPDATE users SET name = ?, mobile = ?, email = ?, role = ?
	WHERE id = ?
`
var delete_user_query = `
	UPDATE users SET is_deleted = 1 WHERE id = ?
`
var user_role_filter = " AND role = ?"

/*-------------PRODUCER QUERIES---------*/
var get_producers_query = `
	SELECT p.id, p.name, p.mobile, p.email, p.address, p.created_on,
	p.created_by, u.name AS created_by_name
	FROM producers p
	INNER JOIN users u ON u.id = p.created_by
	WHERE p.is_deleted = 0
`
var count_producers_query = `
	SELECT COUNT(p.id) AS COUNT
	FROM producers p
	INNER JOIN users u ON u.id = p.created_by
	WHERE p.is_deleted = 0
`
var insert_producer_query = `INSERT INTO producers(name, mobile, email, address, created_by) VALUES(?, ?, ?, ?, ?)`
var update_producer_query = `
	UPDATE producers SET name = ?, mobile = ?, email = ? address = ?
	WHERE id = ?
`
var delete_producer_query = `
	UPDATE producers SET is_deleted = 1 WHERE id = ?
`
var producer_id_filter = " AND p.id = ?"
var producer_name_filter = " AND p.name LIKE ?"
var producer_mobile_filter = " AND p.mobile = ?"
var producer_email_filter = " AND p.email = ?"
var order_producer_filter = " ORDER BY p.created_on DESC"

/*-------------CONTENT QUERIES---------*/
var get_contents_query = `
	SELECT contents.id, contents.name, genre, runtime, release_year, synopsis,
	price, producer_id, producers.name AS producer_name,
	thumbnail, trailer, video, category, price, rating, content_type,
	contents.created_on, contents.created_by, users.name AS created_by_name
	FROM contents
	INNER JOIN producers ON producers.id = contents.producer_id
	INNER JOIN users ON users.id = contents.created_by
	WHERE contents.is_deleted = 0
`
var get_content_price_query = "SELECT id, price FROM contents WHERE id = ?"
var get_home_contents_query = `
	SELECT contents.id, contents.name, thumbnail, release_year, genre, rating,
	price, category
	FROM contents
	INNER JOIN producers ON producers.id = contents.producer_id
	INNER JOIN users ON users.id = contents.created_by
	WHERE contents.is_deleted = 0 AND contents.is_published = 1
`
var count_contents_query = `
	SELECT COUNT(contents.id) AS COUNT
	FROM contents
	INNER JOIN producers ON producers.id = contents.producer_id
	INNER JOIN users ON users.id = contents.created_by
	WHERE contents.is_deleted = 0
`
var insert_content_query = `
	INSERT INTO contents(
		name, genre, runtime, release_year, synopsis, price,
		producer_id, thumbnail, trailer, video, category, content_type,
		is_published, created_by
	) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
var update_content_query = `
	UPDATE contents SET name = ?, genre = ?, runtime = ?,
	release_year = ?, synopsis = ?, price = ?, producer_id = ?, thumbnail = ?,
	trailer = ?, video = ?, category = ?, content_type = ?, is_published = ?
	WHERE id = ?
`
var update_content_average_rating_query = "UPDATE contents SET rating = ? WHERE id = ?"
var delete_content_query = `UPDATE contents SET is_deleted = 1 WHERE id = ?`
var content_id_filter = " AND contents.id = ?"
var content_name_filter = " AND contents.name LIKE ?"
var content_type_filter = " AND contents.content_type = ?"
var content_category_filter = " AND contents.category = ?"
var content_year_filter = " AND contents.release_year = ?"
var content_published_filter = " AND contents.is_published = ?"
var content_order_filter = " ORDER BY contents.created_on DESC"

/*---------------------------EPISODE QUERIES--------------------------*/
var get_episodes_query = `
	SELECT ep.id, season, episode_number, ep.name, ep.runtime,
	ep.release_year, ep.synopsis, ep.thumbnail, ep.video, ep.rating, ep.created_on,
	ep.created_by, users.name AS created_by_name, content_id,
	contents.name AS content_name
	FROM episodes ep
	INNER JOIN contents ON contents.id = ep.content_id
	INNER JOIN users ON users.id = contents.created_by
	WHERE ep.is_deleted = 0 AND ep.content_id = ?
`
var count_episodes_query = `
	SELECT COUNT(*) AS COUNT
	FROM episodes ep
	INNER JOIN contents ON contents.id = ep.content_id
	INNER JOIN users ON users.id = contents.created_by
	WHERE ep.is_deleted = 0 AND ep.content_id = ?
`
var insert_episode_query = `
	INSERT INTO episodes(
		season, episode_number, name, runtime, release_year,
		synopsis, thumbnail, video, content_id, created_by
	) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
var update_episode_query = `
	UPDATE contents SET season = ?, episode_number = ?, name = ?, runtime = ?,
	release_year = ?, synopsis = ?, thumbnail = ?, video = ?
	WHERE id = ?
`
var delete_episode_query = `UPDATE episodes SET is_deleted = 1 WHERE id = ?`
var episode_id_filter = " AND ep.id = ?"
var episode_name_filter = " AND ep.name LIKE ?"
var episode_season_filter = " AND ep.season = ?"
var episode_order_filter = " ORDER BY ep.created_on DESC"

/*---------------------------RATING & REVIEW QUERIES--------------------------*/
var get_content_ratings_query = `
	SELECT cr.content_id, cr.customer_id, cr.rating, cr.review, cr.created_on,
	con.name AS content_name, cus.name AS customer_name, cus.mobile
	FROM content_ratings cr
	INNER JOIN contents con ON con.id = cr.content_id
	INNER JOIN customers cus ON cus.id = cr.customer_id
	WHERE cr.is_deleted = 0
`
var count_content_rating_query = `
	SELECT COUNT(*) AS COUNT FROM content_ratings cr
	INNER JOIN contents con ON con.id = cr.content_id
	INNER JOIN customers cus ON cus.id = cr.customer_id
	WHERE cr.is_deleted = 0
`
var get_average_rating_query = `
	SELECT ROUND(AVG(CAST(rating AS FLOAT)), 2) AS rating
	FROM content_ratings
	WHERE is_deleted = 0 AND content_id = ?
`
var rating_content_filter = " AND cr.content_id = ?"
var rating_customer_filter = " AND cr.customer_id = ?"
var rating_order_filter = " ORDER BY cr.created_on DESC"
var insert_content_rating_query = `
	INSERT INTO content_ratings(rating, review, content_id, customer_id) values(?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE rating = VALUES(rating), review = VALUES(review)
`
var update_content_rating_query = `
	UPDATE content_ratings SET rating = ? AND review = ? WHERE content_id = ? AND customer_id = ?
`
var delete_content_rating_query = `
	UPDATE content_ratings SET is_deleted = 1 WHERE content_id = ? AND customer_id = ? 
`

/*---------------------------ACTOR QUERIES-------------------------*/
var get_actors_query = `
	SELECT a.id, a.name, a.image, a.gender, a.created_by, a.created_on,
	u.name AS created_by_name
	FROM actors a
	INNER JOIN users u ON u.id = a.created_by
	WHERE a.is_deleted = 0
`
var count_actors_query = `
	SELECT COUNT(a.id) AS COUNT
	FROM actors a
	INNER JOIN users u ON u.id = a.created_by
	WHERE a.is_deleted = 0
`
var actor_id_filter = " AND a.id = ?"
var order_actors_filter = " ORDER BY a.created_on DESC"
var actor_gender_filter = " AND a.gender = ?"
var actor_name_filter = " AND a.name LIKE ?"
var insert_actor_query = `INSERT INTO actors(name, image, gender, created_by) VALUES(?, ?, ?, ?)`
var update_actor_query = `
	UPDATE actors SET name = ?, image = ?, gender = ?
	WHERE id = ?
`
var delete_actor_query = `UPDATE actors SET is_deleted = 1 WHERE id = ?`

/*-------------CUSTOMER QUERIES---------*/
var get_customers_query = `
	SELECT id, name, mobile, avatar, created_on
	FROM customers WHERE is_deleted = 0
`
var count_customers_query = `SELECT COUNT(id) AS COUNT FROM customers WHERE is_deleted = 0`
var insert_customer_query = `INSERT INTO customers(name, mobile) VALUES(?, ?)`
var update_customer_query = `UPDATE customers SET name = ?, avatar = ? WHERE id = ?`
var delete_customer_query = `UPDATE customers SET is_deleted = 1 WHERE id = ?`

/*-------------MOVIE CASTS QUERIES---------*/
var get_casts_query = `
	SELECT actors.id, actors.name, actors.image
	FROM casts
	INNER JOIN actors ON actors.id = casts.actor_id
	WHERE content_id = ?
`
var insert_casts_query = `INSERT INTO casts(content_id, actor_id) VALUES`
var on_duplicate_cast = ` ON DUPLICATE KEY UPDATE content_id = VALUES(content_id),
	actor_id = VALUES(actor_id)
`
var delete_cast_query = "DELETE FROM casts WHERE content_id = ? AND actor_id = ?"

/*-------------AUTH QUERIES---------*/
var user_login_query = `
	SELECT id, name, mobile, email, role, password
	FROM users WHERE is_deleted = 0 AND mobile = ?
`
var get_user_token_query = `
	SELECT ut.token, ut.user_id, ut.created_on,
	u.name, u.mobile, u.email, u.role
	FROM user_tokens ut
	INNER JOIN users u ON u.id = ut.user_id
	WHERE ut.token = ?
`
var insert_user_token_query = `
	INSERT INTO user_tokens(token, user_id) VALUES(?, ?)
`
var delete_user_token_query = "DELETE FROM user_tokens WHERE token = ?"

var customer_login_query = `
	SELECT id, name, mobile FROM customers WHERE is_deleted = 0 AND mobile = ?
`
var get_customer_token_query = `
	SELECT ct.token, ct.customer_id, ct.created_on,
	c.name, c.mobile
	FROM customer_tokens ct
	INNER JOIN customers c ON c.id = ct.customer_id
	WHERE ct.token = ?
`
var insert_customer_token_query = "INSERT INTO customer_tokens(token, customer_id) VALUES(?, ?)"
var delete_customer_token_query = "DELETE FROM customer_tokens WHERE token = ?"
var insert_customer_otp_query = "INSERT INTO customer_otps(otp, mobile, customer_id) VALUES(?, ?, ?)"
var verify_customer_otp_query = `
	SELECT co.otp, co.mobile, co.customer_id, c.name
	FROM customer_otps co
	INNER JOIN customers c ON c.id = co.customer_id
	WHERE co.mobile = ? AND co.otp = ?
`
var delete_customer_otp_query = "DELETE FROM customer_otps WHERE mobile = ?"

/*-------------SUBSCRIPTION PLAN QUERIES---------*/
var get_subscription_plans_query = `
	SELECT sp.id, sp.name, sp.features, sp.price, sp.duration_months, sp.created_on, sp.is_active,
	sp.created_by, u.name AS created_by_name
	FROM subscription_plans sp
	INNER JOIN users u ON sp.created_by = u.id
	WHERE sp.is_deleted = 0
`
var count_subscription_plans_query = `
	SELECT COUNT(*) AS COUNT
	FROM subscription_plans sp
	INNER JOIN users u ON sp.created_by = u.id
	WHERE sp.is_deleted = 0
`
var subscription_plan_id_filter = " AND sp.id = ?"
var subscription_plan_name_filter = " AND sp.name LIKE ?"
var subscription_plan_price_filter = " AND sp.price = ?"
var subscription_plan_duration_filter = " AND sp.duration_months = ?"
var subscription_plan_status_filter = " AND sp.is_active = ?"
var insert_subscription_plan_query = `
	INSERT INTO subscription_plans(
		name, features, price, duration_months, is_active, created_by
	) VALUES(?, ?, ?, ?, ?, ?)
`
var update_subscription_plan_query = `
	UPDATE subscription_plans SET name = ?, features = ?, price = ?, duration_months = ?, is_active = ? WHERE id = ?
`
var delete_subscription_plan_query = "UPDATE subscription_plans SET is_deleted = 1 WHERE id = ?"

/*-------------SUBSCRIPTION QUERIES---------*/
var get_subscriptions_query = `
	SELECT sub.id, sub.customer_id, sub.plan_id, sub.price, sub.expires_on, sub.created_on,
	cus.name AS customer_name, cus.mobile AS customer_mobile, cus.avatar AS customer_avatar,
	sp.name as plan_name
	FROM subscriptions sub
	INNER JOIN customers cus ON sub.customer_id = cus.id
	INNER JOIN subscription_plans sp ON sp.id = sub.plan_id
	WHERE sub.is_deleted = 0
`
var count_subscriptions_query = `
	SELECT COUNT(*) AS COUNT
	FROM subscriptions sub
	INNER JOIN customers cus ON sub.customer_id = cus.id
	INNER JOIN subscription_plans sp ON sp.id = sub.plan_id
	WHERE sub.is_deleted = 0
`
var subscription_id_filter = " AND sub.id = ?"
var subscription_customer_filter = " AND sub.customer_id = ?"
var subscription_plan_filter = " AND sub.plan_id = ?"
var subscription_date_filter = " AND sub.created_on BETWEEN ? AND ?"
var subscription_expired_filter = " AND sub.expires_on < CURDATE()"
var subscription_not_expired_filter = " AND sub.expires_on >= CURDATE()"
var insert_subscription_query = `
	INSERT INTO subscriptions(plan_id, price, customer_id, expires_on) VALUES(?, ?, ?, ?)
`
var update_subscription_query = `
	UPDATE subscriptions SET plan_id = ?, price = ?, expires_on = ? WHERE id = ?
`
var delete_subscription_query = "UPDATE subscriptions SET is_deleted = 1 WHERE id = ?"
var extend_subscription_expiry_query = "UPDATE subscriptions SET expires_on = expires_on + INTERVAL ? DAY WHERE id = ?"

/*------------------------CUSTOMER CONTENT ACTIVITIES--------------------*/
var insert_content_activity_query = `
	INSERT INTO customer_content_activities(customer_id, content_id, episode_id) VALUES(?, ?, ?)
`
var get_trending_contents_query = `
	SELECT cca.content_id, con.name, con.thumbnail, con.release_year,
	con.genre, con.rating, con.price, con.category
	FROM customer_content_activities cca
	INNER JOIN contents con ON con.id = cca.content_id
	where con.is_deleted = 0 AND con.is_published = 1
	AND cca.created_on >= NOW() - INTERVAL 30 DAY
`
var get_continue_watching_query = `
	SELECT DISTINCT(cca.content_id), cca.customer_id, cca.created_on, con.name, con.content_type, con.thumbnail
	FROM customer_content_activities cca
	INNER JOIN contents con ON con.id = cca.content_id
	where con.is_deleted = 0 AND con.is_published = 1
	AND cca.created_on >= NOW() - INTERVAL 30 DAY
`
var get_customer_content_activities_query = `
	SELECT cca.id, cca.customer_id, cca.content_id, cca.episode_id, cca.created_on,
	con.name AS content_name, ep.episode_number, ep.name AS episode_name,
	cus.name AS customer_name, cus.mobile AS customer_mobile
	FROM customer_content_activities cca
	INNER JOIN contents con ON con.id = cca.content_id
	INNER JOIN episodes ep ON ep.id = cca.episode_id
	WHERE con.is_deleted = 0 AND con.is_published = 1
`
var cca_content_type_filter = " AND con.content_type = ?"
var cca_customer_filter = " AND cca.customer_id = ?"
var cca_content_filter = " AND cca.content_id = ?"
var cca_episode_filter = " AND cca.episode_id = ?"
var cca_date_range_filter = " AND cca.created_on BETWEEN ? AND ?"
var cca_order_filter = " ORDER BY cca.created_on DESC"

/*------------------------CUSTOMER FAVOURITES--------------------*/
var get_customer_favourites_query = `
	SELECT fav.customer_id, fav.content_id, fav.created_on,
	cus.name, con.name, con.thumbnail, con.genre,
	con.runtime, con.video, con.release_year, con.rating
	FROM customer_favourites fav
	INNER JOIN customers cus ON cus.id = fav.customer_id
	INNER JOIN contents con ON con.id = fav.content_id
	WHERE fav.customer_id = ?
`
var insert_customer_favourite_query = `
	INSERT IGNORE INTO customer_favourites(customer_id, content_id) VALUES(?, ?)
`
var delete_customer_favourite_query = `
	DELETE FROM customer_favourites WHERE customer_id = ? AND content_id = ?
`
var fav_order_filter = " ORDER BY fav.created_on DESC"
var fav_content_filter = " AND fav.content_id = ?"

/*------------------------PAYMENT ORDER QUERIES--------------------*/
var get_payment_orders_query = `
	SELECT po.order_id, po.payment_id, po.customer_id, po.content_id,
	po.amount, po.status, po.plan_id
	FROM payment_orders po
	WHERE po.order_id IS NOT NULL
`
var insert_payment_order_query = `
	INSERT INTO payment_orders(
		order_id, amount, content_id, plan_id, customer_id
	) VALUES(?, ?, ?, ?, ?)
`
var update_payment_order_query = "UPDATE payment_orders SET status = ?, payment_id = ? WHERE order_id = ?"
var payment_order_id_filter = " AND po.order_id = ?"
var payment_order_payment_filter = " AND po.payment_id = ?"
var payment_order_amount_filter = " AND po.amount = ?"
var payment_order_status_filter = " AND po.status = ?"
var payment_content_filter = " AND po.content_id = ?"
var payment_customer_filter = " AND po.customer_id = ?"

/*------------------------DEVICE TOKEN QUERIES--------------------*/
var get_device_tokens_query = `
	SELECT device_token, customer_id, created_on
	FROM device_tokens
	WHERE customer_id = ?
`
var insert_device_token_query = "INSERT INTO device_tokens(device_token, customer_id) VALUES(?, ?)"
var delete_device_token_query = "DELETE FROM device_tokens WHERE device_token = ?"
