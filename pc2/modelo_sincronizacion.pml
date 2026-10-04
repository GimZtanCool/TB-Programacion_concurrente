/*
   Finite abstraction of PC2's gradient reduction.
   Each worker owns its partial gradient; the coordinator alone applies the update.
*/
#define WORKERS 2

chan gradients = [WORKERS] of { byte };
bool partial_ready[WORKERS];
bool consumed[WORKERS];
byte updates = 0;

proctype Worker(byte id)
{
    /* Private accumulation is abstracted as one completed partial result. */
    atomic {
        assert(!partial_ready[id]);
        partial_ready[id] = true;
    }
    gradients!id;
}

proctype Coordinator()
{
    byte id;
    byte received = 0;

    do
    :: received < WORKERS ->
        gradients?id;
        atomic {
            assert(partial_ready[id]);
            assert(!consumed[id]);
            consumed[id] = true;
            received++;
            if
            :: received == WORKERS ->
                /* Sole writer of shared model parameters. */
                updates++;
                assert(updates == 1);
            :: else -> skip
            fi
        }
    :: else -> break
    od;

    assert(consumed[0] && consumed[1]);
    assert(updates == 1);
}

init {
    atomic {
        run Worker(0);
        run Worker(1);
        run Coordinator();
    }
}
