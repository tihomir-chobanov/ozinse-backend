CREATE TABLE user_favorites(
    "user_id" INT NOT NULL,
    "project_id" INT NOT NULL,
    PRIMARY KEY (user_id, project_id)
);

ALTER TABLE user_favorites ADD CONSTRAINT user_favorites_user_id_foreign FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_favorites ADD CONSTRAINT user_favorites_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;