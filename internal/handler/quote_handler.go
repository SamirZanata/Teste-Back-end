package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/back-end/quote-api/internal/domain"
	"github.com/back-end/quote-api/internal/service"
)

type QuoteHandler struct {
	svc *service.QuoteService
}

func NewQuoteHandler(svc *service.QuoteService) *QuoteHandler {
	return &QuoteHandler{svc: svc}
}

func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req domain.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendValidationError(c, err)
		return
	}

	resp, err := h.svc.CreateQuote(c.Request.Context(), &req)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *QuoteHandler) sendValidationError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, fieldErrorToMessage(e))
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados de entrada inválidos",
			"details": msgs,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados de entrada inválidos",
			"details": []string{"Corpo da requisição inválido. Verifique o JSON enviado (campos obrigatórios e formato)."},
		})
		return
	}
}

func fieldNameInPortuguese(field string) string {
	names := map[string]string{
		"Zipcode":       "CEP (recipient.address.zipcode)",
		"Address":       "Endereço do destinatário (recipient.address)",
		"Recipient":     "Destinatário (recipient)",
		"Volumes":       "Lista de volumes (volumes)",
		"Category":      "Categoria do volume",
		"Amount":        "Quantidade do volume",
		"UnitaryWeight": "Peso unitário (unitary_weight)",
		"Price":         "Preço do volume (price)",
		"Height":        "Altura do volume (height)",
		"Width":         "Largura do volume (width)",
		"Length":        "Comprimento do volume (length)",
	}
	if n, ok := names[field]; ok {
		return n
	}
	return field
}

func fieldErrorToMessage(e validator.FieldError) string {
	field := fieldNameInPortuguese(e.Field())
	switch e.Tag() {
	case "required":
		return field + " é obrigatório"
	case "min":
		return field + " deve ser no mínimo " + e.Param()
	case "max":
		return field + " deve ser no máximo " + e.Param()
	case "len":
		return field + " deve ter exatamente " + e.Param() + " caracteres"
	case "gt":
		return field + " deve ser maior que " + e.Param()
	case "gte":
		return field + " deve ser maior ou igual a " + e.Param()
	default:
		return field + ": " + e.Tag()
	}
}

func (h *QuoteHandler) sendError(c *gin.Context, err error) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "zipcode"):
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	case strings.Contains(msg, "Frete Rápido"):
		c.JSON(http.StatusBadGateway, gin.H{"error": msg})
	case strings.Contains(msg, "salvar"):
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao processar cotação"})
	}
}
