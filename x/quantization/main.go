// !LLM
package main

import (
	"fmt"
	"math"
)

// QuantizationParams holds the parameters needed for quantization/dequantization
type QuantizationParams struct {
	Scale     float32 // Scale factor for mapping float range to int range
	ZeroPoint int8    // Zero point offset
	MinVal    float32 // Original minimum value
	MaxVal    float32 // Original maximum value
}

// Matrix represents a simple 2D matrix
type Matrix struct {
	Data [][]float32
	Rows int
	Cols int
}

// QuantizedMatrix represents a quantized version using int8
type QuantizedMatrix struct {
	Data   [][]int8
	Rows   int
	Cols   int
	Params QuantizationParams
}

// NewMatrix creates a new matrix with the given dimensions
func NewMatrix(rows, cols int) *Matrix {
	data := make([][]float32, rows)
	for i := range data {
		data[i] = make([]float32, cols)
	}
	return &Matrix{Data: data, Rows: rows, Cols: cols}
}

// calculateQuantizationParams determines the scale and zero point for quantization
func calculateQuantizationParams(minVal, maxVal float32) QuantizationParams {
	// For int8: range is -128 to 127 (256 possible values)
	const qMin, qMax int8 = -128, 127
	const qRange = float32(256) // 256 possible values in int8

	// Calculate scale: how much each quantized unit represents in float space
	scale := (maxVal - minVal) / qRange

	// Calculate zero point: where 0.0 maps to in quantized space
	zeroPoint := int8(math.Round(float64(float32(qMin) - minVal/scale)))

	// Clamp zero point to valid range
	if zeroPoint < qMin {
		zeroPoint = qMin
	}
	if zeroPoint > qMax {
		zeroPoint = qMax
	}

	return QuantizationParams{
		Scale:     scale,
		ZeroPoint: zeroPoint,
		MinVal:    minVal,
		MaxVal:    maxVal,
	}
}

// findMinMax finds the minimum and maximum values in the matrix
func (m *Matrix) findMinMax() (float32, float32) {
	if m.Rows == 0 || m.Cols == 0 {
		return 0, 0
	}

	minVal := m.Data[0][0]
	maxVal := m.Data[0][0]

	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			if m.Data[i][j] < minVal {
				minVal = m.Data[i][j]
			}
			if m.Data[i][j] > maxVal {
				maxVal = m.Data[i][j]
			}
		}
	}

	return minVal, maxVal
}

// Quantize converts the float32 matrix to int8 using calculated parameters
func (m *Matrix) Quantize() *QuantizedMatrix {
	minVal, maxVal := m.findMinMax()
	params := calculateQuantizationParams(minVal, maxVal)

	// Create quantized matrix
	qData := make([][]int8, m.Rows)
	for i := range qData {
		qData[i] = make([]int8, m.Cols)
	}

	// Quantization formula: q = round((x - zero_point_float) / scale) + zero_point_int
	// Simplified: q = round(x / scale) + zero_point
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			// Map float value to quantized range
			quantized := math.Round(float64(m.Data[i][j]/params.Scale)) + float64(params.ZeroPoint)

			// Clamp to int8 range
			if quantized < -128 {
				quantized = -128
			}
			if quantized > 127 {
				quantized = 127
			}

			qData[i][j] = int8(quantized)
		}
	}

	return &QuantizedMatrix{
		Data:   qData,
		Rows:   m.Rows,
		Cols:   m.Cols,
		Params: params,
	}
}

// Dequantize converts the quantized matrix back to float32
func (qm *QuantizedMatrix) Dequantize() *Matrix {
	data := make([][]float32, qm.Rows)
	for i := range data {
		data[i] = make([]float32, qm.Cols)
	}

	// Dequantization formula: x = scale * (q - zero_point)
	for i := 0; i < qm.Rows; i++ {
		for j := 0; j < qm.Cols; j++ {
			dequantized := qm.Params.Scale * float32(qm.Data[i][j]-qm.Params.ZeroPoint)
			data[i][j] = dequantized
		}
	}

	return &Matrix{Data: data, Rows: qm.Rows, Cols: qm.Cols}
}

// Print displays the matrix values
func (m *Matrix) Print(name string) {
	fmt.Printf("\n%s (%dx%d):\n", name, m.Rows, m.Cols)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			fmt.Printf("%8.4f ", m.Data[i][j])
		}
		fmt.Println()
	}
}

