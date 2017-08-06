package system

import "testing"

func TestProcessorsRemove(t *testing.T) {
	processors := Processors{1, 2, 3, 4, 5, 6, 7, 8}

	rest := processors.Remove(Processors{1, 5, 7}...)
	for i, v := range []int64{2, 3, 4, 6, 8} {
		if v != rest[i] {
			t.Fatalf("rest[%v] != %v, got %v", i, v, rest[i])
		}
	}

	var removed Processors
	rest = processors.Remove(removed...)
	for i, v := range []int64{1, 2, 3, 4, 5, 6, 7, 8} {
		if v != rest[i] {
			t.Fatalf("rest[%v] != %v, got %v", i, v, rest[i])
		}
	}
}

func loadCpuInfo(nsocket, ncore, ht int) *CpuInfo {
	cpuinfo := CpuInfo{Processors: int64(nsocket * ncore * ht), Sockets: Sockets{}}
	for i := 0; i < nsocket; i++ {
		socket := Socket{Index: int64(i), Cores: Cores{}}
		for j := 0; j < ncore; j++ {
			core := Core{Index: int64(j), Processors: Processors{int64(i*ncore + j), int64(i*ncore + j + nsocket*ncore)}}
			socket.Cores = append(socket.Cores, &core)
		}
		cpuinfo.Sockets = append(cpuinfo.Sockets, &socket)
	}
	return &cpuinfo
}

func TestIsProcessorInSameSocket(t *testing.T) {
	cpuinfo := loadCpuInfo(2, 6, 2)

	if !cpuinfo.IsProcessorInSameSocket(0, 15) {
		t.Fatalf("0 and 15 should be in socket 0")
	}

	if cpuinfo.IsProcessorInSameSocket(0, 6) {
		t.Fatalf("0 in socket 0, 6 in socket 1")
	}

	cpuinfo = loadCpuInfo(2, 2, 2)

	if !cpuinfo.IsProcessorInSameSocket(0, 1) {
		t.Fatalf("0 in socket 0, 1 in socket 0")
	}
}

func TestIsProcessorInSameCore(t *testing.T) {
	cpuinfo := loadCpuInfo(2, 6, 2)

	if cpuinfo.IsProcessorInSameCore(0, 1) {
		t.Fatalf("0 in core 0 socket 0, 1 in core 1 socket 0")
	}

	if cpuinfo.IsProcessorInSameCore(0, 6) {
		t.Fatalf("0 in core 0 socket 0, 6 in core 0 socket 1")
	}

	if cpuinfo.IsProcessorInSameCore(0, 8) {
		t.Fatalf("0 in core 0 socket 0, 8 in core 2 socket 1")
	}

	if !cpuinfo.IsProcessorInSameCore(0, 12) {
		t.Fatalf("0 and 12 in core 0 socket 0")
	}

	cpuinfo = loadCpuInfo(2, 2, 2)

	if !cpuinfo.IsProcessorInSameCore(0, 4) {
		t.Fatalf("0 and 4 in core 0 socket 0")
	}
}
