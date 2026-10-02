/* Finite abstraction of parallelGradient and the subsequent update.
 * Job ids identify disjoint partitions, NOT particular workers.
 * No atomic block covers the protocol or the writer critical section.
 */
#ifndef WORKERS
#define WORKERS 2
#endif
#define JOBS WORKERS

chan jobs = [WORKERS] of { byte };
chan results = [WORKERS] of { byte };
bool jobs_closed = false;
bool results_closed = false;
bool ready[JOBS];
bool consumed[JOBS];
byte workers_done = 0;
byte received = 0;
byte writers = 0;
byte updates = 0;
bool finished = false;

proctype Producer() {
    byte id = 0;
    do
    :: id < JOBS -> jobs!id; id++
    :: else -> break
    od;
    jobs_closed = true
}

proctype Worker(byte worker_id) {
    byte id;
    do
    :: jobs?id ->
        assert(!ready[id]);
        ready[id] = true; /* partial gradient is private */
#ifdef BAD_WRITER
        /* Negative control: workers incorrectly mutate shared parameters. */
        writers++;
        assert(writers == 1);
        writers--;
#endif
        results!id
    :: (jobs_closed && len(jobs) == 0) -> break
    od;
#ifdef BAD_WAIT
    /* Negative control: simulate a missing WaitGroup.Done. */
    if
    :: worker_id == 0 -> skip
    :: else -> workers_done++
    fi
#else
    workers_done++
#endif
}

proctype Closer() {
    workers_done == WORKERS;
    results_closed = true
}

proctype Coordinator() {
    byte id;
    jobs_closed; /* Go enqueues and closes jobs before receiving results. */
    do
    :: results?id ->
        assert(id < JOBS);
        assert(ready[id]);
        assert(!consumed[id]);
        consumed[id] = true;
        received++
    :: (results_closed && len(results) == 0) -> break
    od;
    assert(workers_done == WORKERS);
    assert(received == JOBS);
    id = 0;
    do
    :: id < JOBS -> assert(consumed[id]); id++
    :: else -> break
    od;
    assert(writers == 0);
    writers++;
    assert(writers == 1); /* exclusive access through sole ownership */
    updates++;
    assert(updates == 1);
    writers--;
    finished = true
}

/* Liveness is checked separately from invalid end states and assertions. */
ltl completion { <> finished }

init {
    byte id = 0;
    atomic {
        run Producer();
        do
        :: id < WORKERS -> run Worker(id); id++
        :: else -> break
        od;
        run Closer();
        run Coordinator()
    }
}
