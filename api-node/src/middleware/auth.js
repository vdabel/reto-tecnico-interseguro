const jwt = require('jsonwebtoken');

const JWT_SECRET = process.env.JWT_SECRET || 'interseguro_super_secret_key_2024';

/**
 * Middleware para validar el token JWT en el encabezado Authorization
 */
function authenticateJWT(req, res, next) {
  const authHeader = req.headers.authorization;

  if (!authHeader) {
    return res.status(401).json({
      success: false,
      error: 'Falta el encabezado de autorización (Authorization: Bearer <token>)'
    });
  }

  const parts = authHeader.split(' ');
  if (parts.length !== 2 || parts[0].toLowerCase() !== 'bearer') {
    return res.status(401).json({
      success: false,
      error: "Formato de token inválido. Se espera 'Bearer <token>'"
    });
  }

  const token = parts[1];

  jwt.verify(token, JWT_SECRET, (err, user) => {
    if (err) {
      return res.status(401).json({
        success: false,
        error: 'Token JWT inválido o expirado',
        details: err.message
      });
    }

    req.user = user;
    next();
  });
}

module.exports = {
  authenticateJWT
};
