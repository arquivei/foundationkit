package app

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
)

// nolint: errcheck
func dumpGoroutines(w http.ResponseWriter, r *http.Request) {
	var b [4 * 1024 * 1024]byte

	n := runtime.Stack(b[:], true)

	if n > 0 {
		w.Write(b[0:n])
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to write Stack into the buffer"))
	}
}

// nolint: errcheck
func dumpMemProfile(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	bufferWriter := bufio.NewWriter(&b)

	err := pprof.WriteHeapProfile(bufferWriter)
	bufferWriter.Flush()

	switch {
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to write heap into the buffer: " + err.Error()))
	case b.Len() == 0:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Heap buffer is empty"))
	default:
		w.Write(b.Bytes())
	}
}

func bytesToMegabytes(bytes uint64) uint64 {
	return bytes / (1024 * 1024)
}

// nolint: errcheck
func dumpMemStats(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	bw := bufio.NewWriter(w)

	fmt.Fprint(bw, "*** General statistics ***\n\n")
	fmt.Fprintf(bw, "Alloc (Alloc): %d bytes [%d mb]\n", memStats.Alloc, bytesToMegabytes(memStats.Alloc))
	fmt.Fprintf(bw, "Cumulative Alloc (TotalAlloc): %d bytes [%d mb]\n", memStats.TotalAlloc, bytesToMegabytes(memStats.TotalAlloc))
	fmt.Fprintf(bw, "Obtained from OS (Sys): %d bytes [%d mb]\n", memStats.Sys, bytesToMegabytes(memStats.Sys))
	fmt.Fprintf(bw, "Pointer Lookups (Lookups): %d times\n", memStats.Lookups)
	fmt.Fprintf(bw, "Cumulative Objects Allocated (Mallocs): %d times\n", memStats.Mallocs)
	fmt.Fprintf(bw, "Cumulative Objects Freed (Frees): %d times\n", memStats.Frees)
	fmt.Fprintf(bw, "\tLive Objects (Mallocs - Frees): %d \n", (memStats.Mallocs - memStats.Frees))

	fmt.Fprint(bw, "\n\n*** Heap Memory statistics ***\n\n")
	fmt.Fprintf(bw, "Heap Allocation (HeapAlloc): %d bytes [%d mb]\n", memStats.HeapAlloc, bytesToMegabytes(memStats.HeapAlloc))
	fmt.Fprintf(bw, "Largest Heap Memory Size from OS (HeapSys): %d bytes [%d mb]\n", memStats.HeapSys, bytesToMegabytes(memStats.HeapSys))
	fmt.Fprintf(bw, "Idle Heap Spans (HeapIdle): %d bytes [%d mb]\n", memStats.HeapIdle, bytesToMegabytes(memStats.HeapIdle))
	fmt.Fprintf(bw, "In-use Heap Spans (HeapInuse): %d bytes [%d mb]\n", memStats.HeapInuse, bytesToMegabytes(memStats.HeapInuse))
	fmt.Fprintf(bw, "Heap Released (HeapReleased): %d bytes [%d mb]\n", memStats.HeapReleased, bytesToMegabytes(memStats.HeapReleased))
	fmt.Fprintf(bw, "Heap Objects (HeapObjects): %d \n", memStats.HeapObjects)

	fmt.Fprint(bw, "\n\n*** Stack Memory statistics ***\n\n")
	fmt.Fprintf(bw, "Stack in Use (StackInuse): %d bytes [%d mb]\n", memStats.StackInuse, bytesToMegabytes(memStats.StackInuse))
	fmt.Fprintf(bw, "Stack from OS (StackSys): %d bytes [%d mb]\n", memStats.StackSys, bytesToMegabytes(memStats.StackSys))

	fmt.Fprint(bw, "\n\n*** Off-heap Memory statistics ***\n\n")
	fmt.Fprintf(bw, "Mspan structures memory (MSpanInuse): %d bytes [%d mb]\n", memStats.MSpanInuse, bytesToMegabytes(memStats.MSpanInuse))
	fmt.Fprintf(bw, "Mspan structures memory from OS (MSpanSys): %d bytes [%d mb]\n", memStats.MSpanSys, bytesToMegabytes(memStats.MSpanSys))
	fmt.Fprintf(bw, "MCache structures memory (MCacheInuse): %d bytes [%d mb]\n", memStats.MCacheInuse, bytesToMegabytes(memStats.MCacheInuse))
	fmt.Fprintf(bw, "MCache structures memory from OS (MCacheSys): %d bytes [%d mb]\n", memStats.MCacheSys, bytesToMegabytes(memStats.MCacheSys))
	fmt.Fprintf(bw, "Profiling bucket hash tables size (BuckHashSys): %d bytes [%d mb]\n", memStats.BuckHashSys, bytesToMegabytes(memStats.BuckHashSys))
	fmt.Fprintf(bw, "GC metadata size (GCSys): %d bytes [%d mb]\n", memStats.GCSys, bytesToMegabytes(memStats.GCSys))
	fmt.Fprintf(bw, "Miscellaneous (OtherSys): %d bytes [%d mb]\n", memStats.OtherSys, bytesToMegabytes(memStats.OtherSys))

	fmt.Fprint(bw, "\n\n*** Garbage Collector statistics ***\n\n")
	fmt.Fprintf(bw, "Next GC Target (NextGC): %d \n", memStats.NextGC)
	fmt.Fprintf(bw, "Last GC in UNIX epoch (LastGC): %d \n", memStats.LastGC)
	fmt.Fprintf(bw, "Cumulative ns in GC stop-the-world (PauseTotalNs): %d \n", memStats.PauseTotalNs)
	fmt.Fprintf(bw, "Completed Cycles (NumGC): %d \n", memStats.NumGC)
	fmt.Fprintf(bw, "Forced Cycles (NumForcedGC): %d \n", memStats.NumForcedGC)
	fmt.Fprintf(bw, "CPU Fraction (GCCPUFraction): %f \n", memStats.GCCPUFraction)

	bw.Flush()
}
