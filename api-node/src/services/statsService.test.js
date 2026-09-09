const { describe, it } = require('node:test');
const assert = require('node:assert');
const {
  validateMatrix,
  isDiagonalMatrix,
  calculateSingleMatrixStats,
  processStats
} = require('./statsService');

describe('Stats Service - Validaciones y Cálculos', () => {

  describe('validateMatrix', () => {
    it('debe validar exitosamente una matriz regular', () => {
      const matrix = [
        [1, 2, 3],
        [4, 5, 6]
      ];
      const res = validateMatrix(matrix);
      assert.strictEqual(res.valid, true);
      assert.strictEqual(res.rows, 2);
      assert.strictEqual(res.cols, 3);
    });

    it('debe fallar ante una matriz vacía', () => {
      const res = validateMatrix([]);
      assert.strictEqual(res.valid, false);
    });

    it('debe fallar ante una matriz irregular (jagged)', () => {
      const matrix = [
        [1, 2, 3],
        [4, 5]
      ];
      const res = validateMatrix(matrix);
      assert.strictEqual(res.valid, false);
      assert.match(res.error, /Matriz irregular/);
    });

    it('debe fallar ante elementos no numéricos', () => {
      const matrix = [
        [1, 'dos'],
        [3, 4]
      ];
      const res = validateMatrix(matrix);
      assert.strictEqual(res.valid, false);
      assert.match(res.error, /Elemento no numérico/);
    });
  });

  describe('isDiagonalMatrix', () => {
    it('debe retornar true para una matriz identidad', () => {
      const identity = [
        [1, 0, 0],
        [0, 1, 0],
        [0, 0, 1]
      ];
      assert.strictEqual(isDiagonalMatrix(identity), true);
    });

    it('debe retornar true para una matriz diagonal con números arbitrarios', () => {
      const diag = [
        [7.5, 0],
        [0, -3.2]
      ];
      assert.strictEqual(isDiagonalMatrix(diag), true);
    });

    it('debe retornar true considerando tolerancias de punto flotante (epsilon)', () => {
      const diag = [
        [4, 0.00000001],
        [0.00000005, 9]
      ];
      assert.strictEqual(isDiagonalMatrix(diag), true);
    });

    it('debe retornar false si contiene elementos no nulos fuera de la diagonal', () => {
      const notDiag = [
        [1, 2],
        [0, 3]
      ];
      assert.strictEqual(isDiagonalMatrix(notDiag), false);
    });

    it('debe retornar false para matrices rectangulares no cuadradas', () => {
      const rect = [
        [1, 0, 0],
        [0, 1, 0]
      ];
      assert.strictEqual(isDiagonalMatrix(rect), false);
    });
  });

  describe('calculateSingleMatrixStats', () => {
    it('debe calcular correctamente valor máximo, mínimo, suma y promedio', () => {
      const matrix = [
        [1, 2, 3],
        [4, 5, 6]
      ];
      // sum = 21, count = 6, avg = 3.5, min = 1, max = 6
      const stats = calculateSingleMatrixStats(matrix);
      assert.strictEqual(stats.maxValue, 6);
      assert.strictEqual(stats.minValue, 1);
      assert.strictEqual(stats.totalSum, 21);
      assert.strictEqual(stats.average, 3.5);
      assert.strictEqual(stats.totalElements, 6);
      assert.strictEqual(stats.isDiagonal, false);
    });
  });

  describe('processStats (múltiples matrices)', () => {
    it('debe calcular estadísticas globales y por matriz correctamente', () => {
      const input = {
        matrices: {
          matA: [
            [1, 0],
            [0, 2]
          ], // diagonal, min 0, max 2, sum 3, avg 0.75
          matB: [
            [10, 20],
            [30, 40]
          ] // no diagonal, min 10, max 40, sum 100, avg 25
        }
      };

      const result = processStats(input);

      // Global
      assert.strictEqual(result.global.maxValue, 40);
      assert.strictEqual(result.global.minValue, 0);
      assert.strictEqual(result.global.totalSum, 103);
      assert.strictEqual(result.global.totalElements, 8);
      assert.strictEqual(result.global.average, 12.875);
      assert.strictEqual(result.global.isAnyDiagonal, true);
      assert.deepStrictEqual(result.global.diagonalMatrices, ['matA']);

      // Por matriz
      assert.strictEqual(result.byMatrix.matA.isDiagonal, true);
      assert.strictEqual(result.byMatrix.matB.isDiagonal, false);
      assert.strictEqual(result.byMatrix.matB.maxValue, 40);
    });
  });
});
