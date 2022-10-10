CREATE TABLE IF NOT EXISTS gyms (
    id serial primary key,
    name character varying(255),
    gym_type character varying(50),
    city character varying(50),
    state character varying(50),
    country character varying(50),
    address character varying(255),
    lat double precision,
    long double precision,
    amenities character varying(255),
    created_on timestamp with time zone,
    modified_on timestamp with time zone
);
