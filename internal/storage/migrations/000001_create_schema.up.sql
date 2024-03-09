CREATE TABLE IF NOT EXISTS `merchants` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `api_token` VARCHAR(100) NOT NULL,
  `country` VARCHAR(100) NOT NULL,
  `address` VARCHAR(100) NOT NULL,
  `phone_number` VARCHAR(100) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS `customers` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS `payments` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `order_token` VARCHAR(100) NOT NULL,
  `customer_id` INT NOT NULL,
  `merchant_id` INT NOT NULL,
  `amount` DECIMAL(10, 2) NOT NULL,
  `status` ENUM('pending', 'succeeded', 'failed', 'cancelled', 'refunded', 'processed', 'authorized') NOT NULL DEFAULT 'pending',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(id),
  FOREIGN KEY (merchant_id) REFERENCES merchants(id)
);

CREATE TABLE IF NOT EXISTS `refunds` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `payment_id` INT NOT NULL,
  `amount` DECIMAL(10, 2) NOT NULL,
  `reason` VARCHAR(200),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP,
  FOREIGN KEY (payment_id) REFERENCES payments(id)
);

CREATE TABLE IF NOT EXISTS `credit_cards` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `token` VARCHAR(255) NOT NULL,
  `expiration_month` VARCHAR(20) NOT NULL,
  `expiration_year` VARCHAR(20) NOT NULL,
  `card_holder` VARCHAR(20) NOT NULL,
  `last_four` VARCHAR(20) NOT NULL,
  `card_type` VARCHAR(50) NOT NULL,
  `card_brand` VARCHAR(100) NOT NULL,
  `customer_id` INT NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);