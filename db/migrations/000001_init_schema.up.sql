-- First, create customers
CREATE TABLE customers (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    phone       VARCHAR(20) NOT NULL,
    address     TEXT
);

-- Then, create users
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role         VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Now, create orders since customers exists
CREATE TABLE orders (
    id              SERIAL PRIMARY KEY,
    order_type      VARCHAR(20) CHECK (order_type IN ('Normal', 'Export')),
    customer_id     INT NOT NULL,
    status          VARCHAR(20) CHECK (status IN ('Created', 'Approved', 'Dispatched', 'Out for Delivery', 'Delivered', 'Returned')),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Change products table to use UUID
CREATE TABLE products (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       INT NOT NULL,
    stock       INT NOT NULL CHECK (stock >= 0),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Then order_items, products, stock (since they depend on orders/products)
CREATE TABLE order_items (
    id          SERIAL PRIMARY KEY,
    order_id    INT NOT NULL,
    product_id  INT NOT NULL, -- Fix product_id type to UUID
    quantity    INT NOT NULL CHECK (quantity > 0),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Inserting products data
INSERT INTO products (name, description, price, stock) VALUES
    ('Smart Bin X1', 'An advanced smart waste bin with AI sorting.', 199, 50),
    ('Smart Bin X2', 'A premium smart waste bin with auto-compaction.', 299, 30),
    ('Smart Bin Lite', 'A budget-friendly smart waste bin.', 99, 100);

-- Inserting customers data
INSERT INTO customers (name, email, phone, address) VALUES
('John Doe', 'john.doe@example.com', '+971-50-123-4567', 'Dubai Marina, Dubai, UAE'),
('Jane Smith', 'jane.smith@example.com', '+971-55-987-6543', 'Sheikh Zayed Road, Dubai, UAE'),
('Alice Johnson', 'alice.johnson@example.com', '+971-52-555-1234', 'Al Ain, Abu Dhabi, UAE'),
('Bob Brown', 'bob.brown@example.com', '+971-56-111-2222', 'Al Khobar Street, Sharjah, UAE'),
('Charlie Davis', 'charlie.davis@example.com', '+971-50-444-5555', 'Corniche, Abu Dhabi, UAE');

INSERT INTO users (username, email, role, password) 
VALUES ('admin', 'admin@example.com', 'admin', '$2a$12$4B985qPlHtRs53.8d/WIleFMX.49ooHa2H6/unqWtpnk/fW6XVmUq'),
('manager', 'manager@example.com', 'manager', '$2a$12$4B985qPlHtRs53.8d/WIleFMX.49ooHa2H6/unqWtpnk/fW6XVmUq'),
('staff', 'staff@example.com', 'staff', '$2a$12$4B985qPlHtRs53.8d/WIleFMX.49ooHa2H6/unqWtpnk/fW6XVmUq');
