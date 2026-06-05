package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/Talan-Application/quiz-service/internal/config"
	"github.com/Talan-Application/quiz-service/internal/repository/postgres"
	"github.com/Talan-Application/quiz-service/internal/service"
	grpcserver "github.com/Talan-Application/quiz-service/internal/transport/grpc"
	"github.com/Talan-Application/quiz-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	zapLog := logger.New(cfg.App.Env)
	defer zapLog.Sync()

	db, err := postgres.NewConnection(cfg.Database)
	if err != nil {
		zapLog.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	quizRepo := postgres.NewQuizRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	answerRepo := postgres.NewAnswerRepository(db)
	resultRepo := postgres.NewQuizResultRepository(db)

	quizSvc := service.NewQuizService(quizRepo, zapLog)
	questionSvc := service.NewQuestionService(questionRepo, answerRepo, zapLog)
	resultSvc := service.NewQuizResultService(resultRepo, quizRepo, questionRepo, answerRepo, zapLog)

	grpcSrv := grpcserver.NewServer(cfg.GRPC, cfg.JWT.SecretKey, zapLog, quizSvc, questionSvc, resultSvc)

	go func() {
		if err := grpcSrv.Run(); err != nil {
			zapLog.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcSrv.GracefulStop()
	zapLog.Info("server shut down gracefully")
}
