// PC2 benchmark: weighted logistic regression over the chronological PaySim split.
// The sequential and worker-pool paths use the same features, weights and update rule.
package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	featureCount = 15
	trainRatio   = 0.8
)

var featureNames = [featureCount]string{
	"step", "amount", "oldbalanceOrg", "newbalanceOrig", "oldbalanceDest", "newbalanceDest",
	"errorBalanceOrig", "errorBalanceDest", "hora_del_dia", "dia_de_la_semana",
	"type_CASH_OUT", "type_DEBIT", "type_PAYMENT", "type_TRANSFER", "dest_type_M",
}

type dataset struct {
	trainX []float32
	trainY []uint8
	testX  []float32
	testY  []uint8
	mean   [featureCount]float64
	std    [featureCount]float64
	cap    float64
}

type csvRow struct {
	step, amount, oldOrig, newOrig, oldDest, newDest float64
	typ, destType                                    string
	y                                                uint8
}

func main() {
	path := flag.String("input", "../paysim.csv", "path to the full PaySim CSV")
	mode := flag.String("mode", "bench", "bench or evaluate")
	repeats := flag.Int("repeats", 30, "benchmark repetitions")
	epochs := flag.Int("epochs", 5, "training epochs for evaluation")
	benchEpochs := flag.Int("bench-epochs", 1, "epochs per timed benchmark")
	workersArg := flag.String("workers", "", "comma-separated worker counts; default 1,2,4,NumCPU")
	outPath := flag.String("out", "results.csv", "benchmark CSV output path")
	flag.Parse()

	data, err := loadDataset(*path)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("Dataset: %d train, %d test, %d features; cap amount=%.2f\n", len(data.trainY), len(data.testY), featureCount, data.cap)
	fmt.Printf("Training frauds: %d; test frauds: %d; logical CPUs: %d\n", countOnes(data.trainY), countOnes(data.testY), runtime.NumCPU())
	fmt.Println("Feature order:", strings.Join(featureNames[:], ", "))

	switch *mode {
	case "bench":
		if *repeats < 1 || *benchEpochs < 1 {
			fatal(errors.New("repeats and bench-epochs must be positive"))
		}
		workers, err := parseWorkers(*workersArg)
		if err != nil {
			fatal(err)
		}
		if err := runBench(data, workers, *repeats, *benchEpochs, *outPath); err != nil {
			fatal(err)
		}
	case "evaluate":
		if *epochs < 1 {
			fatal(errors.New("epochs must be positive"))
		}
		if err := evaluate(data, *epochs); err != nil {
			fatal(err)
		}
	default:
		fatal(fmt.Errorf("unknown mode %q (use bench or evaluate)", *mode))
	}
}

func parseWorkers(s string) ([]int, error) {
	if strings.TrimSpace(s) == "" {
		seen := map[int]bool{}
		out := []int{1}
		seen[1] = true
		for _, n := range []int{2, 4, runtime.NumCPU()} {
			if n > 1 && !seen[n] {
				out = append(out, n)
				seen[n] = true
			}
		}
		return out, nil
	}
	seen := map[int]bool{}
	var out []int
	for _, field := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil || n < 1 {
			return nil, fmt.Errorf("invalid worker count %q", field)
		}
		if !seen[n] {
			out = append(out, n)
			seen[n] = true
		}
	}
	sort.Ints(out)
	return out, nil
}

func openCSV(path string) (*csv.Reader, *os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(f)
	r.ReuseRecord = true
	if _, err := r.Read(); err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("read CSV header: %w", err)
	}
	return r, f, nil
}

func columnIndex(header []string, name string) int {
	for i, v := range header {
		if v == name {
			return i
		}
	}
	return -1
}

func readHeader(path string) ([]string, map[string]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	h, err := r.Read()
	if err != nil {
		return nil, nil, err
	}
	idx := make(map[string]int, len(h))
	for i, name := range h {
		idx[name] = i
	}
	for _, required := range []string{"step", "type", "amount", "oldbalanceOrg", "newbalanceOrig", "oldbalanceDest", "newbalanceDest", "nameDest", "isFraud"} {
		if _, ok := idx[required]; !ok {
			return nil, nil, fmt.Errorf("required CSV column %q is missing", required)
		}
	}
	return h, idx, nil
}

