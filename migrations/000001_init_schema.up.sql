-- ============================================================================
-- 1. TYPES AND INDEPENDENT TABLES (Lookup tables with no foreign dependencies)
-- ============================================================================

-- Custom ENUM type for defining media project classifications
CREATE TYPE project_type AS ENUM ('movie', 'series');

-- Table storing system and management user roles
CREATE TABLE role (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE
);

-- Table managing target age categories and their visual badge references
CREATE TABLE age_category (  
    "id" SERIAL PRIMARY KEY,
    "range" VARCHAR(255) NOT NULL UNIQUE,
    "icon_url" VARCHAR(255) NOT NULL
);

-- Table managing content sections (e.g., Short Films, Cartoons, Documentaries)
CREATE TABLE category (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE
);

-- Table managing media taxonomy genres (e.g., Action, Comedy, Drama)
CREATE TABLE genre (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "icon_url" VARCHAR(255) NOT NULL
);

-- ============================================================================
-- 2. CORE SYSTEM TABLES (Primary data storage models)
-- ============================================================================

-- Table storing user accounts, authentication profiles, and settings configurations
CREATE TABLE users (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(255) NOT NULL UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "full_name" VARCHAR(255),
    "phone" VARCHAR(255),
    "birth_date" DATE NOT NULL,
    "role_id" INT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "image" VARCHAR(255) NOT NULL,
    "reset_token" VARCHAR(255),
    "reset_token_expires_at" TIMESTAMP,
    "language" VARCHAR(50) DEFAULT 'kk',
    "notifications_enabled" BOOLEAN DEFAULT true,
    "dark_mode_enabled" BOOLEAN DEFAULT false
);

-- Table storing high-level descriptive metadata for Movies and TV Series
CREATE TABLE project (         
    "id" SERIAL PRIMARY KEY,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "release_year" INTEGER NOT NULL,
    "cover_image_url" VARCHAR(255) NOT NULL,
    "is_favorite" BOOLEAN NOT NULL DEFAULT false, 
    "is_featured" BOOLEAN DEFAULT false,
    "type" project_type NOT NULL,
    "duration" INT NOT NULL,
    "keywords" VARCHAR(255) NOT NULL,
    "director" VARCHAR(255) NOT NULL,
    "producer" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- 3. DEPENDENT ENTITIES (Sub-structures tightly bound to core tables)
-- ============================================================================

-- Table managing fine-grained permission controls mapped to specific access modules
CREATE TABLE role_permission (
    "id" SERIAL PRIMARY KEY,
    "role_id" INT NOT NULL,
    "module" VARCHAR(100) NOT NULL,
    "access_level" VARCHAR(50) NOT NULL,
    CONSTRAINT unique_role_module UNIQUE (role_id, module)
);

-- Table structure dividing multi-episode serialization projects into seasons
CREATE TABLE season (
    "id" SERIAL PRIMARY KEY,
    "project_id" INTEGER NOT NULL,
    "season_number" INTEGER NOT NULL,
    CONSTRAINT season_project_number_unique UNIQUE (project_id, season_number)
);

-- Table tracking individual playback details and streaming video references
CREATE TABLE episode (
    "id" SERIAL PRIMARY KEY,
    "season_id" INTEGER NOT NULL,
    "episode_number" INTEGER NOT NULL,
    "youtube_video_id" VARCHAR(255) NOT NULL,
    "duration" INTEGER NOT NULL,
    CONSTRAINT episode_season_number_unique UNIQUE (season_id, episode_number)
);

-- Table capturing multiple decorative promotional image assets per project
CREATE TABLE project_screenshot (
    "id" SERIAL PRIMARY KEY,
    "project_id" INT NOT NULL,
    "url_to_image" VARCHAR(255) NOT NULL
);

-- Table identifying layout indexing positions for showcased main landing page carousels
CREATE TABLE main_page_project(
    "id" SERIAL PRIMARY KEY,
    "project_id" INT NOT NULL UNIQUE,
    "position" INT NOT NULL,
    "icon_url" VARCHAR(255) NOT NULL
);

-- ============================================================================
-- 4. RELATIONSHIP JUNCTION TABLES (Handling many-to-many associations)
-- ============================================================================

-- Intermediary junction mapping multiple genres to media projects
CREATE TABLE project_genre (
    "project_id" INT NOT NULL,
    "genre_id" INT NOT NULL
);

-- Intermediary junction mapping multiple age restrictions to media projects
CREATE TABLE project_age_category (
    "project_id" INT NOT NULL,
    "age_category_id" INT NOT NULL
);

-- Intermediary junction mapping structural landing categories to media projects
CREATE TABLE project_category(
    "project_id" INT NOT NULL,
    "category_id" INT NOT NULL
);

-- Intermediary junction tracking individual user watchlists / bookmarks
CREATE TABLE user_favorites(
    "user_id" INT NOT NULL,
    "project_id" INT NOT NULL,
    PRIMARY KEY (user_id, project_id)
);

-- ============================================================================
-- 5. FOREIGN KEY CONSTRAINTS AND CASCADING RULES (Data integrity boundaries)
-- ============================================================================

-- Role and permission isolation keys
ALTER TABLE users ADD CONSTRAINT users_role_id_foreign FOREIGN KEY(role_id) REFERENCES role(id);
ALTER TABLE role_permission ADD CONSTRAINT role_permission_role_id_foreign FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE CASCADE;

-- Cascading associations for project taxonomies
ALTER TABLE project_genre ADD CONSTRAINT project_genre_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_genre ADD CONSTRAINT project_genre_genre_id_foreign FOREIGN KEY(genre_id) REFERENCES genre(id) ON DELETE CASCADE;

ALTER TABLE project_age_category ADD CONSTRAINT project_age_category_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_age_category ADD CONSTRAINT project_age_category_age_category_id_foreign FOREIGN KEY(age_category_id) REFERENCES age_category(id) ON DELETE CASCADE;

ALTER TABLE project_category ADD CONSTRAINT project_category_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE project_category ADD CONSTRAINT project_category_category_id_foreign FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE;

-- Cascading content mappings for media playback nodes
ALTER TABLE season ADD CONSTRAINT season_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE episode ADD CONSTRAINT episode_season_id_foreign FOREIGN KEY(season_id) REFERENCES season(id) ON DELETE CASCADE;

-- Cascading display components and dashboard listings
ALTER TABLE project_screenshot ADD CONSTRAINT project_screenshot_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;
ALTER TABLE main_page_project ADD CONSTRAINT main_page_project_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;

-- Cascading configurations for customized account profiles
ALTER TABLE user_favorites ADD CONSTRAINT user_favorites_user_id_foreign FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_favorites ADD CONSTRAINT user_favorites_project_id_foreign FOREIGN KEY(project_id) REFERENCES project(id) ON DELETE CASCADE;