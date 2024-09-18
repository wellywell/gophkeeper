BEGIN;

CREATE TABLE item (id BIGSERIAL PRIMARY KEY, user_id BIGINT, key VARCHAR(255), info TEXT,
    CONSTRAINT fk_user_id
    FOREIGN KEY(user_id) 
    REFERENCES auth_user(id)
    ON DELETE NO ACTION);

CREATE UNIQUE INDEX key_indx ON item(key);
CREATE INDEX user_idx ON item(user_id);

CREATE TABLE credit_card (id SERIAL PRIMARY KEY, item_id BIGINT, number VARCHAR(255), owner_name VARCHAR(255), valid_till DATE, cvc int,
    CONSTRAINT fk_card_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE NO ACTION);


CREATE INDEX card_item_idx ON credit_card(item_id);

CREATE TABLE text_data(id SERIAL PRIMARY KEY, item_id BIGINT, data text,
    CONSTRAINT fk_text_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE NO ACTION);

CREATE INDEX text_item_idx ON text_data(item_id);

CREATE TABLE binary_data (id SERIAL PRIMARY KEY, item_id BIGINT, data bytea,
    CONSTRAINT fk_binary_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE NO ACTION);

CREATE INDEX binary_item_idx ON binary_data(item_id);


CREATE TABLE logopass (id SERIAL PRIMARY KEY, item_id BIGINT, login VARCHAR(255), password VARCHAR(255),
    CONSTRAINT fk_logopas_item_id
    FOREIGN KEY(item_id) 
    REFERENCES item(id)
    ON DELETE NO ACTION);

CREATE INDEX logopas_item_idx ON logopass(item_id);


COMMIT;