CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    latitude DOUBLE PRECISION NOT NULL CHECK (ABS(latitude) <= 90),
    longitude DOUBLE PRECISION NOT NULL CHECK (ABS(latitude) <= 180)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    address_id INT REFERENCES locations (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS couriers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    address_id INT REFERENCES locations (id) ON DELETE CASCADE NOT NULL,
    working_status INT NOT NULL
);

CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(50) NOT NULL,
    address_id INT REFERENCES locations (id) ON DELETE CASCADE NOT NULL,
    working_status INT NOT NULL,
    image VARCHAR(100) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INT REFERENCES restaurants (id) ON DELETE CASCADE NOT NULL,
    title VARCHAR(50) NOT NULL,
    image VARCHAR(100) NOT NULL DEFAULT '',
    description TEXT,
    price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    restaurant_id INT REFERENCES restaurants (id) ON DELETE CASCADE NOT NULL,
    title VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS category_items (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES categories (id) ON DELETE CASCADE NOT NULL,
    menu_item_id INT REFERENCES menu_items (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    password_hash VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    restaurant_id INT REFERENCES restaurants (id) ON DELETE CASCADE NOT NULL,
    courier_id INT REFERENCES couriers (id) ON DELETE CASCADE,
    delivery_price INT NOT NULL DEFAULT 0 CHECK (delivery_price >= 0),
    total_price INT NOT NULL DEFAULT 0 CHECK (total_price >= 0),
    status INT NOT NULL,
    paid TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders (id) ON DELETE CASCADE NOT NULL,
    menu_item_id INT REFERENCES menu_items (id) ON DELETE CASCADE NOT NULL,
    count INT NULL DEFAULT 1 CHECK (count > 0 AND count < 100)
);

CREATE OR REPLACE FUNCTION get_distance(lat1 float, lon1 float, lat2 float, lon2 float)
RETURNS float AS $$
	SELECT 2*6371*asin(sqrt(power(sin(radians((lat2 - lat1)/2)), 2) + 
		cos(radians(lat1))*cos(radians(lat2))*power(sin(radians((lon2 - lon1)/2)), 2)))
$$ LANGUAGE SQL;