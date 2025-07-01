-- SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
--
-- SPDX-License-Identifier: Apache-2.0
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- SPDX-SnippetBegin
-- SPDX-FileCopyrightText: 2024 Daniel Verite <daniel-at-manitou-mail.org>
-- SPDX-License-Identifier: MIT
-- ============================================================================
-- UUIDv7 Creation
-- ============================================================================
-- Retrieved from:
-- https://github.com/dverite/postgres-uuidv7-sql/blob/main/sql/uuidv7-sql--1.0.sql
/* Main function to generate a uuidv7 value with millisecond precision */
CREATE FUNCTION uuidv7 (timestamptz DEFAULT clock_timestamp()) RETURNS uuid AS $$
  -- Replace the first 48 bits of a uuidv4 with the current
  -- number of milliseconds since 1970-01-01 UTC
  -- and set the "ver" field to 7 by setting additional bits
  select encode(
    set_bit(
      set_bit(
        overlay(uuid_send(gen_random_uuid()) placing
	  substring(int8send((extract(epoch from $1)*1000)::bigint) from 3)
	  from 1 for 6),
	52, 1),
      53, 1), 'hex')::uuid;
$$ LANGUAGE sql volatile parallel safe;

COMMENT ON FUNCTION uuidv7 (timestamptz) IS 'Generate a uuid-v7 value with a 48-bit timestamp (millisecond precision) and 74 bits of randomness';

/* Version with the "rand_a" field containing sub-milliseconds (method 3 of the spec)
clock_timestamp() is hoped to provide enough precision and consecutive
calls to not happen fast enough to output the same values in that field.
The uuid is the concatenation of:
- 6 bytes with the current Unix timestamp (number of milliseconds since 1970-01-01 UTC)
- 2 bytes with
- 4 bits for the "ver" field
- 12 bits for the fractional part after the milliseconds
- 8 bytes of randomness from the second half of a uuidv4
*/
CREATE FUNCTION uuidv7_sub_ms (timestamptz DEFAULT clock_timestamp()) RETURNS uuid AS $$
 select encode(
   substring(int8send(floor(t_ms)::int8) from 3) ||
   int2send((7<<12)::int2 | ((t_ms-floor(t_ms))*4096)::int2) ||
   substring(uuid_send(gen_random_uuid()) from 9 for 8)
  , 'hex')::uuid
  from (select extract(epoch from $1)*1000 as t_ms) s
$$ LANGUAGE sql volatile parallel safe;

COMMENT ON FUNCTION uuidv7_sub_ms (timestamptz) IS 'Generate a uuid-v7 value with a 60-bit timestamp (sub-millisecond precision) and 62 bits of randomness';

/* Extract the timestamp in the first 6 bytes of the uuidv7 value.
Use the fact that 'xHHHHH' (where HHHHH are hexadecimal numbers)
can be cast to bit(N) and then to int8.
*/
CREATE FUNCTION uuidv7_extract_timestamp (uuid) RETURNS timestamptz AS $$
 select to_timestamp(
   right(substring(uuid_send($1) from 1 for 6)::text, -1)::bit(48)::int8 -- milliseconds
    /1000.0);
$$ LANGUAGE sql immutable strict parallel safe;

COMMENT ON FUNCTION uuidv7_extract_timestamp (uuid) IS 'Return the timestamp stored in the first 48 bits of the UUID v7 value';

CREATE FUNCTION uuidv7_boundary (timestamptz) RETURNS uuid AS $$
  /* uuid fields: version=0b0111, variant=0b10 */
  select encode(
    overlay('\x00000000000070008000000000000000'::bytea
      placing substring(int8send(floor(extract(epoch from $1) * 1000)::bigint) from 3)
        from 1 for 6),
    'hex')::uuid;
$$ LANGUAGE sql stable strict parallel safe;

COMMENT ON FUNCTION uuidv7_boundary (timestamptz) IS 'Generate a non-random uuidv7 with the given timestamp (first 48 bits) and all random bits to 0. As the smallest possible uuidv7 for that timestamp, it may be used as a boundary for partitions.';

