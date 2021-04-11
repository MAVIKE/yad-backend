INSERT INTO locations (latitude, longitude) VALUES (30, 50);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant1', '78881007869', 'password', 6, 1, 'img/image1.jpg');

INSERT INTO locations (latitude, longitude) VALUES (35, 50);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant2', '77777777777', 'password', 7, 1, 'img/image2.jpg');

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (1, 'Big Tasty', 'Taratatata', 'Mmmmmmm', 1);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (1, 'Hamburger', 'image2', 'description2', 1);

INSERT INTO menu_items (restaurant_id, title, image, description, price)
VALUES (2, 'Margarita pizza', 'image3', 'description3', 1);