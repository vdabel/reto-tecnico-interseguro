/**
 * Servicio de Cálculo de Estadísticas para Matrices
 */

const EPSILON = 1e-6;

/**
 * Valida si la entrada es una matriz rectangular válida de números
 * @param {Array<Array<number>>} matrix
 * @returns {{ valid: boolean, error?: string, rows?: number, cols?: number }}
 */
function validateMatrix(matrix) {
  if (!Array.isArray(matrix) || matrix.length === 0) {
    return { valid: false, error: 'La matriz debe ser un arreglo no vacío' };
  }

  const rows = matrix.length;
  if (!Array.isArray(matrix[0]) || matrix[0].length === 0) {
    return { valid: false, error: 'Las filas de la matriz no pueden estar vacías' };
  }

  const cols = matrix[0].length;

  for (let i = 0; i < rows; i++) {
    if (!Array.isArray(matrix[i])) {
      return { valid: false, error: `La fila en índice ${i} no es un arreglo válido` };
    }
    if (matrix[i].length !== cols) {
      return { valid: false, error: `Matriz irregular: la fila ${i} tiene ${matrix[i].length} columnas, se esperaban ${cols}` };
    }
    for (let j = 0; j < cols; j++) {
      const val = matrix[i][j];
      if (typeof val !== 'number' || Number.isNaN(val) || !Number.isFinite(val)) {
        return { valid: false, error: `Elemento no numérico en posición [${i}][${j}]` };
      }
    }
  }

  return { valid: true, rows, cols };
}

/**
 * Verifica si una matriz es diagonal (debe ser cuadrada y los elementos fuera de la diagonal deben ser 0)
 * @param {Array<Array<number>>} matrix
 * @returns {boolean}
 */
function isDiagonalMatrix(matrix) {
  const validation = validateMatrix(matrix);
  if (!validation.valid) return false;

  const { rows, cols } = validation;
  // Una matriz diagonal es obligatoriamente cuadrada
  if (rows !== cols) return false;

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j) {
        if (Math.abs(matrix[i][j]) > EPSILON) {
          return false;
        }
      }
    }
  }

  return true;
}

/**
 * Calcula las estadísticas individuales para una sola matriz
 * @param {Array<Array<number>>} matrix
 * @returns {Object}
 */
function calculateSingleMatrixStats(matrix) {
  const validation = validateMatrix(matrix);
  if (!validation.valid) {
    throw new Error(validation.error);
  }

  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  for (let i = 0; i < matrix.length; i++) {
    for (let j = 0; j < matrix[i].length; j++) {
      const val = matrix[i][j];
      if (val > max) max = val;
      if (val < min) min = val;
      sum += val;
      count++;
    }
  }

  const avg = count > 0 ? sum / count : 0;
  const isDiagonal = isDiagonalMatrix(matrix);

  return {
    maxValue: Number(max.toFixed(6)),
    minValue: Number(min.toFixed(6)),
    average: Number(avg.toFixed(6)),
    totalSum: Number(sum.toFixed(6)),
    totalElements: count,
    isDiagonal,
    dimensions: {
      rows: validation.rows,
      cols: validation.cols
    }
  };
}

/**
 * Procesa múltiples matrices o una sola matriz y produce estadísticas consolidadas y por matriz.
 * @param {Object|Array} input
 * @returns {Object}
 */
function processStats(input) {
  if (!input) {
    throw new Error('El cuerpo de la petición no contiene datos de matriz');
  }

  // Normalizar entrada a mapa de matrices { [nombre]: matriz }
  let matrixMap = {};

  if (input.matrices && typeof input.matrices === 'object' && !Array.isArray(input.matrices)) {
    matrixMap = input.matrices;
  } else if (typeof input === 'object' && !Array.isArray(input) && !input.matrix && !input.matrices) {
    matrixMap = input;
  } else if (input.matrix && Array.isArray(input.matrix)) {
    matrixMap = { inputMatrix: input.matrix };
  } else if (Array.isArray(input)) {
    // Si es un arreglo simple de matrices o matriz 2D
    if (input.length > 0 && Array.isArray(input[0]) && typeof input[0][0] === 'number') {
      matrixMap = { inputMatrix: input };
    } else {
      input.forEach((m, idx) => {
        matrixMap[`matrix_${idx + 1}`] = m;
      });
    }
  } else {
    throw new Error('Formato de datos no reconocido para el procesamiento estadístico');
  }

  const keys = Object.keys(matrixMap);
  if (keys.length === 0) {
    throw new Error('No se proporcionaron matrices válidas para calcular estadísticas');
  }

  const byMatrix = {};
  let globalMax = -Infinity;
  let globalMin = Infinity;
  let globalSum = 0;
  let globalCount = 0;
  let isAnyDiagonal = false;
  const diagonalMatrices = [];

  for (const key of keys) {
    const matrix = matrixMap[key];
    const stats = calculateSingleMatrixStats(matrix);
    byMatrix[key] = stats;

    if (stats.maxValue > globalMax) globalMax = stats.maxValue;
    if (stats.minValue < globalMin) globalMin = stats.minValue;
    globalSum += stats.totalSum;
    globalCount += stats.totalElements;

    if (stats.isDiagonal) {
      isAnyDiagonal = true;
      diagonalMatrices.push(key);
    }
  }

  const globalAvg = globalCount > 0 ? globalSum / globalCount : 0;

  return {
    global: {
      maxValue: Number(globalMax.toFixed(6)),
      minValue: Number(globalMin.toFixed(6)),
      average: Number(globalAvg.toFixed(6)),
      totalSum: Number(globalSum.toFixed(6)),
      totalElements: globalCount,
      isAnyDiagonal,
      diagonalMatrices
    },
    byMatrix
  };
}

module.exports = {
  validateMatrix,
  isDiagonalMatrix,
  calculateSingleMatrixStats,
  processStats
};
