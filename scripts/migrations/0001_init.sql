CREATE TABLE account
(
	id          BIGSERIAL PRIMARY KEY,
	nickname    VARCHAR(100) NOT NULL,
	avatar      VARCHAR(200) NOT NULL,
	description TEXT         NOT NULL,
	is_deleted  BOOL                     DEFAULT FALSE,
	created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE post
(
	id         BIGSERIAL PRIMARY KEY,
	author_id  BIGINT       NOT NULL REFERENCES account (id) ON DELETE RESTRICT,
	title      VARCHAR(100) NOT NULL,
	text       TEXT         NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	tags_id    BIGINT[]    NOT NULL  DEFAULT '{}'::BIGINT[]
);

CREATE TABLE post_vote_agg
(
	id            BIGSERIAL PRIMARY KEY,
	post_id       BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	like_count    BIGINT NOT NULL DEFAULT 0,
	dislike_count BIGINT NOT NULL DEFAULT 0,
	UNIQUE (post_id)
);

CREATE INDEX ON post_vote_agg USING btree (post_id);

CREATE TABLE post_vote
(
	id        BIGSERIAL PRIMARY KEY,
	post_id   BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	author_id BIGINT NOT NULL REFERENCES account (id) ON DELETE RESTRICT,
	vote      BOOL   NOT NULL,
	UNIQUE (post_id,author_id)
);

CREATE INDEX ON post_vote USING btree (post_id);

CREATE TABLE comment
(
	id         BIGSERIAL PRIMARY KEY,
	author_id  BIGINT       NOT NULL REFERENCES account (id) ON DELETE RESTRICT,
	post_id    BIGINT       NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	parent_id  BIGINT       NOT NULL    DEFAULT 0,
	text       TEXT         NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX ON comment USING btree (post_id);

CREATE TABLE comment_vote_agg
(
	id            BIGSERIAL PRIMARY KEY,
	comment_id    BIGINT NOT NULL REFERENCES comment (id) ON DELETE RESTRICT,
	like_count    BIGINT NOT NULL DEFAULT 0,
	dislike_count BIGINT NOT NULL DEFAULT 0,
	UNIQUE (comment_id)
);

CREATE INDEX ON comment_vote_agg USING btree (comment_id);

CREATE TABLE comment_vote
(
	id         BIGSERIAL PRIMARY KEY,
	post_id    BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	author_id  BIGINT NOT NULL REFERENCES account (id) ON DELETE RESTRICT,
	comment_id BIGINT NOT NULL REFERENCES comment (id) ON DELETE RESTRICT,
	vote       BOOL   NOT NULL,
	UNIQUE (post_id,comment_id,author_id)
);

CREATE INDEX ON comment_vote USING btree (post_id, comment_id);

CREATE TABLE post_tag
(
	id      BIGSERIAL PRIMARY KEY,
	post_id BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	tag_id  BIGINT NOT NULL REFERENCES tag (id) ON DELETE RESTRICT,
	UNIQUE (post_id,tag_id)
);

CREATE INDEX ON post_tag USING btree (post_id);

CREATE TABLE tag
(
	id        BIGSERIAL PRIMARY KEY,
	author_id BIGINT      NOT NULL REFERENCES account (id) ON DELETE RESTRICT,
	name      VARCHAR(50) NOT NULL
);

CREATE FUNCTION trigger_set_timestamp()
	RETURNS TRIGGER AS
$$
BEGIN
	NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_account
	BEFORE UPDATE
	ON account
	FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_post
	BEFORE UPDATE
	ON post
	FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_comment
	BEFORE UPDATE
	ON comment
	FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE comment_1x
(
	CHECK ( post_id >= 0 AND post_id < 10000 )
) INHERITS (comment);

CREATE TABLE comment_vote_1x
(
	CHECK ( post_id >= 0 AND post_id < 10000 )
) INHERITS (comment_vote);