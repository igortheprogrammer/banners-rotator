BEGIN;

CREATE TABLE slots
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '""',
    CONSTRAINT "slots_pk" PRIMARY KEY (id)
);

CREATE TABLE banners
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '""',
    CONSTRAINT "banners_pk" PRIMARY KEY (id)
);

CREATE TABLE groups
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '""',
    CONSTRAINT "groups_pk" PRIMARY KEY (id)
);

CREATE TABLE rotations
(
    slot_id   bigint NOT NULL,
    banner_id bigint NOT NULL
);

CREATE TABLE views
(
    slot_id   bigint NOT NULL,
    banner_id bigint NOT NULL,
    group_id  bigint NOT NULL,
    date      bigint NOT NULL
);

CREATE TABLE clicks
(
    slot_id   bigint NOT NULL,
    banner_id bigint NOT NULL,
    group_id  bigint NOT NULL,
    date      bigint NOT NULL
);

ALTER TABLE rotations
    ADD CONSTRAINT fk_rotations_slots FOREIGN KEY (slot_id)
        REFERENCES slots (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID;


ALTER TABLE rotations
    ADD CONSTRAINT fk_rotations_banners FOREIGN KEY (banner_id)
        REFERENCES banners (id);


ALTER TABLE views
    ADD CONSTRAINT fk_views_slots FOREIGN KEY (slot_id)
        REFERENCES slots (id);


ALTER TABLE views
    ADD CONSTRAINT fk_views_banners FOREIGN KEY (banner_id)
        REFERENCES banners (id);


ALTER TABLE views
    ADD CONSTRAINT fk_views_groups FOREIGN KEY (group_id)
        REFERENCES groups (id);


ALTER TABLE clicks
    ADD CONSTRAINT fk_clicks_slots FOREIGN KEY (slot_id)
        REFERENCES slots (id);


ALTER TABLE clicks
    ADD CONSTRAINT fk_clicks_banners FOREIGN KEY (banner_id)
        REFERENCES banners (id);


ALTER TABLE clicks
    ADD CONSTRAINT fk_clicks_groups FOREIGN KEY (group_id)
        REFERENCES groups (id);

END;
