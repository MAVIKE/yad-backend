-- Categories
INSERT INTO categories (restaurant_id, title)
VALUES (1, 'Category1');

INSERT INTO categories (restaurant_id, title)
VALUES (1, 'Category2');

INSERT INTO categories (restaurant_id, title)
VALUES (2, 'Category3');

INSERT INTO categories (restaurant_id, title)
VALUES (2, 'Category4');

INSERT INTO categories (restaurant_id, title)
VALUES (3, 'Category5');

-- Category items for menu items in restaurant 1
INSERT INTO category_items(category_id, menu_item_id)
VALUES (1, 1);

INSERT INTO category_items(category_id, menu_item_id)
VALUES (1, 2);

INSERT INTO category_items(category_id, menu_item_id)
VALUES (2, 3);

-- Category items for menu items in restaurant 2
INSERT INTO category_items(category_id, menu_item_id)
VALUES (3, 4);

INSERT INTO category_items(category_id, menu_item_id)
VALUES (4, 5);

INSERT INTO category_items(category_id, menu_item_id)
VALUES (4, 6);

-- Category items for menu items in restaurant 3
INSERT INTO category_items(category_id, menu_item_id)
VALUES (5, 7);

INSERT INTO category_items(category_id, menu_item_id)
VALUES (5, 8);
