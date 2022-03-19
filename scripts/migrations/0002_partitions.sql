CREATE TABLE post_comment_agg
(
	id            BIGSERIAL PRIMARY KEY,
	post_id       BIGINT NOT NULL REFERENCES post (id) ON DELETE RESTRICT,
	comment_count BIGINT NOT NULL DEFAULT 0
);

CREATE OR REPLACE FUNCTION create_new_partition(parent_table_name TEXT,
												first_id BIGINT,
												last_id BIGINT,
												partition_name TEXT) RETURNS VOID AS
$BODY$
DECLARE
	sql TEXT;
BEGIN

	SELECT FORMAT('CREATE TABLE IF NOT EXISTS %s (CHECK (
         post_id > %s AND
         post_id <= %s))
         INHERITS (%I)', partition_name, first_id::TEXT,
				  last_id::TEXT,
				  parent_table_name)
	INTO sql;

	EXECUTE sql;
	PERFORM index_partition(partition_name);
END;
$BODY$
	LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION index_partition(partition_name TEXT) RETURNS VOID AS
$BODY$
BEGIN
	EXECUTE 'CREATE INDEX IF NOT EXISTS ' || partition_name || '_post_id_idx ON ' || partition_name ||
			' USING btree (post_id)';
END;
$BODY$
	LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_comment_row()
	RETURNS TRIGGER AS
$BODY$
DECLARE
	partition_id   BIGINT;
	partition_name TEXT;
	range          BIGINT;
	first_id       BIGINT;
	last_id        BIGINT;
	sql            TEXT;
BEGIN
	partition_id := NEW.post_id;
	range := 10;
	IF partition_id % range = 0 THEN
		first_id := (partition_id - 1) / range * range;
		last_id := ((partition_id - 1) / range + 1) * range;
	ELSE
		first_id := partition_id / range * range;
		last_id := (partition_id/ range + 1) * range;
	END IF;

	range := 10;
	partition_name :=
		FORMAT('%s_%s_%s', TG_TABLE_NAME, first_id::TEXT,
			   last_id::TEXT);

	RAISE NOTICE '% % ', first_id::TEXT,
		last_id::TEXT;

	IF NOT EXISTS(SELECT relname FROM pg_class WHERE relname = partition_name) THEN
		PERFORM pg_advisory_xact_lock(1234);
		PERFORM create_new_partition(TG_TABLE_NAME, first_id, last_id, partition_name);
	END IF;

	SELECT FORMAT('INSERT INTO %s values ($1.*)', partition_name) INTO sql;
	EXECUTE sql USING NEW;

	RETURN NEW;
END;
$BODY$
	LANGUAGE plpgsql;

CREATE TRIGGER comment_insert_trigger
	BEFORE INSERT
	ON comment
	FOR EACH ROW
EXECUTE PROCEDURE insert_comment_row();

CREATE OR REPLACE FUNCTION delete_parent_row_comment()
	RETURNS TRIGGER AS
$BODY$
DECLARE
BEGIN
	DELETE FROM ONLY comment WHERE id = NEW.ID;
	RETURN NULL;
END;
$BODY$
	LANGUAGE plpgsql;

CREATE TRIGGER after_insert_row_trigger_comment
	AFTER INSERT
	ON comment
	FOR EACH ROW
EXECUTE PROCEDURE delete_parent_row_comment();