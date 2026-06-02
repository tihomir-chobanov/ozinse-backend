CREATE TABLE main_page_project(
    "id" SERIAL PRIMARY KEY,
    "project_id" INT NOT NULL UNIQUE,
    "position" INT NOT NULL,
    "icon_url" VARCHAR(255) NOT NULL, 
    CONSTRAINT main_page_project_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE
);

