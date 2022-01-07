package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/app/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/facade"
	"github.com/kostikan/bd_kursovaya/internal/pkg/repo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type appContext struct {
	dig.In
	Forum *forum.Implementation
}

func provideAndInvoke(ctx context.Context, c *dig.Container, providers map[string]interface{}, fn func(ac appContext)) error {
	err := c.Provide(func() context.Context {
		return ctx
	})
	if err != nil {
		return fmt.Errorf("provide application context: %w", err)
	}

	for name, constructor := range providers {
		err := c.Provide(constructor)
		if err != nil {
			logrus.Errorf("provide %v, error %v", name, err)
			return fmt.Errorf("%s %w", name, err)
		}
	}
	return fmt.Errorf("invoke application %w", c.Invoke(fn))
}

type forumProviderOpts struct {
	dig.In
	*facade.Facade
}

func forumServiceProvider(opts forumProviderOpts) *forum.Implementation {
	return forum.NewForum(forum.Opts{
		Facade: opts.Facade,
	})
}

func databaseProvider(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=bd_kursovaya dbname=bd_kursovaya sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	return db, nil
}

type facadeProviderOpts struct {
	dig.In
	*repo.AccountRepo
	*repo.CommentRepo
	*repo.PostRepo
	*repo.TagRepo
}

func facadeProvider(opts facadeProviderOpts) *facade.Facade {
	return facade.NewFacade(facade.Opts{
		AccountRepo: opts.AccountRepo,
		CommentRepo: opts.CommentRepo,
		PostRepo:    opts.PostRepo,
		TagRepo:     opts.TagRepo,
	})
}

func accountRepoProvider(db *sqlx.DB) *repo.AccountRepo {
	return repo.NewAccountRepo(db)
}

func commentRepoProvider(db *sqlx.DB) *repo.CommentRepo {
	return repo.NewCommentRepo(db)
}

func postRepoProvider(db *sqlx.DB) *repo.PostRepo {
	return repo.NewPostRepo(db)
}

func tagRepoProvider(db *sqlx.DB) *repo.TagRepo {
	return repo.NewTagRepo(db)
}
