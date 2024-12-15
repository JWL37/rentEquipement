CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(50) NOT NULL UNIQUE, 
    passwordhash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL 
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY, 
    userID INT NOT NULL, 
    email VARCHAR(255) NOT NULL UNIQUE, 
    phoneNumber VARCHAR(15), 
    address VARCHAR(255), 
    FOREIGN KEY (userID) REFERENCES users (id)
);

CREATE TABLE equipments (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT, 
    pricePerDay DECIMAL(10, 2) NOT NULL,
    isAvailable BOOLEAN DEFAULT TRUE 
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY, 
    customerID INT NOT NULL, 
    equipmentID INT NOT NULL, 
    start_date DATE NOT NULL, 
    endDate DATE NOT NULL, 
    totalCost DECIMAL(10, 2) NOT NULL, 
    FOREIGN KEY (customerID) REFERENCES customers (id),
    FOREIGN KEY (equipmentID) REFERENCES equipments (id)
);


CREATE TABLE maintenances (
    id SERIAL PRIMARY KEY, 
    equipmentID INT NOT NULL, 
    date DATE NOT NULL, 
    description TEXT, 
    FOREIGN KEY (equipmentID) REFERENCES equipments (id)
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY, 
    customerID INT NOT NULL, 
    equipmentID INT NOT NULL, 
    rating INT CHECK (rating BETWEEN 1 AND 5), 
    comment TEXT, 
    reviewDate DATE NOT NULL, 
    FOREIGN KEY (customerID) REFERENCES customers (id),
    FOREIGN KEY (equipmentID) REFERENCES equipments (id)
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    event_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    event_type VARCHAR(50),
    user_id INT,
    message TEXT
);

CREATE OR REPLACE FUNCTION log_user_creation() 
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (event_type, user_id, message)
    VALUES ('User Creation', NEW.id, 'User created with username: ' || NEW.username);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_insert
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION log_user_creation();

CREATE OR REPLACE FUNCTION log_order_creation()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (event_type, user_id, message)
    VALUES ('Order Created', NEW.customerid, 'New order created with equipment ID: ' || NEW.equipmentID || ', start date: ' || NEW.start_date || ', end date: ' || NEW.endDate);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_order_insert
AFTER INSERT ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_creation();

CREATE VIEW orders_view AS
SELECT
    o.id AS order_id,
    o.customerID,
    o.equipmentID,
    e.name AS equipment_name,
    o.start_date,
    o.endDate,
    o.totalCost
FROM
    orders o
INNER JOIN
    equipments e
    ON o.equipmentID = e.id;


