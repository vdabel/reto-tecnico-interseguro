package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-go/config"
	"api-go/services"

	"github.com/gofiber/fiber/v2"
)

type MatrixHandler struct {
	cfg           *config.Config
	matrixService *services.MatrixService
	httpClient    *http.Client
}

func NewMatrixHandler(cfg *config.Config, matrixService *services.MatrixService) *MatrixHandler {
	return &MatrixHandler{
		cfg:           cfg,
		matrixService: matrixService,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ProcessMatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type NodeStatsPayload struct {
	Matrices map[string][][]float64 `json:"matrices"`
}

// ProcessMatrix maneja el endpoint POST para recibir una matriz, realizar la factorización QR y rotación,
// y despachar los resultados a la API de Node.js para el cálculo estadístico.
func (h *MatrixHandler) ProcessMatrix(c *fiber.Ctx) error {
	var req ProcessMatrixRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Cuerpo de petición inválido. Se espera JSON con el campo 'matrix' como array de arrays de números",
		})
	}

	// 1. Factorización QR mediante Gonum
	q, r, err := h.matrixService.FactorizeQR(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Error en procesamiento de matriz: %s", err.Error()),
		})
	}

	// 2. Rotación de 90 grados (solución dual para satisfacer requerimientos de rotación y QR del PDF)
	rotated, err := h.matrixService.Rotate90(req.Matrix)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Error al rotar matriz: %s", err.Error()),
		})
	}

	// 3. Preparar payload para la API de Node.js (estrictamente las matrices Q y R devueltas por QR)
	nodePayload := NodeStatsPayload{
		Matrices: map[string][][]float64{
			"q": q,
			"r": r,
		},
	}

	payloadBytes, err := json.Marshal(nodePayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error al serializar datos para la API de estadísticas",
		})
	}

	// 4. Enviar mediante HTTP POST a la API de Node.js protegiendo con el mismo token JWT
	nodeURL := fmt.Sprintf("%s/api/stats", h.cfg.NodeAPIURL)
	httpReq, err := http.NewRequest("POST", nodeURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error al inicializar la petición hacia la API de Node.js",
		})
	}

	httpReq.Header.Set("Content-Type", "application/json")
	// Reenviar Authorization header para mantener la cadena de seguridad JWT
	if authHeader := c.Get("Authorization"); authHeader != "" {
		httpReq.Header.Set("Authorization", authHeader)
	}

	resp, err := h.httpClient.Do(httpReq)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("No se pudo conectar con la API de Node.js en %s: %s", nodeURL, err.Error()),
			"qr": fiber.Map{
				"q": q,
				"r": r,
			},
			"rotated": rotated,
		})
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"error":   "Error al leer la respuesta de la API de Node.js",
		})
	}

	var statsResult interface{}
	if err := json.Unmarshal(respBody, &statsResult); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success":     false,
			"error":       "Respuesta no válida recibida de la API de Node.js",
			"rawResponse": string(respBody),
		})
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"success":    false,
			"error":      "La API de Node.js devolvió un código de error",
			"nodeErrors": statsResult,
		})
	}

	// 5. Retornar respuesta consolidada de alto nivel
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Matriz procesada y estadísticas calculadas exitosamente",
		"data": fiber.Map{
			"original": req.Matrix,
			"qr": fiber.Map{
				"q": q,
				"r": r,
			},
			"rotated":    rotated,
			"statistics": statsResult,
		},
	})
}
