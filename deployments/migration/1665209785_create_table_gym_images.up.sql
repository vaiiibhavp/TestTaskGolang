CREATE TABLE IF NOT EXISTS gym_images (
    id serial primary key,
    gym_id int,
    image_type character varying(50),
    label character varying(50),
    created_on timestamp with time zone,
    modified_on timestamp with time zone
);