package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"transaction_api/constants"
	"transaction_api/dto"
	"transaction_api/dto/transaction"
	"transaction_api/model"
	"transaction_api/service"
	"transaction_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	accountService     *service.AccountService
	validate           *validator.Validate
}

func NewTransactionHandler(transactionService *service.TransactionService, accountService *service.AccountService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		accountService:     accountService,
		validate:           validator.New(),
	}
}

func (h *TransactionHandler) InsertTransactionHandler(c *gin.Context) {
	var req transaction.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	if err := h.validate.Struct(req.Request); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      utils.FormatValidationErrors(err)["reason"],
		})
		return
	}

	data := req.Request
	checkDebet, err := h.accountService.GetDetailAccountService(c, int64(data.ROWID_SENDER))
	if err != nil {
		log.Println(err)
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedCheckBalance,
		})
		return
	}

	if strings.ToUpper(checkDebet.STATUS) != "ACTIVE" {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Debet account is not active",
		})
		return
	}

	if checkDebet.BALANCE < data.AMOUNT {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedInsufficientBalance,
		})
		return
	}

	_, err = h.transactionService.SaveTransactionService(
		c.Request.Context(),
		data,
	)

	if err != nil {
		log.Println(err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedInsertTransaction,
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessInsertTransaction,
	})
}

func (h *TransactionHandler) GetAllTransactionHandler(c *gin.Context) {
	var req transaction.RequestSearch
	ctx := c.Request.Context()
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      utils.FormatValidationErrors(err)["reason"],
		})
		return
	}

	if req.Request.DATE_TO != "" && req.Request.DATE_FROM != "" {
		from, _ := time.Parse("2006-01-02", req.Request.DATE_FROM)
		to, _ := time.Parse("2006-01-02", req.Request.DATE_TO)
		if to.Sub(from).Hours()/24 > 30 {
			utils.WriteJSON(c, http.StatusBadRequest, model.Response{
				Status:       constants.StatusFailed,
				ResponseCode: constants.CodeFailed,
				Message:      "Date range must not exceed 30 days",
			})
			return
		}

		if from.After(to) {
			utils.WriteJSON(c, http.StatusBadRequest, model.Response{
				Status:       constants.StatusFailed,
				ResponseCode: constants.CodeFailed,
				Message:      "dateFrom must not be greater than dateTo",
			})
			return

		}
	}

	data, total, err := h.transactionService.GetAllTransactionsService(ctx, limit, offset, req.Request.REFERENCE_NUMBER, req.Request.DATE_FROM, req.Request.DATE_TO)
	if err != nil {
		log.Println(err)
		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetTransaction,
		})
		return
	}

	pagination := utils.BuildPagination(page, limit, total)

	utils.WriteJSON(c, http.StatusOK, utils.PaginatedResponse{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetTransaction,
		Data:         data,
		Pagination:   pagination,
	})
}

func (h *TransactionHandler) GetDetailTransactionHandler(c *gin.Context) {
	var req dto.User

	rowId := c.Param("rowid_trx")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      utils.FormatValidationErrors(err)["reason"],
		})
		return
	}

	id, err := strconv.Atoi(rowId)
	if err != nil || id <= 0 {

		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Invalid rowid_trx parameter",
		})
		return
	}

	data, err := h.transactionService.GetDetailTransactionService(c, rowId)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Failed to get detail transaction : " + err.Error(),
		})
		return
	}

	if data == nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      "Transaction not found",
		})
		return
	}

	utils.WriteJSON(c, http.StatusOK, model.Response{
		Status:       constants.StatusSuccess,
		ResponseCode: constants.CodeSuccess,
		Message:      constants.MessageSuccessGetTransaction,
		Data:         data,
	})
}

func (h *TransactionHandler) GetDownloadTransactionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	w := c.Writer

	var req transaction.RequestSearch

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedBindJSON,
		})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      utils.FormatValidationErrors(err)["reason"],
		})
		return
	}

	data, total, err := h.transactionService.GetAllTransactionsService(ctx, -1, 0, req.Request.REFERENCE_NUMBER, req.Request.DATE_FROM, req.Request.DATE_TO)
	if err != nil {

		utils.WriteJSON(c, http.StatusInternalServerError, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MessageFailedGetTransaction,
		})
		return
	}

	if total <= 0 {
		utils.WriteJSON(c, http.StatusBadRequest, model.Response{
			Status:       constants.StatusFailed,
			ResponseCode: constants.CodeFailed,
			Message:      constants.MMessageDataNotFoundTransaction,
		})
		return
	} else {

		f := excelize.NewFile()
		sheet := "Sheet1"
		style, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		styleBorder, err := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		f.SetCellValue(sheet, "A1", "No")
		f.SetCellValue(sheet, "B1", "Reference Number")
		f.SetCellValue(sheet, "C1", "Status")
		f.SetCellValue(sheet, "D1", "Transaction Date")
		f.SetCellValue(sheet, "E1", "Debet Account")
		f.SetCellValue(sheet, "F1", "Debet Name")
		f.SetCellValue(sheet, "G1", "Credit Account")
		f.SetCellValue(sheet, "H1", "Credit Name")
		f.SetCellValue(sheet, "I1", "Amount")
		f.SetCellValue(sheet, "J1", "Remark")

		if err := f.SetCellStyle(sheet, "A1", "J1", style); err != nil {
			log.Fatal(err)
		}

		colWidths := []float64{6.29, 26, 30, 30, 20, 20, 20, 20, 20, 30}
		for i, w := range colWidths {
			col := intToExcelCol(int64(i + 1))
			if err := f.SetColWidth(sheet, col, col, w); err != nil {
				log.Fatal(err)
			}
		}

		if total != 0 {
			if err := f.SetCellStyle(sheet, "A2", "J"+strconv.FormatInt(total+1, 10), styleBorder); err != nil {
				log.Fatal(err)
			}
			for i, row := range data {
				formattedDate := ""
				if row.TRANSACTION_DATE != "" {
					t, err := time.Parse(time.RFC3339Nano, row.TRANSACTION_DATE)
					if err != nil {
						log.Println("Error parsing date:", err)

					}

					formattedDate = t.Format("2006-01-02 15:04:05")
				}

				rowNum := i + 2
				f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
				f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row.REFERENCE_NUMBER)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), fmt.Sprintf("%s : %s", row.STATUS_CODE, row.STATUS_DESC))
				f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), formattedDate)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), row.DEBET_ACCOUNT)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), row.DEBET_NAME)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), row.CREDIT_ACCOUNT)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), row.CREDIT_NAME)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", rowNum), row.AMOUNT)
				f.SetCellValue(sheet, fmt.Sprintf("J%d", rowNum), row.REMARK)

			}
		}

		now := time.Now()
		formatted := now.Format("20060102150405")
		var filename = "Transaction_" + formatted + ".xlsx"
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("File-Name", filename)

		if err := f.Write(w); err != nil {

			log.Println("Error writing file:", err)
			http.Error(w, "could not generate file", http.StatusInternalServerError)
			return
		}

	}
}

func intToExcelCol(n int64) string {
	result := ""
	for n > 0 {
		n--
		result = string(rune('A'+(n%26))) + result
		n /= 26
	}
	return result
}