// Print displays the quantized matrix values
func (qm *QuantizedMatrix) Print(name string) {
	fmt.Printf("\n%s (%dx%d) [int8]:\n", name, qm.Rows, qm.Cols)
	for i := 0; i < qm.Rows; i++ {
		for j := 0; j < qm.Cols; j++ {
			fmt.Printf("%4d ", qm.Data[i][j])
		}
		fmt.Println()
	}
	fmt.Printf("Scale: %.6f, ZeroPoint: %d, Range: [%.4f, %.4f]\n",
		qm.Params.Scale, qm.Params.ZeroPoint, qm.Params.MinVal, qm.Params.MaxVal)
}

// calculateMemoryUsage returns memory usage in bytes for comparison
func (m *Matrix) calculateMemoryUsage() int {
	return m.Rows * m.Cols * 4 // 4 bytes per float32
}

func (qm *QuantizedMatrix) calculateMemoryUsage() int {
	return qm.Rows*qm.Cols*1 + 4 + 1 + 4 + 4 // 1 byte per int8 + params
}

func main() {
	fmt.Println("=== LLM Quantization Demonstration ===")
	fmt.Println("This demonstrates the core concept of quantization used in efficient LLM inference")

	// Create a sample weight matrix (like a small part of an LLM layer)
	rows, cols := 4, 5
	original := NewMatrix(rows, cols)

	// Fill with some typical weight values that might appear in an LLM
	weights := [][]float32{
		{0.1234, -0.5678, 0.9012, -0.3456, 0.7890},
		{-0.2345, 0.6789, -0.1012, 0.4567, -0.8901},
		{0.3456, -0.7890, 0.2123, -0.5678, 0.9012},
		{-0.4567, 0.8901, -0.3234, 0.6789, -0.1012},
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			original.Data[i][j] = weights[i][j]
		}
	}

	// Show original matrix
	original.Print("Original FP32 Matrix")

	// Quantize to int8
	quantized := original.Quantize()
	quantized.Print("Quantized INT8 Matrix")

	// Dequantize back to float32
	dequantized := quantized.Dequantize()
	dequantized.Print("Dequantized FP32 Matrix")

	// Calculate and show memory savings
	originalMemory := original.calculateMemoryUsage()
	quantizedMemory := quantized.calculateMemoryUsage()
	savings := float64(originalMemory-quantizedMemory) / float64(originalMemory) * 100

	fmt.Printf("\n=== Memory Usage Comparison ===\n")
	fmt.Printf("Original FP32: %d bytes\n", originalMemory)
	fmt.Printf("Quantized INT8: %d bytes (including params)\n", quantizedMemory)
	fmt.Printf("Memory savings: %.1f%%\n", savings)

	// Calculate quantization error
	fmt.Printf("\n=== Quantization Error Analysis ===\n")
	var totalError float64
	var maxError float64

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			error := math.Abs(float64(original.Data[i][j] - dequantized.Data[i][j]))
			totalError += error
			if error > maxError {
				maxError = error
			}
		}
	}

	avgError := totalError / float64(rows*cols)
	fmt.Printf("Average quantization error: %.6f\n", avgError)
	fmt.Printf("Maximum quantization error: %.6f\n", maxError)

	fmt.Println("\n=== How This Applies to Real LLM Frameworks ===")
	fmt.Println(`
In production LLM quantization frameworks:

1. **Layer-wise Quantization**: Each layer (attention, MLP) may have different
   quantization parameters optimized for its value distribution.

2. **Weight vs Activation Quantization**:
   - Weights: Quantized once during model conversion (static)
   - Activations: Quantized dynamically during inference

3. **Mixed Precision**: Critical layers might stay in FP16/FP32 while others
   use INT8/INT4 to balance accuracy and speed.

4. **Computation Graph Adaptation**:
   - Some ops work directly on quantized values (element-wise add)
   - Matrix multiplication often requires dequantization
   - Special kernels can perform quantized matrix operations
   - Batch norm and layer norm may need special handling

5. **Framework Integration**:
   - PyTorch: torch.quantization with QConfig
   - TensorRT: Automatic quantization-aware training
   - ONNX Runtime: Post-training quantization

The key insight: Quantization trades precision for memory/speed, requiring
careful calibration to maintain model quality while achieving efficiency gains.`)
}
