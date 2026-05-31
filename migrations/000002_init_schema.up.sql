ALTER TABLE role DROP COLUMN IF EXISTS permissions;

CREATE TABLE role_permission (
    "id" SERIAL PRIMARY KEY,
    "role_id" INT NOT NULL,
    "module" VARCHAR(100) NOT NULL,
    "access_level" VARCHAR(50) NOT NULL,
    CONSTRAINT unique_role_module UNIQUE (role_id, module),
    CONSTRAINT role_permission_role_id_foreign FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE
);