func parseRow(rec []string, idx map[string]int) (csvRow, error) {
	getFloat := func(k string) (float64, error) { return strconv.ParseFloat(rec[idx[k]], 64) }
	var x csvRow
	var err error
	if x.step, err = getFloat("step"); err != nil {
		return x, err
	}
	if x.amount, err = getFloat("amount"); err != nil {
		return x, err
	}
	if x.oldOrig, err = getFloat("oldbalanceOrg"); err != nil {
		return x, err
	}
	if x.newOrig, err = getFloat("newbalanceOrig"); err != nil {
		return x, err
	}
	if x.oldDest, err = getFloat("oldbalanceDest"); err != nil {
		return x, err
	}
	if x.newDest, err = getFloat("newbalanceDest"); err != nil {
		return x, err
	}
	label, err := strconv.Atoi(rec[idx["isFraud"]])
	if err != nil || (label != 0 && label != 1) {
		return x, fmt.Errorf("invalid isFraud value %q", rec[idx["isFraud"]])
	}
	x.y = uint8(label)
	x.typ = rec[idx["type"]]
	nameDest := rec[idx["nameDest"]]
	if len(nameDest) == 0 {
		return x, errors.New("empty nameDest")
	}
	x.destType = nameDest[:1]
	return x, nil
}

func (r csvRow) features(cap float64) [featureCount]float64 {
	amount := math.Min(r.amount, cap)
	return [featureCount]float64{
		r.step, amount, r.oldOrig, r.newOrig, r.oldDest, r.newDest,
		r.newOrig + amount - r.oldOrig,
		r.oldDest + amount - r.newDest,
		math.Mod(r.step, 24), math.Mod(math.Floor(r.step/24), 7),
		boolFloat(r.typ == "CASH_OUT"), boolFloat(r.typ == "DEBIT"), boolFloat(r.typ == "PAYMENT"), boolFloat(r.typ == "TRANSFER"), boolFloat(r.destType == "M"),
	}
}

func boolFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func loadDataset(path string) (*dataset, error) {
	_, idx, err := readHeader(path)
	if err != nil {
		return nil, err
	}

	// Pass 1 computes the exact linear-interpolated 99.99th percentile used by NumPy.
	r, f, err := openCSV(path)
	if err != nil {
		return nil, err
	}
	amounts := make([]float64, 0, 6400000)
	var lastStep float64
	rows := 0
	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			f.Close()
			return nil, fmt.Errorf("CSV row %d: %w", rows+2, e)
		}
		x, e := parseRow(rec, idx)
		if e != nil {
			f.Close()
			return nil, fmt.Errorf("CSV row %d: %w", rows+2, e)
		}
		if rows > 0 && x.step < lastStep {
			f.Close()
			return nil, fmt.Errorf("input is not chronological at data row %d", rows+1)
		}
		lastStep = x.step
		amounts = append(amounts, x.amount)
		rows++
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	if rows < 2 {
		return nil, errors.New("dataset must contain at least two data rows")
	}
	sort.Float64s(amounts)
	pos := float64(rows-1) * 0.9999
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))
	cap := amounts[lo] + (amounts[hi]-amounts[lo])*(pos-float64(lo))
	amounts = nil
	runtime.GC()

	split := int(float64(rows) * trainRatio)
	d := &dataset{cap: cap, trainX: make([]float32, split*featureCount), trainY: make([]uint8, split), testX: make([]float32, (rows-split)*featureCount), testY: make([]uint8, rows-split)}
	var count [featureCount]int
	var m2 [featureCount]float64
	classCounts := [2]int{}
	r, f, err = openCSV(path)
	if err != nil {
		return nil, err
	}
	for i := 0; i < rows; i++ {
		rec, e := r.Read()
		if e != nil {
			f.Close()
			return nil, fmt.Errorf("read row %d: %w", i+2, e)
		}
		x, e := parseRow(rec, idx)
		if e != nil {
			f.Close()
			return nil, e
		}
		if int(x.y) > 1 {
			f.Close()
			return nil, errors.New("invalid class label")
		}
		classCounts[x.y]++
		if i >= split {
			continue
		}
		features := x.features(cap)
		for j, v := range features {
			count[j]++
			delta := v - d.mean[j]
			d.mean[j] += delta / float64(count[j])
			m2[j] += delta * (v - d.mean[j])
		}
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	for j := range d.std {
		d.std[j] = math.Sqrt(m2[j] / float64(count[j]))
		if d.std[j] == 0 {
			d.std[j] = 1
		}
	}

	// Pass 3 stores float32 features to keep the full train/test matrix within practical RAM.
	r, f, err = openCSV(path)
	if err != nil {
		return nil, err
	}
	for i := 0; i < rows; i++ {
		rec, e := r.Read()
		if e != nil {
			f.Close()
			return nil, fmt.Errorf("read row %d: %w", i+2, e)
		}
		x, e := parseRow(rec, idx)
		if e != nil {
			f.Close()
			return nil, e
		}
		features := x.features(cap)
		if i < split {
			d.trainY[i] = x.y
			base := i * featureCount
			for j, v := range features {
				d.trainX[base+j] = float32((v - d.mean[j]) / d.std[j])
			}
		} else {
			k := i - split
			d.testY[k] = x.y
			base := k * featureCount
			for j, v := range features {
				d.testX[base+j] = float32((v - d.mean[j]) / d.std[j])
			}
		}
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	if classCounts[0] == 0 || classCounts[1] == 0 {
		return nil, errors.New("both classes must occur in the dataset")
	}
	fmt.Printf("Chronological split: %d/%d; scaler fitted on training rows only\n", split, rows-split)
	return d, nil
}

type gradient struct {
	values [featureCount]float64
	bias   float64
}
type model struct {
	weights [featureCount]float64
	bias    float64
}

func sigmoid(z float64) float64 {
	if z >= 0 {
		e := math.Exp(-z)
		return 1 / (1 + e)
	}
	e := math.Exp(z)
	return e / (1 + e)
}

func accumulate(x []float32, y []uint8, start, end int, w model) gradient {
	var g gradient
	for i := start; i < end; i++ {
		base := i * featureCount
		z := w.bias
		for j := 0; j < featureCount; j++ {
			z += w.weights[j] * float64(x[base+j])
		}
		p := sigmoid(z)
		weight := 0.5006
		if y[i] == 1 {
			weight = 387.35
		}
		err := (p - float64(y[i])) * weight
		g.bias += err
		for j := 0; j < featureCount; j++ {
			g.values[j] += err * float64(x[base+j])
		}
	}
	return g
}

func sequentialGradient(x []float32, y []uint8, w model) gradient {
	return accumulate(x, y, 0, len(y), w)
}

func parallelGradient(x []float32, y []uint8, w model, workers int) gradient {
	if workers <= 1 {
		return sequentialGradient(x, y, w)
	}
	if workers > len(y) {
		workers = len(y)
	}
	type job struct{ start, end int }
	jobs := make(chan job, workers)
	results := make(chan gradient, workers)
	var wg sync.WaitGroup
	for n := 0; n < workers; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				results <- accumulate(x, y, task.start, task.end, w)
			}
		}()
	}
	for n := 0; n < workers; n++ {
		start := len(y) * n / workers
		end := len(y) * (n + 1) / workers
		jobs <- job{start, end}
	}
	close(jobs)
	go func() { wg.Wait(); close(results) }()
	var total gradient
	for partial := range results {
		total.bias += partial.bias
		for j := range total.values {
			total.values[j] += partial.values[j]
		}
	}
	return total
}

func update(w *model, g gradient, n int, rate float64) {
	for j := range w.weights {
		w.weights[j] -= rate * (g.values[j]/float64(n) + 0.0001*w.weights[j])
	}
	w.bias -= rate * g.bias / float64(n)
}

func train(x []float32, y []uint8, workers, epochs int) model {
	var w model
	for e := 0; e < epochs; e++ {
		var g gradient
		if workers <= 1 {
			g = sequentialGradient(x, y, w)
		} else {
			g = parallelGradient(x, y, w, workers)
		}
		update(&w, g, len(y), 0.1)
	}
	return w
}

