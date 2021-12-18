CREATE FUNCTION partition_type_comment()
	RETURNS SMALLINT
	IMMUTABLE
	LANGUAGE SQL AS
'SELECT 1 :: SMALLINT';

CREATE FUNCTION partition_type_comment_vote()
	RETURNS SMALLINT
	IMMUTABLE
	LANGUAGE SQL AS
'SELECT 2 :: SMALLINT';

CREATE TABLE partitions
(
	id                    BIGSERIAL PRIMARY KEY,
	partition_type        SMALLINT NOT NULL,
	partition_start_index BIGINT   NOT NULL,
	partition_end_index   BIGINT   NOT NULL
);

CREATE TABLE post_comment_agg
(
	id            BIGSERIAL PRIMARY KEY,
	post_id       BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	comment_count BIGINT NOT NULL DEFAULT 0
);