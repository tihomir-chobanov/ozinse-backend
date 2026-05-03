package handler

import (
    "net/http"
    "strconv"
    "ozinse-backend/internal/model"
    "ozinse-backend/internal/service"
    "github.com/gin-gonic/gin"
)

type CategoryHandler struct {
    service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
    return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
    categories, err := h.service.GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    category, err := h.service.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "category with id " + c.Param("id") + " not found"})
        return
    }
    c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Create(c *gin.Context) {
    var category model.Category
    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.service.Create(&category); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var category model.Category
    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    category.ID = id
    if err := h.service.Update(&category); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		// Here we are returning a custom error message
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
}

/*

Core Purpose
The Handler's main job is to translate between the world of HTTP (web browsers, Postman, mobile apps) and the world of Go (structs, methods, and logic). It ensures that the rest of your application (the Service and Repository) doesn't have to deal with web-specific details like headers, cookies, or JSON parsing.

Key Responsibilities
Routing: It listens for specific "paths" and "methods." For example, it tells the app: "When you see a DELETE request at /api/categories/5, execute the Delete function."

Request Unmarshalling (Binding): It takes the raw JSON text you send in Postman and converts it into a Go Model. It uses the ShouldBindJSON method in Gin to do this.

Parameter Extraction: it pulls information out of the URL. In /api/categories/:id, the Handler is responsible for grabbing the :id and converting it from a string to an integer so the code can use it.

Input Validation: It performs "shallow" validation. It checks if the ID is a number or if the required JSON fields are present. If something is wrong, it stops the request immediately before it ever reaches the database.

Response Marshalling: After the Service and Repository finish their work, the Handler takes the Go result and converts it back into JSON to send to the user.

Status Code Assignment: The Handler decides the "tone" of the response by setting HTTP status codes:

200 OK: "Everything went great."

201 Created: "I successfully made the new category."

400 Bad Request: "You sent me the wrong data format."

404 Not Found: "That ID doesn't exist."

500 Internal Server Error: "Something crashed on my end."

NB--------------------------------------------------------

Imagine a user sending a POST request to create a category:

The Handler receives the raw JSON. It uses the Model as a template to turn that JSON into a Go object.

The Handler passes that Model object to the Service.

The Service performs logic and passes the Model to the Repository.

The Repository saves the data from the Model into the database.

Finally, the Handler takes the updated Model (now with an ID) and sends it back to the user.
*/