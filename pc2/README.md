# PC2 - PaySim weighted logistic regression

Integrated team report: [`CC65-PC2-202620-Equipo.pdf`](CC65-PC2-202620-Equipo.pdf).

## Reproduce

From this directory, with Go installed:

```powershell
go run . -input ..\paysim.csv -mode bench -repeats 30 -bench-epochs 1 -workers 1,2,4,8 -out results.csv
go run . -input ..\paysim.csv -mode evaluate -epochs 5
```

The program reads the full original CSV, checks chronological order, computes the
99.99th percentile amount cap, fits population standardization statistics on the
first 80% of records, and applies the same transformation to the held-out final
20%. The 15 features and class weights follow the PC1 preprocessing description.

Benchmark timings cover gradient computation, worker coordination, reduction,
and model update, but exclude CSV loading and standardization. Raw observations
are written to `results.csv`; the console prints the 10% trimmed mean for each
configuration. Set worker counts to values appropriate for the machine. The
program reports logical CPU count, Go heap allocation, and process CPU time per run.

## Promela

`sync.pml` abstracts two workers and the coordinator. Each worker publishes one
private partial result; only the coordinator consumes results and updates shared
parameters. With Spin installed, verify assertions and deadlock absence using:

```powershell
spin -a sync.pml
gcc -O2 -o pan.exe pan.c
.\pan.exe -a
```
