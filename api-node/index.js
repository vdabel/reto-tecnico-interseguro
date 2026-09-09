const express = require('express');
const cors = require('cors');
const { authenticateJWT } = require('./src/middleware/auth');
const { processStats } = require('./src/services/statsService');

const app = express();
const PORT = process.env.PORT || 4000;

// Middlewares
app.use(cors());
app.use(express.json({ limit: '10mb' }));

// Logging simple para auditoría
app.use((req, res, next) => {
  const start = Date.now();
  res.on('finish', () => {
    const duration = Date.now() - start;
    console.log(`[${new Date().toISOString()}] ${req.method} ${req.originalUrl} ${res.statusCode} - ${duration}ms`);
  });
  next();
});

// Endpoint de salud
app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'api-node',
    version: '1.0.0'
  });
});

// Endpoint principal protegido por JWT para cálculo estadístico
app.post('/api/stats', authenticateJWT, (req, res) => {
  try {
    const stats = processStats(req.body);

    res.json({
      success: true,
      data: stats
    });
  } catch (error) {
    res.status(400).json({
      success: false,
      error: error.message || 'Error procesando estadísticas de matriz'
    });
  }
});

// Manejador global de errores
app.use((err, req, res, next) => {
  console.error('Error no controlado:', err);
  res.status(500).json({
    success: false,
    error: 'Error interno del servidor en api-node'
  });
});

app.listen(PORT, () => {
  console.log(`🚀 API Node.js (Estadísticas) escuchando en el puerto ${PORT}...`);
});
