BEGIN;

INSERT INTO "user_roles" ("title") VALUES ('employee'), ('admin'), ('superadmin');

INSERT INTO "users" ("email", "username", "password", "role_id")
VALUES 
  ('emp001@mail.com', 'employee001', '$2y$10$SvkeV0XlhJYx5HdYJM96reaU76WfZngkn2iWfdJG94Mdld9F.akVW', '1'),
  ('admin001@mail.com', 'admin001', '$2y$10$SvkeV0XlhJYx5HdYJM96reaU76WfZngkn2iWfdJG94Mdld9F.akVW', '2'),
  ('super@mail.com', 'superadmin001', '$2y$10$aROrDjTllIQYVmCFUEK7F.ubjQA7nzhEzOFeLntdSqK6OiyeuQb6m', '3');

INSERT INTO "categories" ("title", "desc")
VALUES
  ('food & beverage', ''),
  ('fashion', 'all fashion in the world'),
  ('gadget', 'just a gadget');

INSERT INTO "products" ("title", "desc", "price", "discount", "stock", "category_id")
VALUES
    ('Coffee', 'Just a food & beverage product', 150, 30, 999, 1),
    ('Steak', 'Just a food & beverage product', 200, 0, 100, 1),
    ('Shirt', 'Just a fashion product', 590, 90, 20, 2),
    ('Phone', 'Just a gadget product', 33400, 3400, 8, 3),
    ('Computer', 'Just a gadget product', 49000, 500, 3, 3);

COMMIT;