CREATE TABLE users (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(100),
	email VARCHAR(100),
	mobile VARCHAR(10),
	password VARCHAR(100),
	role ENUM('ADMIN', 'STAFF') DEFAULT 'STAFF',
	is_deleted TINYINT(1) DEFAULT 0,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) AUTO_INCREMENT=1001;

CREATE TABLE producers (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(500),
	mobile VARCHAR(10),
	email VARCHAR(100),
	address VARCHAR(500),
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by INT,
	is_deleted TINYINT(1) DEFAULT 0,
	FOREIGN KEY(created_by) REFERENCES users(id)
) AUTO_INCREMENT = 1001;

CREATE TABLE contents (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(500),
	genre VARCHAR(100),
	runtime INT(4),
	release_year INT(4),
	synopsis VARCHAR(1000),
	price INT(4),
	producer_id INT,
	thumbnail VARCHAR(500),
	trailer VARCHAR(500),
	video VARCHAR(500),
	category ENUM('PREMIERE', 'NORMAL') DEFAULT 'NORMAL',
	content_type ENUM('MOVIE', 'SHUMANG LILA', 'SERIES', 'SHORT FILM') DEFAULT 'MOVIE',
	rating FLOAT(2,1) DEFAULT 0,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0,
	created_by INT,
	FOREIGN KEY(producer_id) REFERENCES producers(id) ON DELETE CASCADE,
	FOREIGN KEY(created_by) REFERENCES users(id)
) AUTO_INCREMENT = 1001;

CREATE TABLE episodes (
	id INT PRIMARY KEY AUTO_INCREMENT,
	content_id INT,
	season INT,
	episode_number INT,
	name VARCHAR(500),
	thumbnail VARCHAR(500),
	runtime INT(4),
	release_year INT(4),
	synopsis VARCHAR(1000),
	video VARCHAR(500),
	rating FLOAT(2, 1) DEFAULT 0,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0,
	created_by INT,
	FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
	FOREIGN KEY(created_by) REFERENCES users(id)
);

CREATE TABLE actors (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(500),
	image VARCHAR(500),
	gender ENUM('MALE', 'FEMALE', 'TRANSGENDER'),
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0,
	created_by INT,
	FOREIGN KEY(created_by) REFERENCES users(id)
) AUTO_INCREMENT = 1001;

CREATE TABLE casts (
	content_id INT,
	actor_id INT,
	PRIMARY KEY(content_id, actor_id),
	FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
	FOREIGN KEY(actor_id) REFERENCES actors(id)
);

CREATE TABLE customers (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(300),
	mobile VARCHAR(10),
    avatar VARCHAR(100),
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0
) AUTO_INCREMENT = 1001;

CREATE TABLE content_ratings (
	content_id INT,
	customer_id INT,
	rating INT DEFAULT 0,
	review VARCHAR(500),
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0,
	PRIMARY KEY(content_id, customer_id),
	FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
	FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE user_tokens (
	token VARCHAR(100) PRIMARY KEY,
	user_id INT,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE customer_tokens (
	token VARCHAR(100) PRIMARY KEY,
	customer_id INT,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE customer_otps (
	otp INT PRIMARY KEY,
	mobile VARCHAR(10),
	customer_id INT,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE subscription_plans (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(200),
	features VARCHAR(500),
	price INT,
	duration_months INT DEFAULT 1,
	is_active TINYINT(1) DEFAULT 1,
	is_deleted TINYINT(1) DEFAULT 0,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by INT,
	FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
) AUTO_INCREMENT=1001;

CREATE TABLE subscriptions (
	id INT PRIMARY KEY AUTO_INCREMENT,
	customer_id INT,
	plan_id INT,
	price INT,
	expires_on DATE,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_deleted TINYINT(1) DEFAULT 0,
	FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE,
	FOREIGN KEY(plan_id) REFERENCES subscription_plans(id) ON DELETE CASCADE
) AUTO_INCREMENT = 1001;

CREATE TABLE purchases (
    customer_id INT,
    content_id INT,
    price INT,
    discount INT,
    net_amount INT,
    payment_mode ENUM("CARD", "UPI", "NET BANKING") DEFAULT "UPI",
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(customer_id, content_id),
    FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY(content_id) REFERENCES contents(id)
);

CREATE TABLE customer_content_activities(
	id INT PRIMARY KEY AUTO_INCREMENT,
	customer_id INT,
	content_id INT,
	episode_id INT DEFAULT NULL,
	created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE,
	FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
	FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE CASCADE
) AUTO_INCREMENT = 1001;

CREATE TABLE customer_favourites(
    customer_id INT,
    content_id INT,
    created_on TIMESTAMP NOT NULL,
    PRIMARY KEY(customer_id, content_id),
    FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE
);
