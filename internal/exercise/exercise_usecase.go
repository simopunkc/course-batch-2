package exercise

import (
	"course/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseUsecase struct {
	db *gorm.DB
}

func NewExerciseUsecase(db *gorm.DB) *ExerciseUsecase {
	return &ExerciseUsecase{db}
}

func (exUsecase ExerciseUsecase) GetExerciseByID(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}

	type Question struct {
		ID            int       `json:"id"`
		ExerciseID    int       `json:"-"`
		Body          string    `json:"body"`
		OptionA       string    `json:"option_a"`
		OptionB       string    `json:"option_b"`
		OptionC       string    `json:"option_c"`
		OptionD       string    `json:"option_d"`
		CorrectAnswer string    `json:"-"`
		Score         int       `json:"score"`
		CreatorID     int       `json:"-"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	type Exercise struct {
		ID          int        `json:"id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Questions   []Question `json:"questions"`
	}

	var exercise Exercise
	err = exUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "exercise not found",
		})
		return
	}

	c.JSON(200, exercise)
}

func (exUsecase ExerciseUsecase) CalculateUserScore(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = exUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "exercise not found",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = exUsecase.db.Where("user_id = ?", userID).Find(&answers).Error
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error when find answers",
		})
		return
	}
	if len(answers) == 0 {
		c.JSON(200, map[string]interface{}{
			"score": 0,
		})
		return
	}

	mapQuestion := make(map[int]domain.Question)
	for _, question := range exercise.Questions {
		mapQuestion[question.ID] = question
	}

	var score int
	for _, answer := range answers {
		if strings.EqualFold(answer.Answer, mapQuestion[answer.QuestionID].CorrectAnswer) {
			score += mapQuestion[answer.QuestionID].Score
		}
	}
	c.JSON(200, map[string]interface{}{
		"score": score,
	})
}

func (eu ExerciseUsecase) CreateExercise(c *gin.Context) {
	var exercise domain.Exercise
	err := c.ShouldBind(&exercise)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if exercise.Title == "" {
		c.JSON(400, gin.H{
			"message": "field title must required",
		})
		return
	}

	if exercise.Description == "" {
		c.JSON(400, gin.H{
			"message": "field description must required",
		})
		return
	}

	err = eu.db.Create(&exercise).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create exercise",
		})
		return
	}

	c.JSON(201, gin.H{
		"id":          exercise.ID,
		"title":       exercise.Title,
		"description": exercise.Description,
	})
}

func (eu ExerciseUsecase) CreateQuestion(c *gin.Context) {
	var question domain.Question

	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}
	question.ExerciseID = id

	question.CreatorID = int(c.Request.Context().Value("user_id").(float64))

	err = c.ShouldBind(&question)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if question.Body == "" {
		c.JSON(400, gin.H{
			"message": "field body must required",
		})
		return
	}

	if question.OptionA == "" {
		c.JSON(400, gin.H{
			"message": "field option_a must required",
		})
		return
	}

	if question.OptionB == "" {
		c.JSON(400, gin.H{
			"message": "field option_b must required",
		})
		return
	}

	if question.OptionC == "" {
		c.JSON(400, gin.H{
			"message": "field option_c must required",
		})
		return
	}

	if question.OptionD == "" {
		c.JSON(400, gin.H{
			"message": "field option_d must required",
		})
		return
	}

	if question.CorrectAnswer == "" {
		c.JSON(400, gin.H{
			"message": "field correct_answer must required",
		})
		return
	}

	if question.Score == 0 {
		c.JSON(400, gin.H{
			"message": "field score must required",
		})
		return
	}

	err = eu.db.Create(&question).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create question",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "berhasil membuat question",
	})
}

func (eu ExerciseUsecase) CreateAnswer(c *gin.Context) {
	var answer domain.Answer

	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}
	answer.ExerciseID = id

	paramIdQuestion := c.Param("qid")
	qid, err := strconv.Atoi(paramIdQuestion)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid question id",
		})
		return
	}
	answer.QuestionID = qid

	answer.UserID = int(c.Request.Context().Value("user_id").(float64))

	err = c.ShouldBind(&answer)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if answer.Answer == "" {
		c.JSON(400, gin.H{
			"message": "field answer must required",
		})
		return
	}

	err = eu.db.Create(&answer).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create answer",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "berhasil membuat answer",
	})
}
