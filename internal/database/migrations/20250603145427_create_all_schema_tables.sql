-- +goose Up
CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TYPE order_status AS ENUM ('order_placed', 'in_progress', 'shipped', 'delivered', 'cancelled');

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	first_name VARCHAR NOT NULL,
	last_name VARCHAR,
	image TEXT,
	email varchar UNIQUE,
	role role,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE local_auth (
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	phone_number VARCHAR NOT NULL UNIQUE,
	is_phone_verified BOOLEAN DEFAULT false,
	password_hash TEXT NOT NULL,

	PRIMARY KEY (user_id)
);

CREATE TABLE oauth (
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	provider VARCHAR NOT NULL,
	provider_id UUID UNIQUE,

	PRIMARY KEY (user_id)
);

CREATE TABLE sessions (
	id SERIAL PRIMARY KEY,
	revoked BOOLEAN DEFAULT false,
	user_agent varchar UNIQUE NOT NULL,
	refresh_token TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE phone_verification_codes (
	id SERIAL PRIMARY KEY,
	otp_code VARCHAR NOT NULL,
	is_used BOOLEAN DEFAULT false,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	
	user_id INT NOT NULL UNIQUE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE password_resets (
	id SERIAL PRIMARY KEY,
	otp_code VARCHAR NOT NULL,
	is_used BOOLEAN DEFAULT false,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE categories (
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE NOT NULL,
	parent_id INT REFERENCES categories(id)
);

CREATE TABLE brands(
	id SERIAL PRIMARY KEY,
	brand VARCHAR UNIQUE NOT NULL
);

CREATE TABLE products (
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE NOT NULL,
	details TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	
	brand_id int NOT NULL,
	product_category INT NOT NULL,	
	FOREIGN KEY (product_category) REFERENCES categories(id) ON DELETE CASCADE,
	FOREIGN KEY (brand_id) REFERENCES brands(id) 
);

CREATE TABLE sizes (
	id SERIAL PRIMARY KEY,
	size VARCHAR NOT NULL,
	label VARCHAR NOT NULL,
	
	UNIQUE(size, label)
);

CREATE TABLE colors (
	id SERIAL PRIMARY KEY,
	color VARCHAR NOT NULL UNIQUE
);

CREATE TABLE product_variants (
	id SERIAL PRIMARY KEY,
	quantity INT NOT NULL,
	price INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	product_id INT NOT NULL,
	color_id INT NOT NULL,
	size_id INT NOT NULL,

	UNIQUE(product_id, color_id, size_id),
	FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
	FOREIGN KEY(color_id) REFERENCES colors(id),
	FOREIGN KEY(size_id) REFERENCES sizes(id)
);

CREATE TABLE rating_review (
	id SERIAL PRIMARY KEY,
	rating INT NOT NULL,
	review VARCHAR,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	user_id INT NOT NULL,
	product_id INT NOT NULL,

	UNIQUE(product_id, user_id),
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE images (
	id SERIAL PRIMARY KEY,
	image TEXT UNIQUE NOT NULL,
	low_res_image TEXT UNIQUE NOT NULL,

	product_id INT NOT NULL,
	FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE carts (
	id SERIAL PRIMARY KEY,
	total_price INT NOT NULL,
	quantity INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	user_id INT NOT NULL UNIQUE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE cart_items (
	id SERIAL PRIMARY KEY,
	quantity INT NOT NULL,
	total_price INT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	cart_id INT NOT NULL,
	product_id INT NOT NULL UNIQUE,

	FOREIGN KEY(cart_id) REFERENCES carts(id) ON DELETE CASCADE,
	FOREIGN KEY(product_id) REFERENCES product_variants(id) ON DELETE CASCADE
);

CREATE TABLE wishlists (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	user_id INT NOT NULL,
	product_id INT NOT NULL UNIQUE,

	UNIQUE(user_id, product_id),
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TABLE cities (
	id SERIAL PRIMARY KEY,
	city VARCHAR UNIQUE
);

CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
	town VARCHAR NOT NULL,
	street VARCHAR NOT NULL,
	address VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	phone_number VARCHAR NOT NULL,
	total_price INTEGER NOT NULL,
	order_status order_status,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	cancelled_at TIMESTAMP,
	
	cities_id INTEGER NOT NULL,
	user_id INTEGER,

	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
	FOREIGN KEY (cities_id) REFERENCES cities(id)
);

CREATE TABLE order_details (
	id SERIAL PRIMARY KEY,
	quantity INT NOT NULL,
	total_price INT NOT NULL,
	
	product_id INT NOT NULL,
	order_id INT NOT NULL,

	UNIQUE(product_id, order_id),
	FOREIGN KEY (product_id) REFERENCES product_variants(id) ON DELETE SET NULL, -- null or cascade
	FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL -- null or cascade
);

CREATE TABLE discounts (
	id SERIAL PRIMARY KEY,
	discount_type VARCHAR NOT NULL,
	discount_value DECIMAL(10, 2) NOT NULL,
	start_date date NOT NULL,
	end_date date NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE variant_discount (
	id SERIAL PRIMARY KEY,
	
	discount_id INT NOT NULL,
	variant_id INT NOT NULL,

	UNIQUE(discount_id, variant_id),
	FOREIGN KEY (discount_id) REFERENCES discounts(id),
	FOREIGN KEY (variant_id) REFERENCES product_variants(id)
);
-- +goose Down
DROP TABLE cart_items;
DROP TABLE carts;
DROP TABLE order_details;
DROP TABLE orders;
DROP TABLE wishlists;
DROP TABLE rating_review;
DROP TABLE variant_discount;
DROP TABLE cities;
DROP TABLE discounts;
DROP TABLE product_variants;
DROP TABLE images;
DROP TABLE products;
DROP TABLE sizes;
DROP TABLE colors;
DROP TABLE brands;
DROP TABLE categories;
DROP TABLE phone_verification_codes;
DROP TABLE sessions;
DROP TABLE oauth;
DROP TABLE password_resets;
DROP TABLE local_auth;
DROP TABLE users;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS role;
