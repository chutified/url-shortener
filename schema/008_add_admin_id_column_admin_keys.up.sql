ALTER TABLE admin_keys
ADD COLUMN admin_id UUID NOT NULL UNIQUE;

ALTER TABLE admin_keys
ADD CONSTRAINT fk_admin_id
FOREIGN KEY (admin_id)
REFERENCES admins (admin_id)
ON DELETE CASCADE;
