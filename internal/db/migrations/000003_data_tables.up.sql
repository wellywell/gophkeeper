BEGIN;


CREATE TYPE item_type AS ENUM ('text', 'logopass', 'credit_card', 'binary');

CREATE TABLE item (id BIGSERIAL PRIMARY KEY, user_id BIGINT, key VARCHAR(255), info TEXT, "item_type" item_type,
    CONSTRAINT fk_user_id
    FOREIGN KEY(user_id) 
    REFERENCES auth_user(id)
    ON DELETE NO ACTION);

CREATE UNIQUE INDEX user_key_indx ON item(user_id, key);

CREATE TABLE credit_card (id SERIAL PRIMARY KEY, item_id BIGINT, number VARCHAR(255), owner_name VARCHAR(255), valid_till DATE, cvc int,
    CONSTRAINT fk_card_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE CASCADE);


CREATE INDEX card_item_idx ON credit_card(item_id);

CREATE TABLE text_data(id SERIAL PRIMARY KEY, item_id BIGINT, data text,
    CONSTRAINT fk_text_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE CASCADE);

CREATE INDEX text_item_idx ON text_data(item_id);

CREATE TABLE binary_data (id SERIAL PRIMARY KEY, item_id BIGINT, data bytea,
    CONSTRAINT fk_binary_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE CASCADE);

CREATE INDEX binary_item_idx ON binary_data(item_id);


CREATE TABLE logopass (id SERIAL PRIMARY KEY, item_id BIGINT, login VARCHAR(255), password VARCHAR(255),
    CONSTRAINT fk_logopas_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE CASCADE);

CREATE INDEX logopas_item_idx ON logopass(item_id);


COMMIT;