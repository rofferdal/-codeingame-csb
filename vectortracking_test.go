package main

import (
	"fmt"
	"math"
	"testing"
)

func assertTrue(t *testing.T, what string, condition bool) {
	if !condition {
		t.Errorf("%s failed, should be true", what)
	}
}

func assertFloatPrettyEqual(t *testing.T, what string, expected float64, actual float64) {
	if math.Abs(expected - actual) > 0.00000001   {
		t.Errorf("%s failed. Expected %f, got %f", what, expected, actual)
	}
}

func TestSubtractVector(t *testing.T) {
	t.Run("NullVMinusNullV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(0, 0)
		sv2 := NewSmartVectorCartesian(0, 0)
		resV := sv1.subtractVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, resV.x)
		assertFloatPrettyEqual(t, "y value", 0, resV.y)
		assertFloatPrettyEqual(t, "length", 0, resV.length)
		assertTrue(t,"angleDegrees is NaN", math.IsNaN(resV.angleDegrees))
		assertTrue(t,"angleRadians is NaN", math.IsNaN(resV.angleRadians))
	})
	t.Run("NullVMinusNormV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(0, 0)
		sv2 := NewSmartVectorCartesian(10, 10)
		resV := sv1.subtractVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", -10, resV.x)
		assertFloatPrettyEqual(t, "y value", -10, resV.y)
		assertFloatPrettyEqual(t, "length", sv2.length, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", sv2.angleDegrees - 180, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", sv2.angleRadians - math.Pi, resV.angleRadians)
	})
	t.Run("NegVMinusNormV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(-10, -10)
		sv2 := NewSmartVectorCartesian(10, 10)
		resV := sv1.subtractVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", -20, resV.x)
		assertFloatPrettyEqual(t, "y value", -20, resV.y)
		assertFloatPrettyEqual(t, "length", math.Sqrt(2) * 20, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", -135, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", math.Pi * -0.75, resV.angleRadians)
	})
	t.Run("NegVMinusNegV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(-5, -10)
		sv2 := NewSmartVectorCartesian(-10, -5)
		resV := sv1.subtractVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 5, resV.x)
		assertFloatPrettyEqual(t, "y value", -5, resV.y)
		assertFloatPrettyEqual(t, "length", math.Sqrt(2) * 5, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", -45, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", math.Pi * -0.25, resV.angleRadians)
	})
}

func TestAddVector(t *testing.T) {
	t.Run("NullVPlusNullV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(0, 0)
		sv2 := NewSmartVectorCartesian(0, 0)
		resV := sv1.addVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, resV.x)
		assertFloatPrettyEqual(t, "y value", 0, resV.y)
		assertFloatPrettyEqual(t, "length", 0, resV.length)
		assertTrue(t,"angleDegrees is NaN", math.IsNaN(resV.angleDegrees))
		assertTrue(t,"angleRadians is NaN", math.IsNaN(resV.angleRadians))
	})
	t.Run("NullVPlusNormV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(0, 0)
		sv2 := NewSmartVectorCartesian(10, 10)
		resV := sv1.addVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 10, resV.x)
		assertFloatPrettyEqual(t, "y value", 10, resV.y)
		assertFloatPrettyEqual(t, "length", sv2.length, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", sv2.angleDegrees, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", sv2.angleRadians, resV.angleRadians)
	})
	t.Run("NegVPlusNormV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(-10, -10)
		sv2 := NewSmartVectorCartesian(10, 10)
		resV := sv1.addVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, resV.x)
		assertFloatPrettyEqual(t, "y value", 0, resV.y)
		assertFloatPrettyEqual(t, "length", 0, resV.length)
		assertTrue(t,"angleDegrees is NaN", math.IsNaN(resV.angleDegrees))
		assertTrue(t,"angleRadians is NaN", math.IsNaN(resV.angleRadians))
	})
	t.Run("NegVPlusNegV", func(t *testing.T) {
		sv1 := NewSmartVectorCartesian(-5, -10)
		sv2 := NewSmartVectorCartesian(-10, -5)
		resV := sv1.addVector(sv2)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", -15, resV.x)
		assertFloatPrettyEqual(t, "y value", -15, resV.y)
		assertFloatPrettyEqual(t, "length", math.Sqrt(2) * 15, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", -135, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", math.Pi * -0.75, resV.angleRadians)
	})
}

