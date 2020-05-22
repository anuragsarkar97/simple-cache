
# Simple Cache

## Introduction 

The **simple cache** is literally the simplest cache possible
built for go. It is lightweight and requires no external dependency. 
It is super fast with no extra feature what so ever. 

It does **SET** and **GET** very well. You caan store data to a persistent file if 
required but that's about it.

## Usage

```go
c := CreateNewCache("cache_data", 1000, false) // required
c.setTTL(5) // required
c.Set("some key", "some value", 10)
c.Get("some key")
c.close() // required
````

## Design

It uses a Heap, and a HashMap to maintain order for expiry and that's about it.

## Bench marking

```text
Showing nodes accounting for 4.13s, 96.27% of 4.29s total
Dropped 40 nodes (cum <= 0.02s)
Showing top 50 nodes out of 55
      flat  flat%   sum%        cum   cum%
     1.60s 37.30% 37.30%      1.60s 37.30%  runtime.pthread_cond_signal
     1.26s 29.37% 66.67%      1.27s 29.60%  runtime.usleep
     0.70s 16.32% 82.98%      0.70s 16.32%  runtime.pthread_cond_wait
     0.20s  4.66% 87.65%      0.21s  4.90%  runtime.nanotime1
     0.10s  2.33% 89.98%      0.10s  2.33%  runtime.madvise
     0.05s  1.17% 91.14%      0.05s  1.17%  runtime.procyield
     0.04s  0.93% 92.07%      0.06s  1.40%  runtime.walltime1
     0.03s   0.7% 92.77%      0.05s  1.17%  exercise.RandStringRunes
     0.03s   0.7% 93.47%      0.03s   0.7%  runtime.pthread_kill
     0.03s   0.7% 94.17%      0.09s  2.10%  runtime.scanobject
     0.02s  0.47% 94.64%      0.24s  5.59%  exercise.(*SimpleCache).processExpiry
     0.02s  0.47% 95.10%      0.04s  0.93%  runtime.libcCall
     0.02s  0.47% 95.57%      0.03s   0.7%  sync.(*Mutex).Unlock (inline)
     0.02s  0.47% 96.04%      0.12s  2.80%  time.Now
     0.01s  0.23% 96.27%      0.15s  3.50%  sync.(*Mutex).lockSlow
         0     0% 96.27%      0.11s  2.56%  exercise.(*SimpleCache).Set
         0     0% 96.27%      0.24s  5.59%  exercise.(*SimpleCache).concurrentProcessChecks
         0     0% 96.27%      0.16s  3.73%  exercise.BenchmarkSimpleCache
         0     0% 96.27%      0.16s  3.73%  exercise.benchmark
         0     0% 96.27%      0.11s  2.56%  runtime.(*mheap).alloc.func1
         0     0% 96.27%      0.11s  2.56%  runtime.(*mheap).allocSpan
         0     0% 96.27%      0.07s  1.63%  runtime.checkTimers
         0     0% 96.27%      2.05s 47.79%  runtime.findrunnable
         0     0% 96.27%      0.10s  2.33%  runtime.gcBgMarkWorker
         0     0% 96.27%      0.11s  2.56%  runtime.gcBgMarkWorker.func2
         0     0% 96.27%      0.11s  2.56%  runtime.gcDrain
         0     0% 96.27%      1.46s 34.03%  runtime.goready.func1
         0     0% 96.27%      2.20s 51.28%  runtime.mcall
         0     0% 96.27%      1.59s 37.06%  runtime.mstart
         0     0% 96.27%      0.21s  4.90%  runtime.nanotime (inline)
         0     0% 96.27%      0.70s 16.32%  runtime.notesleep
         0     0% 96.27%      1.61s 37.53%  runtime.notewakeup
         0     0% 96.27%      2.20s 51.28%  runtime.park_m
         0     0% 96.27%      0.03s   0.7%  runtime.preemptM (inline)
         0     0% 96.27%      0.03s   0.7%  runtime.preemptone
         0     0% 96.27%      1.46s 34.03%  runtime.ready
         0     0% 96.27%      0.15s  3.50%  runtime.resetspinning
         0     0% 96.27%      1.27s 29.60%  runtime.runqgrab
         0     0% 96.27%      1.27s 29.60%  runtime.runqsteal
         0     0% 96.27%      2.20s 51.28%  runtime.schedule
         0     0% 96.27%      0.70s 16.32%  runtime.semasleep
         0     0% 96.27%      1.61s 37.53%  runtime.semawakeup
         0     0% 96.27%      0.03s   0.7%  runtime.signalM (inline)
         0     0% 96.27%      1.61s 37.53%  runtime.startm
         0     0% 96.27%      0.70s 16.32%  runtime.stopm
         0     0% 96.27%      0.10s  2.33%  runtime.sysUsed (inline)
         0     0% 96.27%      1.70s 39.63%  runtime.systemstack
         0     0% 96.27%      1.61s 37.53%  runtime.wakep (inline)
         0     0% 96.27%      0.06s  1.40%  runtime.walltime (inline)
         0     0% 96.27%      0.15s  3.50%  sync.(*Mutex).Lock (inline)

```