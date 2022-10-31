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
	bufferWriter.WriteString(
		fmt.Sprintf("Alloc (Alloc): %d bytes [%d mb]\n", memStats.Alloc, bytesToMegabytes(memStats.Alloc)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Cumulative Alloc (TotalAlloc): %d bytes [%d mb]\n", memStats.TotalAlloc, bytesToMegabytes(memStats.TotalAlloc)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Obtained from OS (Sys): %d bytes [%d mb]\n", memStats.Sys, bytesToMegabytes(memStats.Sys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Pointer Lookups (Lookups): %d times\n", memStats.Lookups),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Cumulative Objects Allocated (Mallocs): %d times\n", memStats.Mallocs),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Cumulative Objects Freed (Frees): %d times\n", memStats.Frees),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("\tLive Objects (Mallocs - Frees): %d \n", (memStats.Mallocs - memStats.Frees)),
	)

	bufferWriter.WriteString("\n\n*** Heap Memory statistics ***\n\n")
	bufferWriter.WriteString(
		fmt.Sprintf("Heap Allocation (HeapAlloc): %d bytes [%d mb]\n", memStats.HeapAlloc, bytesToMegabytes(memStats.HeapAlloc)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Largest Heap Memory Size from OS (HeapSys): %d bytes [%d mb]\n", memStats.HeapSys, bytesToMegabytes(memStats.HeapSys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Idle Heap Spans (HeapIdle): %d bytes [%d mb]\n", memStats.HeapIdle, bytesToMegabytes(memStats.HeapIdle)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("In-use Heap Spans (HeapInuse): %d bytes [%d mb]\n", memStats.HeapInuse, bytesToMegabytes(memStats.HeapInuse)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Heap Released (HeapReleased): %d bytes [%d mb]\n", memStats.HeapReleased, bytesToMegabytes(memStats.HeapReleased)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Heap Objects (HeapObjects): %d \n", memStats.HeapObjects),
	)

	bufferWriter.WriteString("\n\n*** Stack Memory statistics ***\n\n")
	bufferWriter.WriteString(
		fmt.Sprintf("Stack in Use (StackInuse): %d bytes [%d mb]\n", memStats.StackInuse, bytesToMegabytes(memStats.StackInuse)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Stack from OS (StackSys): %d bytes [%d mb]\n", memStats.StackSys, bytesToMegabytes(memStats.StackSys)),
	)

	bufferWriter.WriteString("\n\n*** Off-heap Memory statistics ***\n\n")
	bufferWriter.WriteString(
		fmt.Sprintf("Mspan structures memory (MSpanInuse): %d bytes [%d mb]\n", memStats.MSpanInuse, bytesToMegabytes(memStats.MSpanInuse)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Mspan structures memory from OS (MSpanSys): %d bytes [%d mb]\n", memStats.MSpanSys, bytesToMegabytes(memStats.MSpanSys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("MCache structures memory (MCacheInuse): %d bytes [%d mb]\n", memStats.MCacheInuse, bytesToMegabytes(memStats.MCacheInuse)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("MCache structures memory from OS (MCacheSys): %d bytes [%d mb]\n", memStats.MCacheSys, bytesToMegabytes(memStats.MCacheSys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Profiling bucket hash tables size (BuckHashSys): %d bytes [%d mb]\n", memStats.BuckHashSys, bytesToMegabytes(memStats.BuckHashSys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("GC metadata size (GCSys): %d bytes [%d mb]\n", memStats.GCSys, bytesToMegabytes(memStats.GCSys)),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Miscellaneous (OtherSys): %d bytes [%d mb]\n", memStats.OtherSys, bytesToMegabytes(memStats.OtherSys)),
	)

	bufferWriter.WriteString("\n\n*** Garbage Collector statistics ***\n\n")
	bufferWriter.WriteString(
		fmt.Sprintf("Next GC Target (NextGC): %d \n", memStats.NextGC),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Last GC in UNIX epoch (LastGC): %d \n", memStats.LastGC),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Cumulative ns in GC stop-the-world (PauseTotalNs): %d \n", memStats.PauseTotalNs),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Completed Cycles (NumGC): %d \n", memStats.NumGC),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("Forced Cycles (NumForcedGC): %d \n", memStats.NumForcedGC),
	)
	bufferWriter.WriteString(
		fmt.Sprintf("CPU Fraction (GCCPUFraction): %f \n", memStats.GCCPUFraction),
	)

	bufferWriter.Flush()
	w.Write(b.Bytes())
}