func TestMultiplyByNumber(t *testing.T) {
	t.Run("NullVX5", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, 0)
		resV := sv.multiplyNumber(5)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, resV.x)
		assertFloatPrettyEqual(t, "y value", 0, resV.y)
		assertFloatPrettyEqual(t, "length", 0, resV.length)
		assertTrue(t,"angleDegrees is NaN", math.IsNaN(resV.angleDegrees))
		assertTrue(t,"angleRadians is NaN", math.IsNaN(resV.angleRadians))
	})
	t.Run("NormVX0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-10, 10)
		resV := sv.multiplyNumber(0)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, resV.x)
		assertFloatPrettyEqual(t, "y value", 0, resV.y)
		assertFloatPrettyEqual(t, "length", 0, resV.length)
		assertTrue(t,"angleDegrees is NaN", math.IsNaN(resV.angleDegrees))
		assertTrue(t,"angleRadians is NaN", math.IsNaN(resV.angleRadians))
	})
	t.Run("NormVX10", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-10, 10)
		resV := sv.multiplyNumber(10)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", -100, resV.x)
		assertFloatPrettyEqual(t, "y value", 100, resV.y)
		assertFloatPrettyEqual(t, "length", sv.length * 10, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", sv.angleDegrees, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", sv.angleRadians, resV.angleRadians)
	})
	t.Run("NormVX-10", func(t *testing.T) {
		sv := NewSmartVectorCartesian(10, 10)
		resV := sv.multiplyNumber(-10)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", -100, resV.x)
		assertFloatPrettyEqual(t, "y value", -100, resV.y)
		assertFloatPrettyEqual(t, "length", sv.length * 10, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", sv.angleDegrees - 180, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", sv.angleRadians - math.Pi, resV.angleRadians)
	})
	t.Run("NegVX-10", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-10, -10)
		resV := sv.multiplyNumber(-10)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", resV.x, resV.y, resV.angleDegrees, resV.angleRadians)
		assertFloatPrettyEqual(t, "x value", 100, resV.x)
		assertFloatPrettyEqual(t, "y value", 100, resV.y)
		assertFloatPrettyEqual(t, "length", sv.length * 10, resV.length)
		assertFloatPrettyEqual(t, "angleDegrees", sv.angleDegrees + 180, resV.angleDegrees)
		assertFloatPrettyEqual(t, "angleRadians", sv.angleRadians + math.Pi, resV.angleRadians)
	})
}

func TestNewSmartVectorPolar(t *testing.T) {
	t.Run("shouldGeneratePolar_0_0", func(t *testing.T) {
		sv := NewSmartVectorPolar(0, 0)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_0", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 0)
		assertFloatPrettyEqual(t, "x value", 1, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_90", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 90)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", 1, sv.y)
	})
	t.Run("shouldGeneratePolar_1_180", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 180)
		assertFloatPrettyEqual(t, "x value", -1, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_-90", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, -90)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", -1, sv.y)
	})
	t.Run("shouldGeneratePolar_2_45", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, 45)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2), sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2), sv.y)
	})
	t.Run("shouldGeneratePolar_2_-45", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, -45)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2), sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2) * -1, sv.y)
	})
	t.Run("shouldGeneratePolar_2_135", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, 135)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2) * -1, sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2), sv.y)
	})
	t.Run("shouldGeneratePolar_2_-135", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, -135)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2) * -1, sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2) * -1, sv.y)
	})
}


func TestNewSmartVectorCartesian(t *testing.T) {
	t.Run("shouldGenerate_0_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, 0)
		assertTrue(t,"Length is one", sv.length == 0)
		assertTrue(t,"Angle is NaN", math.IsNaN(sv.angleDegrees))
	})
	t.Run("shouldGenerate_1_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, 0)
		assertTrue(t,"Length is one", sv.length == 1)
		assertTrue(t,"Angle is 0 degrees", sv.angleDegrees == 0)
	})
	t.Run("shouldGenerate_0_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, 1)
		assertTrue(t,"Length is one", sv.length == 1)
		assertTrue(t,"Angle is 90 degrees", sv.angleDegrees == 90)
	})
	t.Run("shouldGenerate_-1_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, 0)
		assertTrue(t,"Length is one", sv.length == 1)
		assertTrue(t,"Angle is 180 degrees", sv.angleDegrees == 180)
	})
	t.Run("shouldGenerate_0_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, -1)
		assertTrue(t,"Length is one", sv.length == 1)
		assertTrue(t,"Angle is -90 degrees", sv.angleDegrees == -90)
	})
	t.Run("shouldGenerate_1_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, 1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		assertTrue(t,"Angle is 45 degrees", sv.angleDegrees == 45)
	})
	t.Run("shouldGenerate_1_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, -1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		assertTrue(t,"Angle is -45 degrees", sv.angleDegrees == -45)
	})
	t.Run("shouldGenerate_-1_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, 1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		assertTrue(t,"Angle is 135 degrees", sv.angleDegrees == 135)
	})
	t.Run("shouldGenerate_-1_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, -1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		assertTrue(t,"Angle is -135 degrees", sv.angleDegrees == -135)
	})
}