INSERT INTO locations (latitude, longitude) VALUES (30, 50);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant1', '78881007869', 'password', 1, 1, 'image1');

INSERT INTO locations (latitude, longitude) VALUES (35, 50);

INSERT INTO restaurants (name, phone, password_hash, address_id, working_status, image)
VALUES ('Restaurant2', '77777777777', 'password', 2, 1, 'image2');

INSERT INTO menu_items (restaurant_id, title, image, description)
VALUES (1, 'Big Tasty', 'Taratatata', 'Mmmmmmm');

INSERT INTO menu_items (restaurant_id, title, image, description)
VALUES (1, 'Hamburger', 'image2', 'description2');

INSERT INTO menu_items (restaurant_id, title, image, description)
VALUES (2, 'Margarita pizza', 'image3', 'description3');