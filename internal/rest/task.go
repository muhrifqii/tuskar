package rest

import (
	"context"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/rest/rest_utils"
	"go.uber.org/zap"
)

type (
	TaskService interface {
		GetAllTasks(ctx context.Context, params *domain.TaskQueryParams) ([]domain.Task, int, error)
		GetByID(ctx context.Context, id string) (domain.Task, error)
		CreateTask(ctx context.Context, task *domain.Task) error
		UpdateTask(ctx context.Context, task *domain.Task) error
		DeleteTask(ctx context.Context, id string) error
	}

	TaskHandler struct {
		taskService TaskService
		validator   *validator.Validate
		log         *zap.Logger
	}
)

func NewTaskHandler(
	router fiber.Router,
	svc TaskService,
	params rest_utils.HandlerParams,
) *TaskHandler {
	handler := &TaskHandler{
		taskService: svc,
		validator:   params.Validator,
		log:         params.Logger,
	}

	taskRoute := router.Group("/tasks")
	taskRoute.Get("/:id", handler.GetByID)
	taskRoute.Get("", handler.GetAllTasks)
	taskRoute.Post("", handler.CreateTask)
	taskRoute.Put("/:id", handler.UpdateTask)
	taskRoute.Delete("/:id", handler.DeleteTask)

	return handler
}

func (h *TaskHandler) GetAllTasks(c *fiber.Ctx) error {
	params := domain.TaskQueryParams{
		Page:  1,
		Limit: 10,
	}
	if err := c.QueryParser(&params); err != nil {
		return rest_utils.NewApiErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	h.log.Debug("param check", zap.Any("param", params))
	if err := h.validator.Struct(&params); err != nil {
		return err
	}
	tasks, total, err := h.taskService.GetAllTasks(c.Context(), &params)
	if err != nil {
		return err
	}
	totalPage := int(math.Ceil(float64(total) / float64(params.Limit)))
	paging := domain.Pagination{
		CurrentPage: params.Page,
		TotalPage:   totalPage,
		TotalTask:   total,
	}
	return rest_utils.ReturnAnyResponse(c, fiber.StatusOK, domain.TaskPaginatedResponse{
		Tasks:      tasks,
		Pagination: paging,
	})
}

func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := h.taskService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return rest_utils.ReturnOkResponse(c, task)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var task domain.Task
	if err := c.BodyParser(&task); err != nil {
		return rest_utils.NewApiErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(&task); err != nil {
		return err
	}
	err := h.taskService.CreateTask(c.Context(), &task)
	if err != nil {
		return err
	}
	return rest_utils.ReturnAnyResponse(c, fiber.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task domain.Task
	if err := c.BodyParser(&task); err != nil {
		return rest_utils.NewApiErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	task.ID = id
	if err := h.validator.Struct(&task); err != nil {
		return err
	}
	err := h.taskService.UpdateTask(c.Context(), &task)
	if err != nil {
		return err
	}
	return rest_utils.ReturnAnyResponse(c, fiber.StatusOK, domain.TaskMutationResponse{
		Message: "Task updated successfully",
		Task:    &task,
	})
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.taskService.DeleteTask(c.Context(), id)
	if err != nil {
		return err
	}
	return rest_utils.ReturnAnyResponse(c, fiber.StatusOK, domain.TaskMutationResponse{
		Message: "Task deleted successfully",
	})
}
