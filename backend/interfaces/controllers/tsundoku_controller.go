package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/yot-sailing/TSUNTSUN/domain"
	"github.com/yot-sailing/TSUNTSUN/interfaces/database"
	"github.com/yot-sailing/TSUNTSUN/usecase"
)

type TsundokuController struct {
	Interactor usecase.TsundokuInteractor
}

func NewTsundokuController(sqlHandler database.SqlHandler) *TsundokuController {
	return &TsundokuController{
		Interactor: usecase.TsundokuInteractor{
			TsundokuRepository: &database.TsundokuRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *TsundokuController) CreateTsundoku(c echo.Context, userID int) {
	tsundoku := domain.Tsundoku{}
	t, _ := time.Parse("2006-01-02", c.FormValue("deadline"))
	c.Bind(&tsundoku)
	tsundoku.Deadline = t
	tsundoku.UserID = userID
	controller.Interactor.Add(tsundoku)
	createdTsundokus := controller.Interactor.GetInfo(userID)
	c.JSON(201, createdTsundokus)
	return
}

func (controller *TsundokuController) GetFreeTsundoku(c echo.Context, userID int, free_time int) []domain.Tsundoku {
	res := controller.Interactor.GetInfo(userID)
	results := []domain.Tsundoku{}
	for _, element := range res {
		if element.Category == "site" {
			need_time := strings.Replace(element.RequiredTime, "min", "", -1)
			required_time, _ := strconv.Atoi(need_time)
			if free_time >= required_time {
				results = append(results, element)
			}
		}
	}
	return results
}

func (controller *TsundokuController) GetTsundoku(userID int) []domain.Tsundoku {
	tsundokus := controller.Interactor.GetInfo(userID)
	return tsundokus
}

func (controller *TsundokuController) Delete(id int) {
	controller.Interactor.Delete(id)
}
