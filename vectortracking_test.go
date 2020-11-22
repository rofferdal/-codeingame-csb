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

func TestNewSmartVectorPolar(t *testing.T) {
	t.Run("shouldGeneratePolar_0_0", func(t *testing.T) {
		sv := NewSmartVectorPolar(0, 0)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_0", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 0)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", 1, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_90", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 90)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", 1, sv.y)
	})
	t.Run("shouldGeneratePolar_1_180", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, 180)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", -1, sv.x)
		assertFloatPrettyEqual(t, "y value", 0, sv.y)
	})
	t.Run("shouldGeneratePolar_1_-90", func(t *testing.T) {
		sv := NewSmartVectorPolar(1, -90)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", 0, sv.x)
		assertFloatPrettyEqual(t, "y value", -1, sv.y)
	})
	t.Run("shouldGeneratePolar_2_45", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, 45)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2), sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2), sv.y)
	})
	t.Run("shouldGeneratePolar_2_-45", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, -45)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2), sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2) * -1, sv.y)
	})
	t.Run("shouldGeneratePolar_2_135", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, 135)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2) * -1, sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2), sv.y)
	})
	t.Run("shouldGeneratePolar_2_-135", func(t *testing.T) {
		sv := NewSmartVectorPolar(2, -135)
		fmt.Printf("x: %f, y: %f, angleDegrees: %f, angleRadians: %f\n", sv.x, sv.y, sv.angleDegrees, sv.angleRadians)
		assertFloatPrettyEqual(t, "x value", math.Sqrt(2) * -1, sv.x)
		assertFloatPrettyEqual(t, "y value", math.Sqrt(2) * -1, sv.y)
	})
}


func TestNewSmartVectorCartesian(t *testing.T) {
	t.Run("shouldGenerate_0_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, 0)
		assertTrue(t,"Length is one", sv.length == 0)
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is NaN", math.IsNaN(sv.angleDegrees))
	})
	t.Run("shouldGenerate_1_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, 0)
		assertTrue(t,"Length is one", sv.length == 1)
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is 0 degrees", sv.angleDegrees == 0)
	})
	t.Run("shouldGenerate_0_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, 1)
		assertTrue(t,"Length is one", sv.length == 1)
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is 90 degrees", sv.angleDegrees == 90)
	})
	t.Run("shouldGenerate_-1_0", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, 0)
		assertTrue(t,"Length is one", sv.length == 1)
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is 180 degrees", sv.angleDegrees == 180)
	})
	t.Run("shouldGenerate_0_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(0, -1)
		assertTrue(t,"Length is one", sv.length == 1)
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is -90 degrees", sv.angleDegrees == -90)
	})
	t.Run("shouldGenerate_1_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, 1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is 45 degrees", sv.angleDegrees == 45)
	})
	t.Run("shouldGenerate_1_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(1, -1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is -45 degrees", sv.angleDegrees == -45)
	})
	t.Run("shouldGenerate_-1_1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, 1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is 135 degrees", sv.angleDegrees == 135)
	})
	t.Run("shouldGenerate_-1_-1", func(t *testing.T) {
		sv := NewSmartVectorCartesian(-1, -1)
		assertTrue(t,"Length is one", sv.length == math.Sqrt(2))
		fmt.Printf("angleDegrees: %f, angleRadians: %f\n", sv.angleDegrees, sv.angleRadians)
		assertTrue(t,"Angle is -135 degrees", sv.angleDegrees == -135)
	})
}