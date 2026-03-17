package image

import (
	"math"
	"testing"
)

func TestTransform(t *testing.T) {
	tests := []struct {
		name      string
		data      []float64
		transform string
		want      float64
	}{
		{"mean of ints", []float64{1, 2, 3, 4, 5}, "mean", 3.0},
		{"mean single", []float64{42}, "mean", 42.0},
		{"max", []float64{1, 5, 3, 2, 4}, "max", 5.0},
		{"max negative", []float64{-10, -5, -1}, "max", -1.0},
		{"min", []float64{1, 5, 3, 2, 4}, "min", 1.0},
		{"min negative", []float64{-10, -5, -1}, "min", -10.0},
		{"maxabs positive", []float64{1, -5, 3}, "maxabs", 5.0},
		{"maxabs all negative", []float64{-1, -2, -3}, "maxabs", 3.0},
		{"first", []float64{7, 8, 9}, "first", 7.0},
		{"first single", []float64{99}, "first", 99.0},
		{"unknown defaults to first", []float64{10, 20}, "bogus", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Transform(tt.data, tt.transform)
			if got != tt.want {
				t.Errorf("Transform(%v, %q) = %v, want %v", tt.data, tt.transform, got, tt.want)
			}
		})
	}
}

func TestTransformNaN(t *testing.T) {
	data := []float64{math.NaN()}
	transforms := []string{"mean", "max", "min", "maxabs", "first"}
	for _, tr := range transforms {
		t.Run(tr, func(t *testing.T) {
			got := Transform(data, tr)
			if got != 0 {
				t.Errorf("Transform(NaN, %q) = %v, want 0", tr, got)
			}
		})
	}
}

func TestApplyCXmodeComplex(t *testing.T) {
	// Complex data: pairs of (real, imag)
	// 3+4i => magnitude=5, phase=atan2(4,3)
	datain := []float64{3, 4}

	t.Run("Ma", func(t *testing.T) {
		out := ApplyCXmode(datain, "Ma", true)
		if len(out) != 1 {
			t.Fatalf("len = %d, want 1", len(out))
		}
		if math.Abs(out[0]-5.0) > 1e-10 {
			t.Errorf("Ma = %v, want 5.0", out[0])
		}
	})

	t.Run("Ph", func(t *testing.T) {
		out := ApplyCXmode(datain, "Ph", true)
		want := math.Atan2(4, 3)
		if math.Abs(out[0]-want) > 1e-10 {
			t.Errorf("Ph = %v, want %v", out[0], want)
		}
	})

	t.Run("Re", func(t *testing.T) {
		out := ApplyCXmode(datain, "Re", true)
		if out[0] != 3.0 {
			t.Errorf("Re = %v, want 3.0", out[0])
		}
	})

	t.Run("Im", func(t *testing.T) {
		out := ApplyCXmode(datain, "Im", true)
		if out[0] != 4.0 {
			t.Errorf("Im = %v, want 4.0", out[0])
		}
	})
}

func TestApplyCXmodeReal(t *testing.T) {
	datain := []float64{-3, 5, 0}

	t.Run("Ma real", func(t *testing.T) {
		out := ApplyCXmode(datain, "Ma", false)
		if len(out) != 3 {
			t.Fatalf("len = %d, want 3", len(out))
		}
		if out[0] != 3.0 {
			t.Errorf("Ma(-3) = %v, want 3.0", out[0])
		}
		if out[1] != 5.0 {
			t.Errorf("Ma(5) = %v, want 5.0", out[1])
		}
	})

	t.Run("Ph real", func(t *testing.T) {
		out := ApplyCXmode(datain, "Ph", false)
		// atan2(0, -3) = pi
		if math.Abs(out[0]-math.Pi) > 1e-10 {
			t.Errorf("Ph(-3) = %v, want %v", out[0], math.Pi)
		}
		// atan2(0, 5) = 0
		if out[1] != 0.0 {
			t.Errorf("Ph(5) = %v, want 0", out[1])
		}
	})

	t.Run("Re real", func(t *testing.T) {
		out := ApplyCXmode(datain, "Re", false)
		if out[0] != -3.0 || out[1] != 5.0 {
			t.Errorf("Re = %v, want original", out)
		}
	})

	t.Run("Im real", func(t *testing.T) {
		out := ApplyCXmode(datain, "Im", false)
		for i, v := range out {
			if v != 0 {
				t.Errorf("Im[%d] = %v, want 0", i, v)
			}
		}
	})
}

func TestDownSampleLineInY(t *testing.T) {
	// 6 elements, outxsize=3 => 2 lines of 3 elements
	// line0: [1,2,3], line1: [4,5,6]
	// For column 0: transform([1,4], "mean") = 2.5
	// For column 1: transform([2,5], "mean") = 3.5
	// For column 2: transform([3,6], "mean") = 4.5
	data := []float64{1, 2, 3, 4, 5, 6}
	out := DownSampleLineInY(data, 3, "mean")
	if len(out) != 3 {
		t.Fatalf("len = %d, want 3", len(out))
	}
	expected := []float64{2.5, 3.5, 4.5}
	for i, want := range expected {
		if math.Abs(out[i]-want) > 1e-10 {
			t.Errorf("out[%d] = %v, want %v", i, out[i], want)
		}
	}
}

func TestDownSampleLineInYSingleLine(t *testing.T) {
	data := []float64{10, 20, 30}
	out := DownSampleLineInY(data, 3, "first")
	if len(out) != 3 {
		t.Fatalf("len = %d, want 3", len(out))
	}
	for i, want := range data {
		if out[i] != want {
			t.Errorf("out[%d] = %v, want %v", i, out[i], want)
		}
	}
}

func TestCreateOutputRGBA(t *testing.T) {
	data := []float64{0, 0.5, 1.0}
	out := CreateOutput(data, "RGBA", 0, 1.0, "Greyscale")
	// 3 data points * 4 bytes (RGBA) = 12 bytes
	if len(out) != 12 {
		t.Fatalf("len = %d, want 12", len(out))
	}
	// Each pixel should have alpha=255
	for i := 0; i < 3; i++ {
		alpha := out[i*4+3]
		if alpha != 255 {
			t.Errorf("pixel %d alpha = %d, want 255", i, alpha)
		}
	}
}

func TestCreateOutputRGBAEqualZminZmax(t *testing.T) {
	data := []float64{5, 5, 5}
	out := CreateOutput(data, "RGBA", 5, 5, "Greyscale")
	if len(out) != 12 {
		t.Fatalf("len = %d, want 12", len(out))
	}
	// When zmin==zmax, all pixels use colorPalette[0]
	// All pixels should be identical
	for i := 1; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if out[i*4+j] != out[j] {
				t.Errorf("pixel %d byte %d = %d, want %d (same as pixel 0)", i, j, out[i*4+j], out[j])
			}
		}
	}
}
