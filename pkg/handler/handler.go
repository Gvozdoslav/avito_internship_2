package handler

import (
	"avito2/pkg/handler/response"
	"avito2/pkg/service"
	"avito2/pkg/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	id   = "id"
	slug = "slug"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/all", handler.GetAllUsers)
			user.GET("/get/:id", handler.GetUserById)
			user.GET("/get/active_segments/:id", handler.GetUserActiveSegments)
			user.POST("/create/:id", handler.CreateUser)
			user.POST("/add_segment", handler.AddUserToSegment)
			user.POST("/add_segments", handler.AddUserToSegments)
			user.PUT("/update", handler.UpdateUserSegments)
			user.DELETE("/remove_segment", handler.RemoveUserSegment)
			user.DELETE("/remove_segments", handler.RemoveUserSegments)
			user.DELETE("/delete/:id", handler.DeleteUser)
		}

		segment := api.Group("/segment")
		{
			segment.GET("/get/:slug", handler.GetSegment)
			segment.GET("/get/all", handler.GetAllSegments)
			segment.POST("/create/:slug", handler.CreateSegment)
			segment.DELETE("/delete/:slug", handler.DeleteSegment)
		}
	}

	return router
}

func (handler *Handler) GetAllUsers(ctx *gin.Context) {

	users, err := handler.services.GetUsers()
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong :(")
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (handler *Handler) GetUserById(ctx *gin.Context) {

	userId, err := parseIntFromString(ctx.Param(id))
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User id entered incorrectly!")
		return
	}

	user, err := handler.services.GetUserById(userId)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusNotFound, "User was not found!")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) GetUserActiveSegments(ctx *gin.Context) {

	userId, err := parseIntFromString(ctx.Param(id))
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User id entered incorrectly!")
		return
	}

	user, err := handler.services.GetUserActiveSegments(userId)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusNotFound, "User was not found!")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) CreateUser(ctx *gin.Context) {

	userId, err := parseIntFromString(ctx.Param(id))
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User id entered incorrectly!")
		return
	}

	user, err := handler.services.CreateUser(userId)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User already exist!")
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (handler *Handler) AddUserToSegment(ctx *gin.Context) {

	var userSingleSegmentDto dto.UserSingleSegmentDto
	if err := ctx.BindJSON(&userSingleSegmentDto); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Given body is not valid")
		return
	}

	user, err := handler.services.AddUserToSegment(&userSingleSegmentDto)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due adding the segment")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) AddUserToSegments(ctx *gin.Context) {

	var userDto dto.UserDto
	if err := ctx.BindJSON(&userDto); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Given body is not valid")
		return
	}

	user, err := handler.services.AddUserToSegments(&userDto)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due adding segments")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) RemoveUserSegment(ctx *gin.Context) {

	var userSingleSegmentDto dto.UserSingleSegmentDto
	if err := ctx.BindJSON(&userSingleSegmentDto); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Given body is not valid")
		return
	}

	user, err := handler.services.RemoveUserFromSegment(&userSingleSegmentDto)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due removing the segment")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) RemoveUserSegments(ctx *gin.Context) {

	var userDto dto.UserDto
	if err := ctx.BindJSON(&userDto); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Given body is not valid")
		return
	}

	user, err := handler.services.RemoveUserFromSegments(&userDto)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due removing the segments")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) UpdateUserSegments(ctx *gin.Context) {

	var userDto dto.UserDto
	if err := ctx.BindJSON(&userDto); err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Given body is not valid")
		return
	}

	user, err := handler.services.UpdateUserSegments(&userDto)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due updating the user segments")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) DeleteUser(ctx *gin.Context) {

	userId, err := parseIntFromString(ctx.Param(id))
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User id entered incorrectly!")
		return
	}

	if err = handler.services.DeleteUser(userId); err != nil {
		response.NewErrorResponse(ctx, http.StatusNotFound, "Something went wrong due deleting the user!")
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (handler *Handler) GetSegment(ctx *gin.Context) {

	slug := ctx.Param(slug)
	if slug == "" || &slug == nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Slug is not specified!")
		return
	}

	segment, err := handler.services.GetSegment(slug)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due getting segment!")
		return
	}

	ctx.JSON(http.StatusOK, segment)
}

func (handler *Handler) GetAllSegments(ctx *gin.Context) {

	segments, err := handler.services.GetAllSegments()
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due getting segments!")
		return
	}

	ctx.JSON(http.StatusOK, segments)
}

func (handler *Handler) CreateSegment(ctx *gin.Context) {

	slug := ctx.Param(slug)
	if slug == "" || &slug == nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Slug is not specified!")
		return
	}

	segment, err := handler.services.CreateSegment(slug)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due creating segment!")
		return
	}

	ctx.JSON(http.StatusOK, segment)
}

func (handler *Handler) DeleteSegment(ctx *gin.Context) {

	slug := ctx.Param(slug)
	if slug == "" || &slug == nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "Slug is not specified!")
		return
	}

	if err := handler.services.DeleteSegment(slug); err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due creating segment!")
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func parseIntFromString(s string) (int, error) {
	userId, err := strconv.ParseInt(s, 10, 32)
	return int(userId), err
}
