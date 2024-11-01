package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/techschool/simplebank/api"
	"github.com/techschool/simplebank/dbsq"
	_ "github.com/techschool/simplebank/doc/statik"
	"github.com/techschool/simplebank/gapi"
	"github.com/techschool/simplebank/mail"
	"github.com/techschool/simplebank/pb"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	// Load configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Context and cancellation on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// Initialize connection pool
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := dbsq.NewStore(connPool)
	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// Initialize wait group for goroutine management
	waitGroup, ctx := errgroup.WithContext(ctx)

	// Run services concurrently
	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)

	// Wait for all goroutines to complete
	if err := waitGroup.Wait(); err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config util.Config, redisOpt asynq.RedisClientOpt, store dbsq.Store) {
	mailSender := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailSender)

	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown of task processor")
		taskProcessor.Shutdown()
		log.Info().Msg("task processor stopped")
		return nil
	})
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migration up")
	}
	log.Info().Msg("database migrated successfully")
}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store dbsq.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener for gRPC server")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
		if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown of gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server stopped")
		return nil
	})
}

func runGatewayServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store dbsq.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gateway server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{UseProtoNames: true},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	if err := pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik filesystem")
	}
	mux.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(statikFS)))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},         // 允许来自前端应用的跨域请求
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // 允许的方法
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // 允许的请求头
		AllowCredentials: true,
	})
	handler := c.Handler(gapi.HttpLogger(mux))

	httpServer := &http.Server{Handler: handler, Addr: config.HTTPServerAddress}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown of HTTP gateway server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}
		log.Info().Msg("HTTP gateway server stopped")
		return nil
	})
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store dbsq.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create Gin server")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start Gin server at %s", config.HTTPServerAddress)
		if err := server.Start(config.HTTPServerAddress); err != nil {
			log.Error().Err(err).Msg("Gin server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown of Gin server")
		// Shutdown logic here if implemented
		log.Info().Msg("Gin server stopped")
		return nil
	})
}
