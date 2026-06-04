package service

import (
	"context"
	"math"
	"slices"

	"go.uber.org/zap"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

const entMaxPoint = 2

type QuizResultService struct {
	questionsRepo repository.QuestionRepository
	answersRepo   repository.AnswerRepository
	quizRepo      repository.QuizRepository
	repo          repository.QuizResultRepository
	logger        *zap.Logger
}

func NewQuizResultService(
	repo repository.QuizResultRepository,
	quizRepo repository.QuizRepository,
	questionRepo repository.QuestionRepository,
	answersRepo repository.AnswerRepository,
	logger *zap.Logger,
) *QuizResultService {
	return &QuizResultService{
		repo:          repo,
		quizRepo:      quizRepo,
		questionsRepo: questionRepo,
		answersRepo:   answersRepo,
		logger:        logger,
	}
}

func (s *QuizResultService) Submit(ctx context.Context, userID int64, req *quizv1.SubmitQuizRequest) (*domain.QuizResult, error) {
	if len(req.GetAnswers()) == 0 {
		return nil, domain.ErrEmptyAnswers
	}

	quiz, err := s.quizRepo.GetById(ctx, req.QuizId)
	if err != nil {
		return nil, err
	}

	questions, err := s.questionsRepo.GetByQuizId(ctx, req.QuizId)
	if err != nil {
		return nil, err
	}

	answers := findSkippedQuestions(req.GetAnswers(), questions)

	questionIDs := make([]int64, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	dbAnswers, err := s.answersRepo.GetByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, err
	}

	questionsWithAnswers := QuestionToQuestionWithCorrectAnswer(questions, dbAnswers)
	questionMap := make(map[int64]domain.QuestionWithCorrectAnswer, len(questionsWithAnswers))
	for _, q := range questionsWithAnswers {
		questionMap[q.ID] = q
	}

	var result *domain.QuizResult
	if quiz.IsEntStandard {
		result = computeResultByEntStandard(questionMap, answers)
	} else {
		result = computeResult(questionMap, answers)
	}

	result.QuizID = req.QuizId
	result.UserID = userID
	result.Answers = buildResultAnswers(questionMap, answers, quiz.IsEntStandard)

	return s.repo.Save(ctx, result)
}

func (s *QuizResultService) GetResults(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error) {
	if userID != 0 {
		results, err := s.repo.GetByQuizAndUser(ctx, quizID, userID)
		if err != nil {
			s.logger.Error("failed to get quiz results",
				zap.Int64("quiz_id", quizID),
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
			return nil, err
		}
		return results, nil
	}

	results, err := s.repo.GetByQuiz(ctx, quizID)
	if err != nil {
		s.logger.Error("failed to get quiz results", zap.Int64("quiz_id", quizID), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func computeResult(questionMap map[int64]domain.QuestionWithCorrectAnswer, answers []*quizv1.AnswerSubmission) *domain.QuizResult {
	questionScore := make(map[int64]int, len(questionMap))
	for qID := range questionMap {
		questionScore[qID] = 0
	}

	unanswered := 0

	for _, a := range answers {
		q, ok := questionMap[a.QuestionId]
		if !ok || len(q.CorrectAnswerIDs) == 0 {
			continue
		}
		if len(a.AnswerIds) == 0 {
			unanswered++
			continue
		}
		correctCount := 0
		incorrectCount := 0
		for _, id := range a.AnswerIds {
			if slices.Contains(q.CorrectAnswerIDs, id) {
				correctCount++
			} else {
				incorrectCount++
			}
		}
		if correctCount > 0 && incorrectCount == 0 {
			questionScore[a.QuestionId] = 1
		} else {
			questionScore[a.QuestionId] = -1
		}
	}

	var correct, incorrect int
	var score float64
	for _, s := range questionScore {
		switch {
		case s > 0:
			correct++
			score += float64(s)
		case s < 0:
			incorrect++
		}
	}

	return &domain.QuizResult{
		TotalQuestionsCount:   len(questionMap),
		CorrectAnswersCount:   correct,
		IncorrectAnswersCount: incorrect,
		UnansweredQuestions:   unanswered,
		Score:                 score,
		MaxScore:              float64(len(questionMap)),
	}
}

func computeResultByEntStandard(questionMap map[int64]domain.QuestionWithCorrectAnswer, answers []*quizv1.AnswerSubmission) *domain.QuizResult {
	questionScore := make(map[int64]int, len(questionMap))
	maxPoints := make(map[int64]int, len(questionMap))

	for qID, q := range questionMap {
		questionScore[qID] = 0
		switch {
		case q.TotalAnswersCount == 0:
		case q.TotalAnswersCount <= 4:
			maxPoints[qID] = 1
		default:
			maxPoints[qID] = entMaxPoint
		}
	}

	unanswered := 0

	for _, a := range answers {
		qID := a.QuestionId
		q, ok := questionMap[qID]
		if !ok || len(q.CorrectAnswerIDs) == 0 {
			continue
		}
		if len(a.AnswerIds) == 0 {
			unanswered++
			continue
		}

		correctCount := 0
		incorrectCount := 0
		for _, id := range a.AnswerIds {
			if slices.Contains(q.CorrectAnswerIDs, id) {
				correctCount++
			} else {
				incorrectCount++
			}
		}

		if q.TotalAnswersCount <= 4 {
			if correctCount > 0 && incorrectCount == 0 {
				questionScore[qID] = 1
			} else {
				questionScore[qID] = -1
			}
			continue
		}

		maxCorrect := len(q.CorrectAnswerIDs)
		if maxCorrect == 0 {
			continue
		}
		userAnswerPoint := int(math.Floor(float64(correctCount) / float64(maxCorrect) * entMaxPoint))
		if maxCorrect == 1 {
			questionScore[qID] = userAnswerPoint - incorrectCount
		} else if userAnswerPoint == 1 && incorrectCount == 1 {
			questionScore[qID] = userAnswerPoint
		} else {
			questionScore[qID] = userAnswerPoint - incorrectCount
		}
	}

	var correct, incorrect int
	var score, maxScore float64
	for qID, s := range questionScore {
		maxScore += float64(maxPoints[qID])
		switch {
		case s > 0:
			correct++
			score += float64(s)
		case s < 0:
			incorrect++
		}
	}

	return &domain.QuizResult{
		TotalQuestionsCount:   len(questionMap),
		CorrectAnswersCount:   correct,
		IncorrectAnswersCount: incorrect,
		UnansweredQuestions:   unanswered,
		Score:                 score,
		MaxScore:              maxScore,
	}
}

func buildResultAnswers(questionMap map[int64]domain.QuestionWithCorrectAnswer, answers []*quizv1.AnswerSubmission, isEntStandard bool) []domain.QuizResultAnswer {
	resultAnswers := make([]domain.QuizResultAnswer, 0, len(answers))
	for _, a := range answers {
		q, ok := questionMap[a.QuestionId]

		var correctAnswerIDs []int64
		var score, maxScore float64

		if ok {
			correctAnswerIDs = q.CorrectAnswerIDs
			score, maxScore = computeAnswerScore(q, a.AnswerIds, isEntStandard)
		}

		selectedAnswerIDs := a.AnswerIds
		if selectedAnswerIDs == nil {
			selectedAnswerIDs = []int64{}
		}
		if correctAnswerIDs == nil {
			correctAnswerIDs = []int64{}
		}

		resultAnswers = append(resultAnswers, domain.QuizResultAnswer{
			QuestionID:        a.QuestionId,
			SelectedAnswerIDs: selectedAnswerIDs,
			CorrectAnswerIDs:  correctAnswerIDs,
			Score:             score,
			MaxScore:          maxScore,
		})
	}
	return resultAnswers
}

func computeAnswerScore(q domain.QuestionWithCorrectAnswer, selectedIDs []int64, isEntStandard bool) (score, maxScore float64) {
	if !isEntStandard || q.TotalAnswersCount <= 4 {
		maxScore = 1
		if len(selectedIDs) == 0 {
			return 0, maxScore
		}
		correctCount, incorrectCount := countCorrectIncorrect(q.CorrectAnswerIDs, selectedIDs)
		if correctCount > 0 && incorrectCount == 0 {
			score = 1
		} else {
			score = -1
		}
		return
	}

	maxScore = entMaxPoint
	if len(selectedIDs) == 0 {
		return 0, maxScore
	}

	maxCorrect := len(q.CorrectAnswerIDs)
	if maxCorrect == 0 {
		return 0, maxScore
	}

	correctCount, incorrectCount := countCorrectIncorrect(q.CorrectAnswerIDs, selectedIDs)
	userAnswerPoint := int(math.Floor(float64(correctCount) / float64(maxCorrect) * entMaxPoint))

	var rawScore int
	if maxCorrect == 1 {
		rawScore = userAnswerPoint - incorrectCount
	} else if userAnswerPoint == 1 && incorrectCount == 1 {
		rawScore = userAnswerPoint
	} else {
		rawScore = userAnswerPoint - incorrectCount
	}
	score = float64(rawScore)
	return
}

func countCorrectIncorrect(correctIDs, selectedIDs []int64) (correct, incorrect int) {
	for _, id := range selectedIDs {
		if slices.Contains(correctIDs, id) {
			correct++
		} else {
			incorrect++
		}
	}
	return
}

func findSkippedQuestions(answers []*quizv1.AnswerSubmission, questions []domain.Question) []*quizv1.AnswerSubmission {
	answeredQuestions := make(map[int64]struct{}, len(answers))
	for _, a := range answers {
		answeredQuestions[a.QuestionId] = struct{}{}
	}

	for _, q := range questions {
		if _, found := answeredQuestions[q.ID]; !found {
			answers = append(answers, &quizv1.AnswerSubmission{
				QuestionId: q.ID,
				AnswerIds:  nil,
			})
		}
	}
	return answers
}

func QuestionToQuestionWithCorrectAnswer(questions []domain.Question, answers []domain.Answer) []domain.QuestionWithCorrectAnswer {
	answersByQuestion := make(map[int64][]domain.Answer, len(questions))
	for _, a := range answers {
		answersByQuestion[a.QuestionID] = append(answersByQuestion[a.QuestionID], a)
	}

	results := make([]domain.QuestionWithCorrectAnswer, 0, len(questions))
	for _, q := range questions {
		qAnswers := answersByQuestion[q.ID]
		var correctIDs []int64
		for _, a := range qAnswers {
			if a.Correct {
				correctIDs = append(correctIDs, a.ID)
			}
		}
		results = append(results, domain.QuestionWithCorrectAnswer{
			ID:                q.ID,
			CorrectAnswerIDs:  correctIDs,
			TotalAnswersCount: int64(len(qAnswers)),
		})
	}
	return results
}
