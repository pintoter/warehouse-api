CREATE TABLE IF NOT EXISTS warehouse (
  id SERIAL PRIMARY KEY,
  name VARCHAR(25) NOT NULL UNIQUE,
  availability BOOLEAN NOT NULL DEFAULT false
);

INSERT INTO warehouse (name, availability) 
VALUES 
  ('Domodedovo', true),
  ('Sharikovo', true),
  ('Molchanovo', false);

CREATE TYPE GOOD_SIZE AS ENUM ('XS', 'S', 'M', 'L', 'XL', 'XXL', 'XXXL');

CREATE TABLE IF NOT EXISTS product (
  id SERIAL PRIMARY KEY,
  name VARCHAR(25) NOT NULL,
  size GOOD_SIZE NOT NULL DEFAULT 'XS',
  code VARCHAR(25) NOT NULL UNIQUE
);

INSERT INTO product (name, size, code) 
VALUES
  ('Lacoste T-Shirt', 'XS', '12345'),
  ('Lacoste T-Shirt', 'S', '12346'),
  ('Lacoste T-Shirt', 'M', '12347'),
  ('Lacoste T-Shirt', 'L', '12348'),
  ('Lacoste T-Shirt', 'XL', '12349'),
  ('Dads pants', 'L', '1337'),
  ('Dads pants', 'XL', '1338'),
  ('Dads pants', 'XXL', '1339'),
  ('Adidas Hoodie', 'L', '10101011'),
  ('Adidas Hoodie', 'XL', '10101012'),
  ('Adidas Hoodie', 'XXL', '10101013'),
  ('Nike Longsleeve', 'S', '1111131231'),
  ('Nike Longsleeve', 'M', '1111131232'),
  ('Nike Longsleeve', 'L', '1111131233'),
  ('Nike Longsleeve', 'XL', '1111131234');

CREATE TABLE IF NOT EXISTS warehouse_product (
  warehouse_id INTEGER REFERENCES warehouse(id),
  product_id INTEGER REFERENCES product(id),
  quantity SMALLINT,
  PRIMARY KEY (warehouse_id, product_id)
);

INSERT INTO warehouse_product (warehouse_id, product_id, quantity)
VALUES
  (1, 1, 3),
  (1, 2, 5),
  (1, 3, 1),
  (1, 4, 2),
  (2, 1, 3),
  (2, 5, 3),
  (2, 6, 5),
  (2, 7, 1),
  (2, 8, 2),
  (3, 1, 3),
  (3, 9, 3),
  (3, 10, 5),
  (3, 11, 1),
  (3, 12, 2);

CREATE TABLE IF NOT EXISTS reservation (
  id SERIAL PRIMARY KEY,
  reservation_id UUID NOT NULL,
  warehouse_id INTEGER REFERENCES warehouse(id),
  product_id INTEGER REFERENCES product(id),
  quantity SMALLINT,
  reserved_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO reservation (reservation_id, warehouse_id, product_id, quantity)
VALUES
  ('422ab5fa-fbf1-461a-99dc-2c6a49c323f1', 1, 1, 3),
  ('422ab5fa-fbf1-461a-99dc-2c6a49c323f1', 2, 1, 2),
  ('965ac486-0451-4e87-be55-2f985cdbf292', 1, 2, 4);
