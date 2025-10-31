package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kgugunava/gorkycode_backend/internal/services"
)

var taskStore = make(map[string]*RouteTask)

type RouteTask struct {
	ID        string                    `json:"id"`
	Status    string                    `json:"status"` // "processing", "completed", "failed"
	Result    *services.RouteResponse   `json:"result,omitempty"`
	Error     string                    `json:"error,omitempty"`
	CreatedAt time.Time                 `json:"created_at"`
	UserID    uint                      `json:"user_id"`
}

// CreateRouteHandler - создает задачу на расчет маршрута
func (h *RouteHandler) CreateRouteHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request services.SendRouteInfoRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Генерируем ID задачи
	taskID := fmt.Sprintf("task_%d_%d", userID.(uint), time.Now().Unix())
	
	// Создаем задачу
	task := &RouteTask{
		ID:        taskID,
		Status:    "processing",
		CreatedAt: time.Now(),
		UserID:    userID.(uint),
	}
	taskStore[taskID] = task

	// Запускаем асинхронную обработку
	go h.processRouteAsync(taskID, request, userID.(uint))

	c.JSON(http.StatusAccepted, gin.H{
		"task_id": taskID,
		"status":  "processing",
		"message": "Маршрут рассчитывается, это может занять несколько минут",
	})
}

// RouteStatusHandler - возвращает статус задачи
func (h *RouteHandler) RouteStatusHandler(c *gin.Context) {
	taskID := c.Param("taskId")
	
	task, exists := taskStore[taskID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	response := gin.H{
		"task_id": task.ID,
		"status":  task.Status,
	}
	
	if task.Status == "completed" {
		response["result"] = task.Result
	} else if task.Status == "failed" {
		response["error"] = task.Error
	}
	
	c.JSON(http.StatusOK, response)
}

// processRouteAsync - асинхронно обрабатывает маршрут
func (h *RouteHandler) processRouteAsync(taskID string, request services.SendRouteInfoRequest, userID uint) {
    defer func() {
        if r := recover(); r != nil {
            task := taskStore[taskID]
            task.Status = "failed"
            task.Error = fmt.Sprintf("Panic: %v", r)
        }
    }()

    // Имитируем долгий расчет
    fmt.Printf("🔄 Начинаем расчет маршрута для задачи %s\n", taskID)
    time.Sleep(5 * time.Second)

    // Создаем места в формате, соответствующем ML-сервису
    places := []json.RawMessage{
        json.RawMessage(`{
            "addres": "",
            "coordinate": [56.310043, 44.001603],
            "description": "",
            "time_to_come": -1,
            "time_to_visit": 30,
            "title": "UserPoint",
            "url": ""
        }`),
        json.RawMessage(`{
            "addres": "Нижний Новгород, Октябрьская улица, 9Б",
            "coordinate": [56.321791, 44.00199],
            "description": "Галерея 9Б открылась 16 марта 2018 года в Нижнем Новгороде. Основатели - коллекционеры искусства Георгий Смирнов и Елена Тальянская.",
            "time_to_come": -1,
            "time_to_visit": 30,
            "title": "Галерея 9Б",
            "url": ""
        }`),
        json.RawMessage(`{
            "addres": "Нижний Новгород, Верхне-Волжская набережная, 7",
            "coordinate": [56.331091, 44.016678],
            "description": "Усадьба Рукавишниковых — одно из самых красивых зданий в историческом центре Нижнего Новгорода.",
            "time_to_come": -1,
            "time_to_visit": 45,
            "title": "Усадьба Рукавишниковых",
            "url": ""
        }`),
    }

    // Создаем результат согласно вашей структуре RouteResponse
    result := &services.RouteResponse{
        Description: "Пешеходный маршрут по историческому центру Нижнего Новгорода",
        Time:        120, // общее время маршрута в минутах
        CountPlaces: len(places),
        Places:      places,
    }

    // Обновляем задачу
    task := taskStore[taskID]
    task.Status = "completed"
    task.Result = result
    
    fmt.Printf("✅ Расчет маршрута завершен для задачи %s\n", taskID)
    
    // Логируем результат
    fmt.Printf("📝 Описание: %s\n", result.Description)
    fmt.Printf("⏱️ Общее время: %d минут\n", result.Time)
    fmt.Printf("📍 Количество мест: %d\n", result.CountPlaces)
}