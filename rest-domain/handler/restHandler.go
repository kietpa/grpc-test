package handler

import (
	"context"
	"net/http"
	"ngc-grpc/helpers"
	"ngc-grpc/model"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Grpc  model.UserServiceClient
	Cache *redis.Client
}

func NewHandler(grpc model.UserServiceClient, cache *redis.Client) *Handler {
	return &Handler{
		Grpc:  grpc,
		Cache: cache,
	}
}

func (h *Handler) CreateUser(c echo.Context) error {
	// bind json input
	var input model.UserAll
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request" + err.Error(),
		})
	}

	in := model.UserRegister{
		Id:       input.ID,
		Name:     input.Name,
		Password: input.Password,
	}

	user, err := h.Grpc.CreateUser(context.Background(), &in)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"Status": "Berhasil",
		"ID":     user.Id,
		"Name":   user.Name,
	})
}

func (h *Handler) LoginUser(c echo.Context) error {
	ctx := context.Background()

	// bind json input
	var input model.UserNamePass
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request" + err.Error(),
		})
	}

	in := model.UserLogin{
		Name:     input.Name,
		Password: input.Password,
	}

	resp, err := h.Grpc.VerifyUserCredentials(context.Background(), &in)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	// store in redis
	user := model.UserCache{
		ID:    resp.Id,
		Name:  resp.Name,
		Token: resp.Token,
	}

	err = helpers.SetUserCache(h.Cache, &user, ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": resp.Token,
	})
}

func (h *Handler) GetUserList(c echo.Context) error {
	list, err := h.Grpc.GetUserList(context.Background(), &model.Empty{})
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": err,
		})
	}
	var result []*model.UserModel
	for _, v := range list.List {
		u := model.UserModel{
			ID:   v.Id,
			Name: v.Name,
		}
		result = append(result, &u)
	}

	return c.JSON(200, result)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	ctx := context.Background()

	// check cache if user info exists
	userCache, err := helpers.GetUserCache(h.Cache, id, ctx)
	if err == nil {
		return c.JSON(200, model.UserModel{
			ID:   userCache.ID,
			Name: userCache.Name,
		})
	}

	in := &model.ID{Id: id}

	user, err := h.Grpc.GetUserByID(context.Background(), in)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": "not found",
		})
	}

	// store in redis
	userStore := model.UserCache{
		ID:    user.Id,
		Name:  user.Name,
		Token: c.Request().Header.Get("Authorization"),
	}

	err = helpers.SetUserCache(h.Cache, &userStore, ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	return c.JSON(200, model.UserModel{
		ID:   user.Id,
		Name: user.Name,
	})
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	in := &model.ID{Id: id}

	_, err := h.Grpc.DeleteUser(context.Background(), in)
	if err != nil {
		return c.JSON(400, echo.Map{
			"message": err,
		})
	}

	return c.JSON(200, echo.Map{"Status": "Berhasil Menghapus Data"})
}

func (h *Handler) UpdateUser(c echo.Context) error {
	oldID := c.Param("id")

	var input model.UserModel
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request" + err.Error(),
		})
	}

	in := model.UserUpdate{
		Id:    oldID,
		Name:  input.Name,
		Newid: input.ID,
	}

	user, err := h.Grpc.UpdateUser(context.Background(), &in)
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": err,
		})
	}

	return c.JSON(200, echo.Map{
		"Status": "Berhasil Perbarui Data",
		"ID":     user.Id,
		"Name":   user.Name,
	})
}
