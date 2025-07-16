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

	var b bytes.Buffer
	bufferWriter := bufio.NewWriter(&b)

	bufferWriter.WriteString("*** General statistics ***\n\n")
	fmt.Fprintf(bufferWriter,
		"Alloc (Alloc): %d bytes [%d mb]\n", memStats.Alloc, bytesToMegabytes(memStats.Alloc))
	fmt.Fprintf(bufferWriter,
		"Cumulative Alloc (TotalAlloc): %d bytes [%d mb]\n", memStats.TotalAlloc, bytesToMegabytes(memStats.TotalAlloc))
	fmt.Fprintf(bufferWriter,
		"Obtained from OS (Sys): %d bytes [%d mb]\n", memStats.Sys, bytesToMegabytes(memStats.Sys))
	fmt.Fprintf(bufferWriter,
		"Pointer Lookups (Lookups): %d times\n", memStats.Lookups)
	fmt.Fprintf(bufferWriter,
		"Cumulative Objects Allocated (Mallocs): %d times\n", memStats.Mallocs)
	fmt.Fprintf(bufferWriter,
		"Cumulative Objects Freed (Frees): %d times\n", memStats.Frees)
	fmt.Fprintf(bufferWriter,
		"\tLive Objects (Mallocs - Frees): %d \n", (memStats.Mallocs - memStats.Frees))

	bufferWriter.WriteString("\n\n*** Heap Memory statistics ***\n\n")
	fmt.Fprintf(bufferWriter,
		"Heap Allocation (HeapAlloc): %d bytes [%d mb]\n", memStats.HeapAlloc, bytesToMegabytes(memStats.HeapAlloc))
	fmt.Fprintf(bufferWriter,
		"Largest Heap Memory Size from OS (HeapSys): %d bytes [%d mb]\n", memStats.HeapSys, bytesToMegabytes(memStats.HeapSys))
	fmt.Fprintf(bufferWriter,
		"Idle Heap Spans (HeapIdle): %d bytes [%d mb]\n", memStats.HeapIdle, bytesToMegabytes(memStats.HeapIdle))
	fmt.Fprintf(bufferWriter,
		"In-use Heap Spans (HeapInuse): %d bytes [%d mb]\n", memStats.HeapInuse, bytesToMegabytes(memStats.HeapInuse))
	fmt.Fprintf(bufferWriter,
		"Heap Released (HeapReleased): %d bytes [%d mb]\n", memStats.HeapReleased, bytesToMegabytes(memStats.HeapReleased))
	fmt.Fprintf(bufferWriter,
		"Heap Objects (HeapObjects): %d \n", memStats.HeapObjects)

	bufferWriter.WriteString("\n\n*** Stack Memory statistics ***\n\n")
	fmt.Fprintf(bufferWriter,
		"Stack in Use (StackInuse): %d bytes [%d mb]\n", memStats.StackInuse, bytesToMegabytes(memStats.StackInuse))
	fmt.Fprintf(bufferWriter,
		"Stack from OS (StackSys): %d bytes [%d mb]\n", memStats.StackSys, bytesToMegabytes(memStats.StackSys))

	bufferWriter.WriteString("\n\n*** Off-heap Memory statistics ***\n\n")
	fmt.Fprintf(bufferWriter,
		"Mspan structures memory (MSpanInuse): %d bytes [%d mb]\n", memStats.MSpanInuse, bytesToMegabytes(memStats.MSpanInuse))
	fmt.Fprintf(bufferWriter,
		"Mspan structures memory from OS (MSpanSys): %d bytes [%d mb]\n", memStats.MSpanSys, bytesToMegabytes(memStats.MSpanSys))
	fmt.Fprintf(bufferWriter,
		"MCache structures memory (MCacheInuse): %d bytes [%d mb]\n", memStats.MCacheInuse, bytesToMegabytes(memStats.MCacheInuse))
	fmt.Fprintf(bufferWriter,
		"MCache structures memory from OS (MCacheSys): %d bytes [%d mb]\n", memStats.MCacheSys, bytesToMegabytes(memStats.MCacheSys))
	fmt.Fprintf(bufferWriter,
		"Profiling bucket hash tables size (BuckHashSys): %d bytes [%d mb]\n", memStats.BuckHashSys, bytesToMegabytes(memStats.BuckHashSys))
	fmt.Fprintf(bufferWriter,
		"GC metadata size (GCSys): %d bytes [%d mb]\n", memStats.GCSys, bytesToMegabytes(memStats.GCSys))
	fmt.Fprintf(bufferWriter,
		"Miscellaneous (OtherSys): %d bytes [%d mb]\n", memStats.OtherSys, bytesToMegabytes(memStats.OtherSys))

	bufferWriter.WriteString("\n\n*** Garbage Collector statistics ***\n\n")
	fmt.Fprintf(bufferWriter,
		"Next GC Target (NextGC): %d \n", memStats.NextGC)
	fmt.Fprintf(bufferWriter,
		"Last GC in UNIX epoch (LastGC): %d \n", memStats.LastGC)
	fmt.Fprintf(bufferWriter,
		"Cumulative ns in GC stop-the-world (PauseTotalNs): %d \n", memStats.PauseTotalNs)
	fmt.Fprintf(bufferWriter,
		"Completed Cycles (NumGC): %d \n", memStats.NumGC)
	fmt.Fprintf(bufferWriter,
		"Forced Cycles (NumForcedGC): %d \n", memStats.NumForcedGC)
	fmt.Fprintf(bufferWriter,
		"CPU Fraction (GCCPUFraction): %f \n", memStats.GCCPUFraction)

	bufferWriter.Flush()
	w.Write(b.Bytes())
}
