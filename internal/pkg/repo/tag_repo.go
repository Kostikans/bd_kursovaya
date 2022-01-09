package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/pkg/model"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
	"github.com/lib/pq"
)

// TagRepo  - tag repo
type TagRepo struct {
	db *sql.Balancer
}

// NewTagRepo - returns new tagRepo
func NewTagRepo(db *sql.Balancer) *TagRepo {
	return &TagRepo{
		db: db,
	}
}

func (p *TagRepo) CreateTag(ctx context.Context, tag model.Tag) (id uint64, err error) {
	query := `INSERT INTO tag(author_id, name) VALUES($1,$2) RETURNING id`

	err = p.db.Write(ctx).Get(&id, query, tag.AuthorID, tag.Name)
	return
}

func (p *TagRepo) CheckAuthorExist(ctx context.Context, tag model.Tag) (exist bool, err error) {
	query := `SELECT EXISTS(SELECT * FROM account WHERE id = $1)`

	err = p.db.Read(ctx).Get(&exist, query, tag.AuthorID)
	return
}

func (p *TagRepo) BulkUpdatePostTags(ctx context.Context, ids []uint64, postID uint64) (err error) {
	colCount := 2
	valueStrings := make([]string, 0, len(ids))
	valueArgs := make([]interface{}, 0, len(ids)*colCount)
	cols := "(" + strings.Trim(strings.Repeat("?,", colCount), ",") + ")"
	for _, item := range ids {
		valueStrings = append(valueStrings, cols)
		valueArgs = append(valueArgs, item, postID)
	}

	query := `
INSERT INTO post_tag(tag_id,post_id)
VALUES ` + strings.Join(valueStrings, ",") + ` ON CONFLICT (post_id,tag_id) DO NOTHING`

	_, err = p.db.Write(ctx).Exec(sqlx.Rebind(sqlx.DOLLAR, query), valueArgs...)
	return
}

func (p *TagRepo) CheckPostAndTagsExist(ctx context.Context, ids []uint64, postID uint64) (exist bool, err error) {
	tagIds := make([]string, 0, len(ids))
	for _, id := range ids {
		tagIds = append(tagIds, strconv.FormatUint(id, 10))
	}
	tagQuery := "(" + strings.Join(tagIds, ",") + ")"

	query := fmt.Sprintf(`
SELECT (tab.tag_array @> $1)::bool AND
	   COALESCE((SELECT 1 FROM post WHERE id = $2),0)::bool
FROM (SELECT ARRAY_AGG(id) AS tag_array FROM (SELECT id FROM tag WHERE id IN %s) tmp) AS tab`, tagQuery)

	err = p.db.Read(ctx).Get(&exist, query, pq.Array(ids), postID)
	return
}

func (p *TagRepo) UpdatePostTags(ctx context.Context, ids []uint64, postID uint64) (err error) {
	query := `UPDATE post SET tags_id = $1 WHERE id = $2`

	_, err = p.db.Write(ctx).Exec(query, pq.Array(ids), postID)
	return
}

func (p *TagRepo) GetTagsByIDs(ctx context.Context, ids []uint64) (rows []model.Tag, err error) {
	tagIds := make([]string, 0, len(ids))
	for _, id := range ids {
		tagIds = append(tagIds, strconv.FormatUint(id, 10))
	}
	tagQuery := "(" + strings.Join(tagIds, ",") + ")"

	query := fmt.Sprintf(`
	SELECT id,author_id,name FROM tag WHERE id IN %s`, tagQuery)

	err = p.db.Read(ctx).Select(&rows, query)
	return
}
