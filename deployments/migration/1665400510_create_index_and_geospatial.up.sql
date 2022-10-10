CREATE INDEX gym_low_idx ON gyms USING btree
(lower((name)::text));

CREATE EXTENSION earthdistance CASCADE;