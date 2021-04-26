-- user1 orders
-- Order 1
INSERT INTO orders (user_id, restaurant_id, courier_id, delivery_price, total_price, status)
VALUES (1, 1, NULL, 100, 900, 0);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (1, 1, 2);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (1, 2, 3);

-- Order 2
INSERT INTO orders (user_id, restaurant_id, courier_id, delivery_price, total_price, status)
VALUES (1, 2, 1, 200, 650, 5);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (2, 4, 3);

-- user2 orders
-- Order 3
INSERT INTO orders (user_id, restaurant_id, courier_id, delivery_price, total_price, status)
VALUES (2, 2, 2, 100, 800, 1);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (3, 4, 3);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (3, 5, 1);