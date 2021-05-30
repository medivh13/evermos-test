DROP TABLE IF EXISTS products CASCADE;
CREATE TABLE products (
    p_id serial PRIMARY KEY,
    p_name VARCHAR(50) NOT NULL,
    p_desc text NOT NULL,
    p_stocks integer
);

CREATE INDEX idx_products_1 ON products (p_id, p_name);

INSERT INTO products (p_name, p_desc, p_stocks) VALUES ('PRODUCT A', 'TEST PRODUCT A', 10);
INSERT INTO products (p_name, p_desc, p_stocks) VALUES ('PRODUCT B', 'TEST PRODUCT B', 10);
INSERT INTO products (p_name, p_desc, p_stocks) VALUES ('PRODUCT C', 'TEST PRODUCT C', 10);
INSERT INTO products (p_name, p_desc, p_stocks) VALUES ('PRODUCT D', 'TEST PRODUCT D', 10);

DROP TABLE IF EXISTS customers CASCADE;
CREATE TABLE customers (
    c_id serial PRIMARY KEY,
    c_email VARCHAR(50) NOT NULL,
    c_password text NOT NULL,
    c_tokens text
);

CREATE INDEX idx_customers_1 ON customers (c_id, c_email, c_tokens);

DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders (
    o_id serial PRIMARY KEY,
    o_c_id integer,
    o_status integer,
    o_due_date timestamp(6) without time zone NOT NULL
);
CREATE INDEX idx_orders_1 ON orders (o_id, o_status);

DROP TABLE IF EXISTS order_details CASCADE;
CREATE TABLE order_details (
    od_id serial PRIMARY KEY,
    od_o_id integer,
    od_p_id integer,
    od_qty integer
);
CREATE INDEX idx_order_details_1 ON order_details (od_id, od_o_id);