-- SPDX-SnippetEnd
-- ============================================================================
-- USERS TABLE
-- ============================================================================
-- This table stores persistent user data synchronized with Auth0
-- Auth0 handles authentication, we store business-specific user information
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  auth0_user_id VARCHAR(255) UNIQUE,
  email VARCHAR(255) UNIQUE NOT NULL,
  email_verified BOOLEAN DEFAULT FALSE,
  password_hash TEXT NOT NULL,
  total_purchases_amount INTEGER DEFAULT 0,
  total_purchases_count INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  last_login_at TIMESTAMP WITH TIME ZONE,
  -- Soft delete support (for GDPR compliance and data recovery)
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  CONSTRAINT users_email_format CHECK (
    email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
  )
);

-- ============================================================================
-- PRODUCT CATEGORIES TABLE
-- ============================================================================
-- SPEC: v1.0.0-s3.1.2
CREATE TABLE product_categories (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  name VARCHAR(100) UNIQUE NOT NULL,
  parent_category_id UUID REFERENCES product_categories (id),
  sort_order INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- PRODUCTS TABLE
-- ============================================================================
CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  product_type VARCHAR(20) NOT NULL DEFAULT 'plugin',
  name VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  short_description VARCHAR(500),
  price_id VARCHAR(255),
  price INTEGER NOT NULL,
  product_id VARCHAR(255),
  category_id UUID REFERENCES product_categories (id),
  credits TEXT,
  download_filename VARCHAR(255) NOT NULL,
  download_filesize BIGINT,
  download_checksum VARCHAR(128),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  released_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT products_price_positive CHECK (price >= 0),
  CONSTRAINT products_type_valid CHECK (product_type IN ('plugin', 'merchandise'))
);

-- ============================================================================
-- ORDERS TABLE
-- ============================================================================
-- SPEC: v1.0.0-s3.2.1
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  user_id UUID NOT NULL REFERENCES users (id),
  stripe_payment_intent_id VARCHAR(255) UNIQUE NOT NULL,
  stripe_customer_id VARCHAR(255),
  total_amount INTEGER NOT NULL,
  currency VARCHAR(3) DEFAULT 'USD',
  tax_amount INTEGER DEFAULT 0,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  -- Of the from BDE-YYYY-MM-DD-XXXXXXXXXXXXXXXXXXXXXXXX
  order_number VARCHAR(39) UNIQUE NOT NULL,
  billing_email VARCHAR(255) NOT NULL,
  billing_name VARCHAR(200),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT orders_total_positive CHECK (total_amount > 0),
  CONSTRAINT orders_status_valid CHECK (
    status IN (
      'pending',
      'processing',
      'completed',
      'failed',
      'refunded',
      'cancelled'
    )
  )
);

-- ============================================================================
-- ORDER ITEMS TABLE
-- ============================================================================
-- Individual products within each order (shopping cart items)
CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products (id),
  -- Item details at the time of purchase
  -- These will not change if product information changes
  product_name VARCHAR(200) NOT NULL,
  product_price INTEGER NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 1,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  line_total INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  CONSTRAINT order_items_quantity_positive CHECK (quantity > 0),
  CONSTRAINT order_items_price_positive CHECK (product_price >= 0),
  CONSTRAINT order_items_total_positive CHECK (line_total >= 0),
  CONSTRAINT order_items_status_valid CHECK (
    status IN (
      'pending',
      'processing',
      'completed',
      'failed',
      'refunded',
      'cancelled'
    )
  )
);

