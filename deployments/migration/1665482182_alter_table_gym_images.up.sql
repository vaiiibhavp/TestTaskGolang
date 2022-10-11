CREATE TYPE type_enum AS ENUM('amenities', 'gym');

ALTER TABLE gym_images
ADD COLUMN type type_enum DEFAULT 'gym';