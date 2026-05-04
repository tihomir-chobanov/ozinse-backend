-- 1. types and parent tables
CREATE TYPE project_type AS ENUM ('movie', 'series');

CREATE TABLE role (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "permissions" JSON NOT NULL
);

CREATE TABLE age_category (  
    "id" SERIAL PRIMARY KEY,
    "range" VARCHAR(255) NOT NULL UNIQUE,
    "icon_url" VARCHAR(255) NOT NULL
);

CREATE TABLE category (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE genre (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "icon_url" VARCHAR(255) NOT NULL
);

CREATE TABLE project (         
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "release_year" INTEGER NOT NULL,
    "cover_image_url" VARCHAR(255) NOT NULL,
    "is_featured" BOOLEAN NOT NULL, 
    "type" project_type NOT NULL,
    "duration" INT NOT NULL,
    "keywords" VARCHAR(255) NOT NULL,
    "director" VARCHAR(255) NOT NULL,
    "producer" VARCHAR(255) NOT NULL
);

-- 2. child tables
CREATE TABLE users (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "full_name" VARCHAR(255) NOT NULL,
    "phone" VARCHAR(255) NOT NULL UNIQUE,
    "birth_date" DATE NOT NULL,
    "role_id" INT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "image" VARCHAR(255) NOT NULL
);

CREATE TABLE season (
    "id" SERIAL PRIMARY KEY,
    "project_id" INTEGER NOT NULL,
    "season_number" INTEGER NOT NULL
);

CREATE TABLE episode (
    "id" SERIAL PRIMARY KEY,
    "season_id" INTEGER NOT NULL,
    "episode_number" INTEGER NOT NULL,
    "youtube_video_id" VARCHAR(255) NOT NULL,
    "duration" INTEGER NOT NULL
);

CREATE TABLE project_screenshot (
    "id" SERIAL PRIMARY KEY,
    "project_id" INT NOT NULL,
    "url_to_image" VARCHAR(255) NOT NULL
);

-- 3. (Many-to-Many)
CREATE TABLE project_genre (
    "project_id" INT NOT NULL,
    "genre_id" INT NOT NULL
);

CREATE TABLE project_age_category (
    "project_id" INT NOT NULL,
    "age_category_id" INT NOT NULL
);

CREATE TABLE project_category(
    "project_id" INT NOT NULL,
    "category_id" INT NOT NULL
);

-- 4. (Foreign Keys)

ALTER TABLE users ADD CONSTRAINT users_role_id_foreign FOREIGN KEY(role_id) REFERENCES role(id);

-- (CASCADE)
ALTER TABLE project_genre ADD CONSTRAINT project_genre_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_genre ADD CONSTRAINT project_genre_genre_id_foreign FOREIGN KEY(genre_id) REFERENCES genre(id) ON DELETE CASCADE;

ALTER TABLE project_age_category ADD CONSTRAINT project_age_category_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_age_category ADD CONSTRAINT project_age_category_age_category_id_foreign FOREIGN KEY(age_category_id) REFERENCES age_category(id) ON DELETE CASCADE;

ALTER TABLE project_category ADD CONSTRAINT project_category_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_category ADD CONSTRAINT project_category_category_id_foreign FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE;

ALTER TABLE season ADD CONSTRAINT season_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE episode ADD CONSTRAINT episode_season_id_foreign FOREIGN KEY(season_id) REFERENCES season(id) ON DELETE CASCADE;

ALTER TABLE project_screenshot ADD CONSTRAINT project_screenshot_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;

-- 5.(Unique constraints)
ALTER TABLE season ADD CONSTRAINT season_project_number_unique UNIQUE (project_id, season_number);
ALTER TABLE episode ADD CONSTRAINT episode_season_number_unique UNIQUE (season_id, episode_number);