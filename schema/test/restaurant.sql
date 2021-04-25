-- Restaurants
INSERT INTO locations (latitude, longitude) VALUES (52, 85);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant1', '71234567891', 'password', 9, 1, 'img/image1.jpg');

INSERT INTO locations (latitude, longitude) VALUES (55, 85);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant2', '71234567892', 'password', 10, 1, 'img/image1.jpg');

INSERT INTO locations (latitude, longitude) VALUES (56, 87);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant2', '71234567893', 'password', 11, 2, 'img/image1.jpg');

-- Menu items for restaurant 1
INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (1, 'Title1', 'img/image1.jpg', 'Descrption1', 100);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (1, 'Title2', 'img/image1.jpg', 'description2', 200);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (1, 'Title3', 'img/image1.jpg', 'description3', 300);

-- Menu items for restaurant 2
INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (2, 'Title4', 'img/image1.jpg', 'Descrption4', 150);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (2, 'Title5', 'img/image1.jpg', 'description5', 250);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (2, 'Title6', 'img/image1.jpg', 'description6', 350);

-- Menu items for restaurant 3
INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (3, 'Title7', 'img/image1.jpg', 'description7', 100);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (3, 'Title8', 'img/image1.jpg', 'description8', 200);