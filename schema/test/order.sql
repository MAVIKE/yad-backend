INSERT INTO orders (user_id, restaurant_id, courier_id, delivery_price, total_price, status)
VALUES (1, 1, 1, 0, 0, 0);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (1, 1, 5);

INSERT INTO orders (user_id, restaurant_id, courier_id, delivery_price, total_price, status)
VALUES (1, 2, 1, 0, 0, 0);

INSERT INTO order_items (order_id, menu_item_id, count)
VALUES (2, 1, 3);