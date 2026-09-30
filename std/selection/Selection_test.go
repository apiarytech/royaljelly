package selection

import (
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

func TestSEL(t *testing.T) {
	t.Run("LINT", func(t *testing.T) {
		if res := SEL(false, iec.LINT(10), iec.LINT(20)); res != iec.LINT(10) {
			t.Errorf("SEL(false) = %v; want 10", res)
		}
		if res := SEL(true, iec.LINT(10), iec.LINT(20)); res != iec.LINT(20) {
			t.Errorf("SEL(true) = %v; want 20", res)
		}
	})
	t.Run("REAL", func(t *testing.T) {
		if res := SEL(true, iec.REAL(1.5), iec.REAL(2.5)); res != iec.REAL(2.5) {
			t.Errorf("SEL(true) = %v; want 2.5", res)
		}
	})
	t.Run("STRING", func(t *testing.T) {
		if res := SEL(false, iec.STRING("a"), iec.STRING("b")); res != iec.STRING("a") {
			t.Errorf("SEL(false) = %v; want 'a'", res)
		}
	})
}

func TestMAX(t *testing.T) {
	t.Run("LINTs", func(t *testing.T) {
		res, err := MAX(iec.LINT(10), iec.LINT(50), iec.LINT(20))
		if err != nil || res != iec.LINT(50) {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, iec.LINT(50))
		}
	})
	t.Run("REALs", func(t *testing.T) {
		res, err := MAX(iec.REAL(10.5), iec.REAL(10.6), iec.REAL(10.1))
		if err != nil || res != iec.REAL(10.6) {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, iec.REAL(10.6))
		}
	})
	t.Run("Strings", func(t *testing.T) {
		res, err := MAX(iec.STRING("apple"), iec.STRING("orange"), iec.STRING("banana"))
		if err != nil || res != iec.STRING("orange") {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, iec.STRING("orange"))
		}
	})
	t.Run("Not enough inputs", func(t *testing.T) {
		_, err := MAX(iec.LINT(10))
		if err == nil {
			t.Error("MAX() with one input should return an error")
		}
	})
}

func TestMIN(t *testing.T) {
	t.Run("LINTs", func(t *testing.T) {
		res, err := MIN(iec.LINT(10), iec.LINT(50), iec.LINT(20))
		if err != nil || res != iec.LINT(10) {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, iec.LINT(10))
		}
	})
	t.Run("REALs", func(t *testing.T) {
		res, err := MIN(iec.REAL(10.5), iec.REAL(10.6), iec.REAL(10.1))
		if err != nil || res != iec.REAL(10.1) {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, iec.REAL(10.1))
		}
	})
	t.Run("Strings", func(t *testing.T) {
		res, err := MIN(iec.STRING("apple"), iec.STRING("orange"), iec.STRING("banana"))
		if err != nil || res != iec.STRING("apple") {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, iec.STRING("apple"))
		}
	})
	t.Run("Not enough inputs", func(t *testing.T) {
		_, err := MIN(iec.LINT(10))
		if err == nil {
			t.Error("MIN() with one input should return an error")
		}
	})
}

func TestLIMIT(t *testing.T) {
	t.Run("LINT", func(t *testing.T) {
		if res := LIMIT(iec.LINT(10), iec.LINT(50), iec.LINT(100)); res != iec.LINT(50) {
			t.Errorf("LIMIT(within) = %v; want 50", res)
		}
		if res := LIMIT(iec.LINT(10), iec.LINT(5), iec.LINT(100)); res != iec.LINT(10) {
			t.Errorf("LIMIT(below) = %v; want 10", res)
		}
		if res := LIMIT(iec.LINT(10), iec.LINT(150), iec.LINT(100)); res != iec.LINT(100) {
			t.Errorf("LIMIT(above) = %v; want 100", res)
		}
	})
	t.Run("REAL", func(t *testing.T) {
		if res := LIMIT(iec.REAL(10.0), iec.REAL(50.5), iec.REAL(100.0)); res != iec.REAL(50.5) {
			t.Errorf("LIMIT(REAL) = %v; want 50.5", res)
		}
	})
	t.Run("STRING", func(t *testing.T) {
		if res := LIMIT(iec.STRING("a"), iec.STRING("b"), iec.STRING("c")); res != iec.STRING("b") {
			t.Errorf("LIMIT(STRING) = %v; want 'b'", res)
		}
	})
}

func TestMUX(t *testing.T) {
	t.Run("Select STRING", func(t *testing.T) {
		res, err := MUX(iec.LINT(1), iec.STRING("a"), iec.STRING("b"), iec.STRING("c"))
		if err != nil || res != iec.STRING("b") {
			t.Errorf("MUX() = %v, err: %v; want 'b', nil", res, err)
		}
	})
	t.Run("Select REAL", func(t *testing.T) {
		res, err := MUX(iec.INT(0), iec.REAL(10.0), iec.REAL(20.0))
		if err != nil || res != iec.REAL(10.0) {
			t.Errorf("MUX() = %v, err: %v; want 10.0, nil", res, err)
		}
	})
	t.Run("Selector out of bounds (negative)", func(t *testing.T) {
		_, err := MUX(iec.LINT(-1), iec.STRING("a"))
		if err == nil {
			t.Error("MUX() with negative selector should return an error")
		}
	})
	t.Run("Selector out of bounds (too high)", func(t *testing.T) {
		_, err := MUX(iec.LINT(2), iec.STRING("a"), iec.STRING("b"))
		if err == nil {
			t.Error("MUX() with out-of-bounds selector should return an error")
		}
	})
	t.Run("No options", func(t *testing.T) {
		_, err := MUX[iec.LINT, iec.STRING](iec.LINT(0))
		if err == nil {
			t.Error("MUX() with no options should return an error")
		}
	})
}
