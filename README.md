# Reto Técnico - Microservicios Interseguro

Este repositorio contiene mi solución al reto técnico. El proyecto está estructurado como un ecosistema de microservicios contenerizados donde la API en **Go (Fiber)** procesa matrices mediante **Factorización QR**, se comunica vía HTTP con la API en **Node.js (Express)** que calcula estadísticas sobre las matrices resultantes, y se incluye un **Frontend** para probar todo interactivamente.

---

## ¿Cómo funciona el flujo?

1. El usuario introduce una matriz rectangular en el **Frontend** o mediante **cURL**.
2. La **API en Go (`:3000`)** recibe la matriz, valida el token JWT y realiza la **Factorización QR** ($A = Q \cdot R$) usando `gonum/mat`, además de una rotación de 90°.
3. La API en Go hace una petición interna por HTTP POST a la **API en Node.js (`:4000`)** enviándole **únicamente las dos matrices generadas por el QR ($Q$ y $R$)**, reenviando la cabecera de autenticación JWT.
4. La **API en Node.js** analiza las matrices $Q$ y $R$, calcula las métricas solicitadas (máximo, mínimo, promedio, suma total y si alguna es diagonal) y responde con un JSON estructurado.
5. La **API en Go** une los resultados matemáticos con las estadísticas recibidas y le devuelve al cliente la respuesta completa para que se renderice en pantalla.

---
## ¿Qué devuelve cada API? (Detalle de Endpoints)

### 1. API en Go (`http://localhost:3000`)

#### A. Generación de Token JWT
* **Ruta:** `POST /api/auth/token`
* **Acceso:** Público (endpoint de utilidad para pruebas y frontend)
* **Body que recibe (opcional):**
  ```json
  { "username": "mi-usuario" }
  ```
* **Qué devuelve:**
  ```json
  {
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires": "2026-09-10T17:30:00Z",
    "user": "mi-usuario"
  }
  ```
---

#### B. Procesamiento de Matriz (QR + Rotación + Estadísticas)
* **Ruta:** `POST /api/matrix/process`
* **Acceso:** Protegido con `Authorization: Bearer <token>`
* **Body que recibe:**
  ```json
  {
    "matrix": [
      [12, -51, 4],
      [6, 167, -68],
      [-4, 24, -41]
    ]
  }
  ```
* **Qué devuelve exactamente:**
  ```json
  {
    "success": true,
    "message": "Matriz procesada y estadísticas calculadas exitosamente",
    "data": {
      "original": [
        [12, -51, 4],
        [6, 167, -68],
        [-4, 24, -41]
      ],
      "qr": {
        "q": [
          [-0.857143, 0.394286, 0.331429],
          [-0.428571, -0.902857, -0.034286],
          [0.285714, -0.171429, 0.942857]
        ],
        "r": [
          [-14, -21, 14],
          [0, -175, 70],
          [0, 0, -35]
        ]
      },
      "rotated": [
        [-4, 6, 12],
        [24, 167, -51],
        [-41, -68, 4]
      ],
      "statistics": {
        "data": {
          "global": {
            "maxValue": 70,
            "minValue": -175,
            "average": -8.968889,
            "totalSum": -161.44,
            "totalElements": 18,
            "isAnyDiagonal": false,
            "diagonalMatrices": []
          },
          "byMatrix": {
            "q": {
              "maxValue": 0.942857,
              "minValue": -0.902857,
              "average": -0.048889,
              "totalSum": -0.44,
              "totalElements": 9,
              "isDiagonal": false,
              "dimensions": { "rows": 3, "cols": 3 }
            },
            "r": {
              "maxValue": 70,
              "minValue": -175,
              "average": -17.888889,
              "totalSum": -161,
              "totalElements": 9,
              "isDiagonal": false,
              "dimensions": { "rows": 3, "cols": 3 }
            }
          }
        },
        "success": true
      }
    }
  }
  ```

**Explicación de los campos devueltos:**
* `data.original`: La matriz tal cual como fue enviada.
* `data.qr.q`: La matriz ortogonal $Q$ obtenida de la factorización.
* `data.qr.r`: La matriz triangular superior $R$ obtenida de la factorización.
* `data.rotated`: La matriz rotada 90° en sentido horario.
* `data.statistics`: La respuesta íntegra calculada por la API de Node.js sobre $Q$ y $R$.

---

### 2. API en Node.js (`http://localhost:4000`)