func countOnes(y []uint8) int {
	n := 0
	for _, v := range y {
		if v == 1 {
			n++
		}
	}
	return n
}

func trimmedMean(v []float64, proportion float64) float64 {
	s := append([]float64(nil), v...)
	sort.Float64s(s)
	cut := int(float64(len(s)) * proportion)
	if 2*cut >= len(s) {
		cut = 0
	}
	total := 0.0
	for _, x := range s[cut : len(s)-cut] {
		total += x
	}
	return total / float64(len(s)-2*cut)
}

func runBench(d *dataset, workers []int, repeats, epochs int, out string) error {
	type result struct {
		mode    string
		workers int
		run     int
		ms      float64
		cpuSec  float64
		heapMB  float64
	}
	var results []result
	configs := []struct {
		mode    string
		workers int
	}{{"sequential", 1}}
	for _, n := range workers {
		if n > 1 {
			configs = append(configs, struct {
				mode    string
				workers int
			}{"worker_pool", n})
		}
	}
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	fmt.Printf("Benchmark: %d repetitions x %d epoch(s), trimmed mean=10%%\n", repeats, epochs)
	for _, c := range configs {
		times := make([]float64, 0, repeats)
		for run := 1; run <= repeats; run++ {
			cpu0 := processCPUSeconds()
			start := time.Now()
			_ = train(d.trainX, d.trainY, c.workers, epochs)
			elapsed := time.Since(start).Seconds() * 1000
			cpu1 := processCPUSeconds()
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			cpuDelta := cpu1 - cpu0
			if cpuDelta < 0 {
				cpuDelta = 0
			}
			results = append(results, result{c.mode, c.workers, run, elapsed, cpuDelta, float64(mem.HeapAlloc) / (1024 * 1024)})
			times = append(times, elapsed)
			fmt.Printf("%s workers=%d run=%d/%d %.3f ms\n", c.mode, c.workers, run, repeats, elapsed)
		}
		fmt.Printf("%s workers=%d trimmed_mean_10pct=%.3f ms\n", c.mode, c.workers, trimmedMean(times, 0.1))
	}
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.Write([]string{"mode", "workers", "run", "elapsed_ms", "process_cpu_seconds", "heap_alloc_mb"}); err != nil {
		return err
	}
	for _, r := range results {
		if err := w.Write([]string{r.mode, strconv.Itoa(r.workers), strconv.Itoa(r.run), fmt.Sprintf("%.6f", r.ms), fmt.Sprintf("%.6f", r.cpuSec), fmt.Sprintf("%.3f", r.heapMB)}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	fmt.Printf("Wrote raw observations to %s; load heap %.1f MiB; benchmark heap %.1f MiB; logical CPUs=%d\n", out, float64(before.HeapAlloc)/(1024*1024), float64(after.HeapAlloc)/(1024*1024), runtime.NumCPU())
	return nil
}

func evaluate(d *dataset, epochs int) error {
	for _, n := range []int{1, runtime.NumCPU()} {
		start := time.Now()
		w := train(d.trainX, d.trainY, n, epochs)
		elapsed := time.Since(start)
		tp, fp, fn := 0, 0, 0
		for i, y := range d.testY {
			base := i * featureCount
			z := w.bias
			for j := 0; j < featureCount; j++ {
				z += w.weights[j] * float64(d.testX[base+j])
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
		precision := float64(tp) / float64(max(1, tp+fp))
		recall := float64(tp) / float64(max(1, tp+fn))
		f1 := 2 * precision * recall / max(1e-12, precision+recall)
		fmt.Printf("workers=%d epochs=%d elapsed=%.3fs test confusion TP=%d FP=%d FN=%d precision=%.5f recall=%.5f F1=%.5f\n", n, epochs, elapsed.Seconds(), tp, fp, fn, precision, recall, f1)
	}
	return nil
}

func fatal(err error) { fmt.Fprintln(os.Stderr, "error:", err); os.Exit(1) }