-- ============================================================================
-- USER DOWNLOADS TABLE
-- ============================================================================
-- Tracks download history for purchased products (SPEC: v1.0.0-s3.4.2)
CREATE TABLE user_downloads (
  id UUID PRIMARY KEY DEFAULT uuidv7 (),
  user_id UUID NOT NULL REFERENCES users (id),
  product_id UUID NOT NULL REFERENCES products (id),
  order_id UUID NOT NULL REFERENCES orders (id),
  download_count INTEGER DEFAULT 0,
  last_downloaded_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT downloads_count_non_negative CHECK (download_count >= 0)
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================
-- These indexes optimize common query patterns for better performance.
-- User lookup indexes (Auth0 integration and profile management)
CREATE INDEX idx_users_auth0_user_id ON users (auth0_user_id);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_created_at ON users (created_at);

-- Product search and browsing indexes
CREATE INDEX idx_products_category ON products (category_id);

CREATE INDEX idx_products_price ON products (price);

CREATE INDEX idx_products_released_at ON products (released_at DESC);

-- Full-text search index for product names and descriptions
CREATE INDEX idx_products_search ON products USING GIN (
  to_tsvector('english', name || ' ' || description)
);

-- Order and transaction indexes
CREATE INDEX idx_orders_user_id ON orders (user_id);

CREATE INDEX idx_orders_status ON orders (status);

CREATE INDEX idx_orders_created_at ON orders (created_at DESC);

CREATE INDEX idx_orders_stripe_payment_intent ON orders (stripe_payment_intent_id);

-- Order items for order details lookup
CREATE INDEX idx_order_items_order_id ON order_items (order_id);

CREATE INDEX idx_order_items_product_id ON order_items (product_id);

-- Download tracking indexes
CREATE INDEX idx_user_downloads_user_id ON user_downloads (user_id);

CREATE INDEX idx_user_downloads_product_id ON user_downloads (product_id);

CREATE INDEX idx_user_downloads_order_id ON user_downloads (order_id);

-- ============================================================================
-- FUNCTIONS AND TRIGGERS
-- ============================================================================
-- Automated database maintenance and business logic. These can be migrated
-- to the application layer.
-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column () RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply the update trigger to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE
UPDATE ON users FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER update_products_updated_at BEFORE
UPDATE ON products FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER update_product_categories_updated_at BEFORE
UPDATE ON product_categories FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER update_orders_updated_at BEFORE
UPDATE ON orders FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

CREATE TRIGGER update_order_items_updated_at BEFORE
UPDATE ON order_items FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column ();

-- Keep this order number generation in the application layer, lol.
--
-- Function to generate human-readable order numbers
-- CREATE OR REPLACE FUNCTION generate_order_number()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- Generate order number in format: BDE-YYYYMMDD-XXXXXXXX
--     NEW.order_number = 'BDE-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
--                        LPAD(EXTRACT(EPOCH FROM NEW.created_at)::INTEGER % 10000, 4, '0');
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- -- Apply order number generation trigger
-- CREATE TRIGGER generate_order_number_trigger BEFORE INSERT ON orders
--     FOR EACH ROW EXECUTE FUNCTION generate_order_number();
-- Function to update user purchase statistics when orders are completed
CREATE OR REPLACE FUNCTION update_user_purchase_stats () RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status != 'completed') THEN
        UPDATE users 
        SET total_purchases_amount = total_purchases_amount + NEW.total_amount,
            total_purchases_count = total_purchases_count + 1
        WHERE id = NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply purchase statistics trigger
CREATE TRIGGER update_user_purchase_stats_trigger
AFTER
UPDATE ON orders FOR EACH ROW
EXECUTE FUNCTION update_user_purchase_stats ();

-- ============================================================================
-- VIEWS FOR COMMON QUERIES
-- ============================================================================
-- Convenient views that encapsulate common business logic
-- User purchase history with product details
CREATE VIEW user_purchase_history AS
SELECT
  u.id as user_id,
  u.email,
  o.id as order_id,
  o.order_number,
  o.total_amount,
  o.status as order_status,
  o.created_at as order_date,
  oi.product_id,
  oi.product_name,
  oi.product_price,
  oi.quantity,
  oi.line_total
FROM
  users u
  JOIN orders o ON u.id = o.user_id
  JOIN order_items oi ON o.id = oi.order_id
WHERE
  u.deleted_at IS NULL;
