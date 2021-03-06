package pcg

// PCG Random Number Generation
// Developed by Melissa O'Neill <oneill@pcg-random.org>
// Paper and details at http://www.pcg-random.org
// Ported to Go by Michael Jones <michael.jones@gmail.com>

const (
	pcg32State      = 0x853c49e6748fea9b //  9600629759793949339
	pcg32Increment  = 0xda3e39cb94b95bdb // 15726070495360670683
	pcg32Multiplier = 0x5851f42d4c957f2d //  6364136223846793005
)

type PCG32 struct {
	state     uint64
	increment uint64
}

func NewPCG32() *PCG32 {
	return &PCG32{pcg32State, pcg32Increment}
}

func (p *PCG32) Seed(state, sequence uint64) *PCG32 {
	p.increment = (sequence << 1) | 1
	p.state = (state+p.increment)*pcg32Multiplier + p.increment
	return p
}

func (p *PCG32) Random() uint32 {
	// Advance 64-bit linear congruential generator to new state
	oldState := p.state
	p.state = oldState*pcg32Multiplier + p.increment

	// Confuse and permute 32-bit output from old state
	xorShifted := uint32(((oldState >> 18) ^ oldState) >> 27)
	rot := uint32(oldState >> 59)
	return (xorShifted >> rot) | (xorShifted << ((-rot) & 31))
}

func (p *PCG32) Bounded(bound uint32) uint32 {
	if bound == 0 {
		return 0
	}
	threshold := -bound % bound
	for {
		r := p.Random()
		if r >= threshold {
			return r % bound
		}
	}
}

// as in int31n, go/src/math/rand/rand.go
// this function uses a single division in the worst case
func (p *PCG32) FastBounded(bound uint32) uint32 {
	v := p.Random()
	prod := uint64(v) * uint64(bound)
	low := uint32(prod)
	if low < uint32(bound) {
		thresh := uint32(-bound) % uint32(bound)
		for low < thresh {
			v = p.Random()
			prod = uint64(v) * uint64(bound)
			low = uint32(prod)
		}
	}
	return uint32(prod >> 32)
}


func (p *PCG32) Advance(delta uint64) *PCG32 {
	p.state = p.advanceLCG64(p.state, delta, pcg32Multiplier, p.increment)
	return p
}

func (p *PCG32) Retreat(delta uint64) *PCG32 {
	return p.Advance(-delta)
}

func (p *PCG32) advanceLCG64(state, delta, curMult, curPlus uint64) uint64 {
	accMult := uint64(1)
	accPlus := uint64(0)
	for delta > 0 {
		if delta&1 != 0 {
			accMult *= curMult
			accPlus = accPlus*curMult + curPlus
		}
		curPlus = (curMult + 1) * curPlus
		curMult *= curMult
		delta /= 2
	}
	return accMult*state + accPlus
}
