package app

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/kostikan/bd_kursovaya/internal/app/api/forum"
	"github.com/kostikan/bd_kursovaya/internal/pkg/facade"
	"github.com/kostikan/bd_kursovaya/internal/pkg/repo"
	"github.com/kostikan/bd_kursovaya/internal/pkg/sql"
	"github.com/kostikan/bd_kursovaya/internal/pkg/updater"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type appContext struct {
	dig.In
	Forum     *forum.Implementation
	Resharder *updater.Resharder
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

func databaseProvider(ctx context.Context) (*sql.Balancer, error) {
	master, err := sqlx.Connect("postgres", "postgres://bd_kursovaya:bd_kursovaya@localhost:5432/bd_kursovaya?sslmode=disable&binary_parameters=yes")
	if err != nil {
		log.Fatalln(err, "failed connect to database")
	}

	//slave, err := sqlx.Connect("postgres", "postgres://master:12345@localhost:5441/forum?sslmode=disable&binary_parameters=yes")
	//if err != nil {
	//	log.Fatalln(err, "failed connect to database")
	//}

	balancer := sql.New()
	balancer.AddNode(sql.Write, master)
	//balancer.AddNode(sql.Read, slave)

	return balancer, nil
}

func resharderProvider(ctx context.Context, repo *repo.ResharderRepo) (*updater.Resharder, error) {
	resharder := updater.NewResharder(repo)
	return resharder, nil
}

func txManagerProvider(ctx context.Context, db *sql.Balancer) (*sql.TxManager, error) {
	txManager := sql.NewTxManager(db)
	return txManager, nil
}

type facadeProviderOpts struct {
	dig.In
	*repo.AccountRepo
	*repo.CommentRepo
	*repo.PostRepo
	*repo.TagRepo

	*sql.TxManager
}

func facadeProvider(opts facadeProviderOpts) *facade.Facade {
	return facade.NewFacade(facade.Opts{
		AccountRepo: opts.AccountRepo,
		CommentRepo: opts.CommentRepo,
		PostRepo:    opts.PostRepo,
		TagRepo:     opts.TagRepo,
		TxManager:   opts.TxManager,
	})
}

func accountRepoProvider(db *sql.Balancer) *repo.AccountRepo {
	return repo.NewAccountRepo(db)
}

func commentRepoProvider(db *sql.Balancer) *repo.CommentRepo {
	return repo.NewCommentRepo(db)
}

func postRepoProvider(db *sql.Balancer) *repo.PostRepo {
	return repo.NewPostRepo(db)
}

func tagRepoProvider(db *sql.Balancer) *repo.TagRepo {
	return repo.NewTagRepo(db)
}

func reshardeRepoProvider(db *sql.Balancer) *repo.ResharderRepo {
	return repo.NewResharderRepo(db)
}
