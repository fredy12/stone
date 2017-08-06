package system

// CpuInfo contains CPU topology of the host system.
type CpuInfo struct {
	// Total usable CPUs (logic processor)
	Processors int64

	// The unique "Physical ID"
	Sockets Sockets
}

func (c *CpuInfo) AddSocket(idx int64) {
	if c.Sockets != nil && c.Sockets.Get(idx) != nil {
		return
	}
	c.Sockets = append(c.Sockets, &Socket{Index: idx, Cores: Cores{}})
}

func (c *CpuInfo) IsProcessorInSameSocket(x, y int64) bool {
	var (
		xSock int64 = -1
		ySock int64 = -1
	)

	for _, socket := range c.Sockets {
		for _, core := range socket.Cores {
			for _, processor := range core.Processors {
				switch processor {
				case x:
					xSock = socket.Index
				case y:
					ySock = socket.Index
				}
				if xSock != -1 && ySock != -1 {
					break
				}
			}
		}
	}

	return xSock != -1 && xSock == ySock
}

func (c *CpuInfo) IsProcessorInSameCore(x, y int64) bool {
	var (
		xCore int64 = -1
		xSock int64 = -1
		yCore int64 = -1
		ySock int64 = -1
	)

	for _, socket := range c.Sockets {
		for _, core := range socket.Cores {
			for _, processor := range core.Processors {
				switch processor {
				case x:
					xCore = core.Index
					xSock = socket.Index
				case y:
					yCore = core.Index
					ySock = socket.Index
				}
				if xCore != -1 && yCore != -1 {
					break
				}
			}
		}
	}

	return xCore != -1 && xCore == yCore && xSock == ySock
}

// Physical ID
type Socket struct {
	// Index of this Socket
	Index int64

	// The unique "Cores" in the Socket
	Cores Cores
}

func (s *Socket) AddCore(idx int64) {
	if s.Cores != nil && s.Cores.Get(idx) != nil {
		return
	}
	s.Cores = append(s.Cores, &Core{Index: idx, Processors: []int64{}})
}

type Sockets []*Socket

func (s Sockets) Get(idx int64) *Socket {
	for _, socket := range s {
		if socket.Index == idx {
			return socket
		}
	}
	return nil
}

// Cores per Socket
type Core struct {
	// Index of this Core
	Index int64

	// The processor IDs
	Processors Processors
}

func (c *Core) AddProcessor(idx int64) {
	if c.Processors != nil && c.Processors.Get(idx) != -1 {
		return
	}
	c.Processors = append(c.Processors, idx)
}

type Cores []*Core

func (c Cores) Get(idx int64) *Core {
	for _, core := range c {
		if core.Index == idx {
			return core
		}
	}
	return nil
}

type Processors []int64

func (p Processors) Get(idx int64) int64 {
	for _, index := range p {
		if index == idx {
			return index
		}
	}
	return -1
}

func (p Processors) Remove(l ...int64) Processors {
	var (
		rest    Processors
		removed = Processors(l)
	)
	for _, i := range p {
		if removed.Get(i) == -1 {
			rest = append(rest, i)
		}
	}
	return rest
}
