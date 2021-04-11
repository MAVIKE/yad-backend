DROP TRIGGER IF EXISTS total_price_update_trigger ON order_items;
DROP FUNCTION IF EXISTS get_total_price(int);
DROP FUNCTION IF EXISTS update_total_price;
DROP FUNCTION IF EXISTS get_distance(float, float, float, float);

DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS admins CASCADE;
DROP TABLE IF EXISTS category_items CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS menu_items CASCADE;
DROP TABLE IF EXISTS restaurants CASCADE;
DROP TABLE IF EXISTS couriers CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
