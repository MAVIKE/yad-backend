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
    working_status INT NOT NULL CHECK (working_status BETWEEN 0 AND 2)
);

CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
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
    count INT NULL DEFAULT 1 CHECK (count > 0 AND count < 100),
    UNIQUE(order_id, menu_item_id)
);

CREATE OR REPLACE FUNCTION get_distance(lat1 float, lon1 float, lat2 float, lon2 float)
RETURNS float AS $$
	SELECT 2*6371*asin(sqrt(power(sin(radians((lat2 - lat1)/2)), 2) + 
		cos(radians(lat1))*cos(radians(lat2))*power(sin(radians((lon2 - lon1)/2)), 2)))
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION get_total_price(cur_order_id int)
RETURNS bigint AS $$
	SELECT SUM(tmp.mul) + (SELECT delivery_price FROM orders WHERE id = cur_order_id) FROM 
	(
		SELECT count * 
			(
				SELECT price FROM menu_items WHERE id = oi.menu_item_id
			) AS mul
		FROM order_items AS oi 
		WHERE order_id = cur_order_id
	) AS tmp
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION update_total_price()
RETURNS trigger AS $$
BEGIN
	IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
		UPDATE orders SET total_price = get_total_price(NEW.order_id) 
			WHERE id = NEW.order_id;
		RETURN NEW;
	ELSIF TG_OP = 'DELETE' THEN
		UPDATE orders SET total_price = get_total_price(OLD.order_id) 
			WHERE id = OLD.order_id;
		RETURN OLD;
	END IF;
RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER total_price_update_trigger
AFTER INSERT OR UPDATE OR DELETE ON order_items 
FOR EACH ROW EXECUTE PROCEDURE update_total_price ();