#### Cálculo Estadístico de Matrices
* **Ruta:** `POST /api/stats`
* **Acceso:** Protegido con `Authorization: Bearer <token>`
* **Body que recibe:**
  Un objeto JSON con las matrices que se desean analizar (en el flujo regular, Go le envía únicamente $Q$ y $R$):
  ```json
  {
    "matrices": {
      "q": [
        [-0.857143, 0.394286, 0.331429],
        [-0.428571, -0.902857, -0.034286],
        [0.285714, -0.171429, 0.942857]
      ],
      "r": [
        [-14, -21, 14],
        [0, -175, 70],
        [0, 0, -35]
      ]
    }
  }
  ```
* **Qué devuelve exactamente:**
  ```json
  {
    "success": true,
    "data": {
      "global": {
        "maxValue": 70,
        "minValue": -175,
        "average": -8.968889,
        "totalSum": -161.44,
        "totalElements": 18,
        "isAnyDiagonal": false,
        "diagonalMatrices": []
      },
      "byMatrix": {
        "q": {
          "maxValue": 0.942857,
          "minValue": -0.902857,
          "average": -0.048889,
          "totalSum": -0.44,
          "totalElements": 9,
          "isDiagonal": false,
          "dimensions": { "rows": 3, "cols": 3 }
        },
        "r": {
          "maxValue": 70,
          "minValue": -175,
          "average": -17.888889,
          "totalSum": -161,
          "totalElements": 9,
          "isDiagonal": false,
          "dimensions": { "rows": 3, "cols": 3 }
        }
      }
    }
  }
  ```

**Explicación de las métricas devueltas:**
* `global.maxValue`: El valor numérico más alto encontrado en el conjunto de matrices evaluadas.
* `global.minValue`: El valor numérico más bajo encontrado en el conjunto de matrices.
* `global.average`: El promedio general (suma total dividida entre el total de elementos de todas las matrices).
* `global.totalSum`: La suma acumulada de todos los números presentes en las matrices.
* `global.totalElements`: Cantidad total de celdas analizadas (ejemplo: para dos matrices de 3x3, suma exactamente 18 elementos).
* `global.isAnyDiagonal`: Bandera booleana (`true` o `false`). Indica si al menos una de las matrices recibidas es una matriz diagonal.
* `global.diagonalMatrices`: Arreglo con los nombres de las matrices que cumplieron la condición diagonal (por ejemplo `["q"]` o `[]`).
* `byMatrix`: El desglose detallado de estas mismas estadísticas calculado para cada matriz de forma independiente, incluyendo sus dimensiones (`rows` y `cols`).

---

## Cómo deplegarlo

### Con Docker

**Ejecuta en la raíz del proyecto:**

```bash
docker compose up --build -d
```

Una vez levantado, ingresa desde tu navegador a: **[http://localhost:8080](http://localhost:8080)**

**Para apagar los contenedores y liberar recursos:**
```bash
docker compose down
```

---

## Pruebas Unitarias

Ambos microservicios cuentan con pruebas unitarias que cubren los casos de éxito, bordes y validaciones de datos:

### Pruebas en Go (Factorización QR y Rotación)
```bash
cd api-go
go test ./services -v
```
*Verifica que $Q$ sea ortogonal, que $R$ sea triangular superior, que $Q \cdot R \approx A$, la rotación de 90° y que se rechacen matrices irregulares o vacías.*

### Pruebas en Node.js (Cálculos Estadísticos y Matriz Diagonal)
```bash
cd api-node
npm test
```
---

## Estructura del Repositorio

```text
reto-interseguro/
├── docker-compose.yml          # Configuración de los 3 contenedores y su red interna
├── README.md                   # Esta documentación
│
├── api-go/                     # Microservicio en Go (Fiber + Gonum)
│   ├── Dockerfile              # Construcción multietapa (golang:alpine -> alpine)
│   ├── go.mod / go.sum         # Dependencias (Fiber, Gonum, JWT)
│   ├── main.go                 # Servidor y rutas
│   ├── config/                 # Lectura de variables de entorno
│   ├── handlers/               # Controladores de matriz y autenticación
│   ├── middleware/             # Middleware de validación JWT
│   └── services/               # Lógica QR, rotación y pruebas unitarias
│
├── api-node/                   # Microservicio en Node.js (Express + JWT)
│   ├── Dockerfile              # Contenedor node:18-alpine
│   ├── package.json            # Dependencias (Express, jsonwebtoken, cors)
│   ├── index.js                # Servidor y endpoint /api/stats
│   └── src/
│       ├── middleware/         # Middleware de validación JWT
│       └── services/           # Lógica estadística y pruebas unitarias
│
└── frontend/                   # Aplicación Web (Panel de Control)
    ├── Dockerfile              # Servidor Nginx Alpine
    ├── nginx.conf              # Configuración de Nginx
    └── index.html              # Interfaz interactiva con Tailwind CSS y Vanilla JS
```
