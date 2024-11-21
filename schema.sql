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

CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    password VARCHAR(50),
    email VARCHAR(100),    
    role INT

);
CREATE TABLE product (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(100),
    image_name VARCHAR(255),
    price DECIMAL(6,2),
    in_stock INT,
    inactive TINYINT DEFAULT 0 
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

INSERT INTO users (first_name, last_name, password, email, role) VALUES ('Frodo', 'Baggins','fb', 'fb@mines.edu',1);
INSERT INTO users (first_name, last_name, password, email, role) VALUES ('Harry', 'Potter','hp', 'hp@mines.edu',2);

INSERT INTO product (product_name, image_name, price, in_stock, inactive) VALUES ('Fork', 'assets/images/fork.jpeg', 1.50, 24, 0);
INSERT INTO product (product_name, image_name, price, in_stock, inactive) VALUES ('Spoon', 'assets/images/spoon.jpeg', 1.00, 3, 0);
INSERT INTO product (product_name, image_name, price, in_stock, inactive) VALUES ('Knife', 'assets/images/knife.jpeg', 2.00, 0, 1);
