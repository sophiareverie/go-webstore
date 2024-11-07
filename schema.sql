use mulholland;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS product;

CREATE TABLE customer (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100)
);
CREATE TABLE product (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(100),
    image_name VARCHAR(255),
    price DECIMAL(6,2),
    in_stock INT
);

CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_id INT,
    customer_id INT,
    quantity INT,
    price DECIMAL(6,2),
    tax DECIMAL(6,2),
    donation DECIMAL(6,2),
    timestamp BIGINT,
    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customer(id) ON DELETE CASCADE
);

INSERT INTO customer (first_name, last_name, email) VALUES ('Jason', 'Derulo', 'jd@mines.edu');
INSERT INTO customer (first_name, last_name, email) VALUES ('Donald', 'Trump', 'dt@mines.edu');

INSERT INTO product (product_name, image_name, price, in_stock) VALUES ('Fork', 'assets/images/fork.jpeg', 1.50, 24);
INSERT INTO product (product_name, image_name, price, in_stock) VALUES ('Spoon', 'assets/images/spoon.jpeg', 1.00, 3);
INSERT INTO product (product_name, image_name, price, in_stock) VALUES ('Knife', 'assets/images/knife.jpeg', 2.00, 10);