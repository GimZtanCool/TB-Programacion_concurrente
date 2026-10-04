package main

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testHeader = "step,type,amount,oldbalanceOrg,newbalanceOrig,oldbalanceDest,newbalanceDest,nameDest,isFraud\n"

func TestParseRowValidation(t *testing.T) {
	idx := map[string]int{}
	for i, name := range strings.Split(strings.TrimSpace(testHeader), ",") {
		idx[name] = i
	}
	valid := []string{"1", "TRANSFER", "100", "100", "0", "0", "100", "C123", "1"}
	if _, err := parseRow(valid, idx); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ col, value string }{
		{"amount", "NaN"}, {"amount", "+Inf"}, {"oldbalanceOrg", "-Inf"},
		{"amount", "-1"}, {"step", "1.5"}, {"type", "OTHER"},
		{"nameDest", "X123"}, {"nameDest", ""}, {"isFraud", "2"},
	} {
		t.Run(tc.col+tc.value, func(t *testing.T) {
			row := append([]string(nil), valid...)
			row[idx[tc.col]] = tc.value
			if _, err := parseRow(row, idx); err == nil {
				t.Fatal("invalid row accepted")
			}
		})
	}
	if _, err := parseRow(valid[:3], idx); err == nil {
		t.Fatal("short row accepted")
	}
}

func TestDatasetValidation(t *testing.T) {
	for _, tc := range []struct{ name, content, want string }{
		{"missing", "step,type\n1,TRANSFER\n", "required CSV column"},
		{"duplicate", testHeader[:len(testHeader)-1] + ",step\n", "duplicate CSV column"},
		{"unordered", testHeader + "2,TRANSFER,10,10,0,0,10,C1,1\n1,PAYMENT,1,1,0,0,1,M1,0\n", "not chronological"},
		{"empty", testHeader, "at least two"},
		{"singleclass", testHeader + "1,PAYMENT,1,1,0,0,1,M1,0\n2,PAYMENT,2,2,0,0,2,M2,0\n", "both classes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "data.csv")
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := loadDataset(path)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want %s", err, tc.want)
			}
		})
	}
}

func TestConcurrentEquivalence(t *testing.T) {
	const rows = 17
	x := make([]float32, rows*featureCount)
	y := make([]uint8, rows)
	for i := range y {
		if i%3 == 0 {
			y[i] = 1
		}
		for j := 0; j < featureCount; j++ {
			x[i*featureCount+j] = float32((i*7+j)%11-5) / 10
		}
	}
	w := model{bias: 0.15}
	for j := range w.weights {
		w.weights[j] = float64(j-7) / 100
	}
	want := sequentialGradient(x, y, w)
	closeEnough := func(a, b float64) bool { return math.Abs(a-b) <= 1e-10*math.Max(1, math.Abs(a)) }
	for _, workers := range []int{1, 2, 3, 4, 20} {
		for repeat := 0; repeat < 10; repeat++ {
			got := parallelGradient(x, y, w, workers)
			if !closeEnough(want.bias, got.bias) {
				t.Fatalf("bias differs with %d workers", workers)
			}
			for j := range got.values {
				if !closeEnough(want.values[j], got.values[j]) {
					t.Fatalf("gradient %d differs", j)
				}
			}
			seq, par := train(x, y, 1, 5), train(x, y, workers, 5)
			if !closeEnough(seq.bias, par.bias) {
				t.Fatal("trained bias differs")
			}
			for j := range seq.weights {
				if !closeEnough(seq.weights[j], par.weights[j]) {
					t.Fatalf("trained weight %d differs", j)
				}
			}
		}
	}
}

func TestWorkersAndTrimmedMean(t *testing.T) {
	for _, s := range []string{"0", "-1", "2,x", "2,"} {
		if _, err := parseWorkers(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
	w, err := parseWorkers("4,2,2,1")
	if err != nil || len(w) != 3 || w[0] != 1 || w[1] != 2 || w[2] != 4 {
		t.Fatalf("unexpected workers %v %v", w, err)
	}
	if got := trimmedMean([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 100}, .1); got != 5.5 {
		t.Fatalf("mean = %v", got)
	}
}

func TestOutputProtectsDataset(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input.csv")
	if err := os.WriteFile(input, []byte(testHeader), 0600); err != nil {
		t.Fatal(err)
	}
	if err := validateOutputPath(input, filepath.Join(dir, ".", "input.csv")); err == nil {
		t.Fatal("source overwrite accepted")
	}
	if err := validateOutputPath(input, filepath.Join(dir, "results.csv")); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(dir, "alias.csv")
	if err := os.Link(input, alias); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}
	if err := validateOutputPath(input, alias); err == nil {
		t.Fatal("source overwrite through alias accepted")
	}
}

// Opt-in integration check: the full CSV is large and is not needed by unit tests.
func TestFullPaySim(t *testing.T) {
	path := os.Getenv("TP_FULL_DATASET")
	if path == "" {
		t.Skip("set TP_FULL_DATASET to the original PaySim CSV")
	}
	d, err := loadDataset(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.trainY) != 5090096 || len(d.testY) != 1272524 || countOnes(d.trainY) != 3959 || countOnes(d.testY) != 4254 {
		t.Fatal("unexpected full dataset split")
	}
	if math.Abs(d.cap-22919326.54) > 0.01 {
		t.Fatalf("unexpected cap %.6f", d.cap)
	}
	for _, workers := range []int{1, 12} {
		w := train(d.trainX, d.trainY, workers, 5)
		tp, fp, fn := 0, 0, 0
		for i, y := range d.testY {
			z := w.bias
			for j := 0; j < featureCount; j++ {
				z += w.weights[j] * float64(d.testX[i*featureCount+j])
			}
			pred := sigmoid(z) >= 0.5
			if pred && y == 1 {
				tp++
			} else if pred && y == 0 {
				fp++
			} else if !pred && y == 1 {
				fn++
			}
		}
		t.Logf("workers=%d epochs=5 TP=%d FP=%d FN=%d", workers, tp, fp, fn)
		if tp != 3731 || fp != 285109 || fn != 523 {
			t.Fatal("full evaluation changed from PC2")
		}
	}
}
