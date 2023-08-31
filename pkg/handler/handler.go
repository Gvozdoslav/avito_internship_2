package handler

import (
	"avito2/pkg/handler/response"
	"avito2/pkg/service"
	"avito2/pkg/service/dto"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.PUT("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.DELETE("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/all", handler.GetAllUsers)
			user.GET("/get/:id", handler.GetUserById)
			user.GET("/get/active_segments/:id", handler.GetUserActiveSegments)
			user.GET("/get/csv/:id", handler.GetUserSegmentsCsv)
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

// GetAllUsers @Summary Get all users
// @Description Get all users with segments they are in
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} []dto.UserDto
// @Failure 500 {object} error
// @Router /api/users [get]
func (handler *Handler) GetAllUsers(ctx *gin.Context) {

	users, err := handler.services.GetUsers()
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong :(")
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetUserById @Summary Get user
// @Description Get user with segments by user id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User Id"
// @Success 200 {object} dto.UserDto
// @Failure 500 {object} error
// @Router /api/user/{id} [get]
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

// GetUserActiveSegments @Summary Get active segments
// @Description Get user with active segments by user id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User Id"
// @Success 200 {object} dto.UserDto
// @Failure 500 {object} error
// @Router /api/user/get/active_segments/{id} [get]
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

// CreateUser @Summary Create user
// @Description Create user by it's id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User Id"
// @Success 201 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/create/{id} [post]
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

// AddUserToSegment @Summary Add user to the segment
// @Description Add user to the segment
// @Tags Users
// @Accept json
// @Produce json
// @Param userId body int true "User Id"
// @Param segment body dto.SegmentDto true "Segment"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/add_segment [post]
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

// AddUserToSegments @Summary Add user to the segments
// @Description Add user to the segments
// @Tags Users
// @Accept json
// @Produce json
// @Param userId body int true "User Id"
// @Param segment body []dto.SegmentDto true "Segments"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/add_segments [post]
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

// RemoveUserSegment @Summary Remove user from segment
// @Description Remove user segment
// @Tags Users
// @Accept json
// @Produce json
// @Param userSingleSegmentDto body dto.UserSingleSegmentDto true "User Single Segment Dto"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/remove_segment [delete]
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

// RemoveUserSegments @Summary Remove user from segments
// @Description Remove user segments
// @Tags Users
// @Accept json
// @Produce json
// @Param userId body int true "User Id"
// @Param segment body []dto.SegmentDto true "Segments"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/remove_segments [delete]
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

// UpdateUserSegments @Summary Update user segments
// @Description SET user segments
// @Tags Users
// @Accept json
// @Produce json
// @Param userId body int true "User Id"
// @Param segment body []dto.SegmentDto true "Segments"
// @Success 200 {object} dto.UserDto
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/user/update [put]
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

// DeleteUser @Summary Delete user
// @Description Delete user by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User Id"
// @Success 200 {object} nil
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/users/delete/{id} [delete]
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

// GetSegment @Summary Get segment
// @Description Get segment by slug
// @Tags Segment
// @Accept json
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} model.Segment
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/segment/{slug} [get]
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

// GetAllSegments @Summary Get all segments
// @Description Get all segments
// @Tags Segment
// @Accept json
// @Produce json
// @Success 200 {object} []model.Segment
// @Failure 500 {object} error
// @Router /api/segment/get/all [get]
func (handler *Handler) GetAllSegments(ctx *gin.Context) {

	segments, err := handler.services.GetAllSegments()
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong due getting segments!")
		return
	}

	ctx.JSON(http.StatusOK, segments)
}

// CreateSegment @Summary Create segment
// @Description Create segment by slug
// @Tags Segment
// @Accept json
// @Produce json
// @Param slug path string true "Slug"
// @Success 201 {object} model.Segment
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/segment/create/{slug} [post]
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

	ctx.JSON(http.StatusCreated, segment)
}

// DeleteSegment @Summary Delete segment
// @Description Delete segment by slug
// @Tags Segment
// @Accept json
// @Produce json
// @Param slug path string true "Slug"
// @Success 201 {object} nil
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/segment/delete/{slug} [delete]
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

func (handler *Handler) GetUserSegmentsCsv(ctx *gin.Context) {

	userId, err := parseIntFromString(ctx.Param(id))
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusBadRequest, "User id is not specified")
		return
	}

	csvUrl, err := handler.services.GetUserSegmentsDataCsvUrl(userId)
	if err != nil {
		response.NewErrorResponse(ctx, http.StatusInternalServerError, "Something wen wrong due getting the csv :/")
		return
	}

	ctx.JSON(http.StatusOK, csvUrl)
}

func parseIntFromString(s string) (int, error) {
	userId, err := strconv.ParseInt(s, 10, 32)
	return int(userId), err
}
