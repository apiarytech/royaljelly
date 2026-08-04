package selection

import (
	"testing"

	. "github.com/apiarytech/royaljelly/core"
)

func TestSEL(t *testing.T) {
	t.Run("LINT", func(t *testing.T) {
		if res := SEL(false, LINT(10), LINT(20)); res != LINT(10) {
			t.Errorf("SEL(false) = %v; want 10", res)
		}
		if res := SEL(true, LINT(10), LINT(20)); res != LINT(20) {
			t.Errorf("SEL(true) = %v; want 20", res)
		}
	})
	t.Run("REAL", func(t *testing.T) {
		if res := SEL(true, REAL(1.5), REAL(2.5)); res != REAL(2.5) {
			t.Errorf("SEL(true) = %v; want 2.5", res)
		}
	})
	t.Run("STRING", func(t *testing.T) {
		if res := SEL(false, STRING("a"), STRING("b")); res != STRING("a") {
			t.Errorf("SEL(false) = %v; want 'a'", res)
		}
	})
}

func TestMAX(t *testing.T) {
	t.Run("LINTs", func(t *testing.T) {
		res, err := MAX(LINT(10), LINT(50), LINT(20))
		if err != nil || res != LINT(50) {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, LINT(50))
		}
	})
	t.Run("REALs", func(t *testing.T) {
		res, err := MAX(REAL(10.5), REAL(10.6), REAL(10.1))
		if err != nil || res != REAL(10.6) {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, REAL(10.6))
		}
	})
	t.Run("Strings", func(t *testing.T) {
		res, err := MAX(STRING("apple"), STRING("orange"), STRING("banana"))
		if err != nil || res != STRING("orange") {
			t.Errorf("MAX() = %v, err: %v; want %v, nil", res, err, STRING("orange"))
		}
	})
	t.Run("Not enough inputs", func(t *testing.T) {
		_, err := MAX(LINT(10))
		if err == nil {
			t.Error("MAX() with one input should return an error")
		}
	})
}

func TestMIN(t *testing.T) {
	t.Run("LINTs", func(t *testing.T) {
		res, err := MIN(LINT(10), LINT(50), LINT(20))
		if err != nil || res != LINT(10) {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, LINT(10))
		}
	})
	t.Run("REALs", func(t *testing.T) {
		res, err := MIN(REAL(10.5), REAL(10.6), REAL(10.1))
		if err != nil || res != REAL(10.1) {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, REAL(10.1))
		}
	})
	t.Run("Strings", func(t *testing.T) {
		res, err := MIN(STRING("apple"), STRING("orange"), STRING("banana"))
		if err != nil || res != STRING("apple") {
			t.Errorf("MIN() = %v, err: %v; want %v, nil", res, err, STRING("apple"))
		}
	})
	t.Run("Not enough inputs", func(t *testing.T) {
		_, err := MIN(LINT(10))
		if err == nil {
			t.Error("MIN() with one input should return an error")
		}
	})
}

func TestLIMIT(t *testing.T) {
	t.Run("LINT", func(t *testing.T) {
		if res := LIMIT(LINT(10), LINT(50), LINT(100)); res != LINT(50) {
			t.Errorf("LIMIT(within) = %v; want 50", res)
		}
		if res := LIMIT(LINT(10), LINT(5), LINT(100)); res != LINT(10) {
			t.Errorf("LIMIT(below) = %v; want 10", res)
		}
		if res := LIMIT(LINT(10), LINT(150), LINT(100)); res != LINT(100) {
			t.Errorf("LIMIT(above) = %v; want 100", res)
		}
	})
	t.Run("REAL", func(t *testing.T) {
		if res := LIMIT(REAL(10.0), REAL(50.5), REAL(100.0)); res != REAL(50.5) {
			t.Errorf("LIMIT(REAL) = %v; want 50.5", res)
		}
	})
	t.Run("STRING", func(t *testing.T) {
		if res := LIMIT(STRING("a"), STRING("b"), STRING("c")); res != STRING("b") {
			t.Errorf("LIMIT(STRING) = %v; want 'b'", res)
		}
	})
}

func TestMUX(t *testing.T) {
	t.Run("Select STRING", func(t *testing.T) {
		res, err := MUX(LINT(1), STRING("a"), STRING("b"), STRING("c"))
		if err != nil || res != STRING("b") {
			t.Errorf("MUX() = %v, err: %v; want 'b', nil", res, err)
		}
	})
	t.Run("Select REAL", func(t *testing.T) {
		res, err := MUX(INT(0), REAL(10.0), REAL(20.0))
		if err != nil || res != REAL(10.0) {
			t.Errorf("MUX() = %v, err: %v; want 10.0, nil", res, err)
		}
	})
	t.Run("Selector out of bounds (negative)", func(t *testing.T) {
		_, err := MUX(LINT(-1), STRING("a"))
		if err == nil {
			t.Error("MUX() with negative selector should return an error")
		}
	})
	t.Run("Selector out of bounds (too high)", func(t *testing.T) {
		_, err := MUX(LINT(2), STRING("a"), STRING("b"))
		if err == nil {
			t.Error("MUX() with out-of-bounds selector should return an error")
		}
	})
	t.Run("No options", func(t *testing.T) {
		_, err := MUX[LINT, STRING](LINT(0))
		if err == nil {
			t.Error("MUX() with no options should return an error")
		}
	})
}
