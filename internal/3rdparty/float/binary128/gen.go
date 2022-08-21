//go:build ignore
// +build ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"text/template"
)

func main() {
	var out string
	flag.StringVar(&out, "o", "extra_test.go", "test cases output path")
	flag.Parse()
	if err := dumpTest(out); err != nil {
		log.Fatalf("%+v", err)
	}
}

func dumpTest(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	t, err := template.ParseFiles("extra_test.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}
	// Use deterministic source for pseudo-random nunmbers.
	rand.Seed(1234)
	// Randomize the exponent since the number of exponent, mantissa combinations
	// otherwise become huge.
	var exps []int
	const nrandExps = 64
	for i := 0; i < nrandExps; i++ {
		// exponent bits: 0x0001 - 0x7FFE
		exp := rand.Intn(0x7FFE) + 1
		exps = append(exps, exp)
	}
	sort.Ints(exps)
	// Randomize the mantissa since we cannot check 112 (48 + 64) bits
	// exhaustively.
	var mants []MantBits
	const nrandMants = 512
	for i := 0; i < nrandMants; i++ {
		// 48 bits.
		a := rand.Uint64() & 0xFFFFFFFFFFFF
		// 64 bits.
		b := rand.Uint64()
		mant := MantBits{a: a, b: b}
		mants = append(mants, mant)
	}
	sort.Slice(mants, func(i, j int) bool {
		if mants[i].a < mants[j].a {
			return true
		}
		return mants[i].b < mants[j].b
	})
	data := map[string][]string{
		"normalized":   getNormalized(exps, mants),
		"denormalized": getDenormalized(mants),
	}
	if err := t.Execute(f, data); err != nil {
		return errors.New(err)
	}
	return nil
}

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 113
	// exponent bias.
	bias = 16383
)

type MantBits struct {
	// 48 bits.
	a uint64
	// 64 bits.
	b uint64
}

func getNormalized(exps []int, mants []MantBits) []string {
	var ns []string
	// normalized
	//
	// exponent bits: 0x0001 - 0x7FFE
	//
	//    (-1)^signbit * 2^(exp-16383) * 1.mant_2
	const lead = 1
	for signbit := 0; signbit <= 1; signbit++ {
		sign := "+"
		if signbit == 1 {
			sign = "-"
		}
		for _, exp := range exps {
			exponent := exp - bias
			// mantissa bits: 112 (48 + 64) bits
			for _, mantBits := range mants {
				mant := fmt.Sprintf("%048b%064b", mantBits.a, mantBits.b)
				s := fmt.Sprintf("%s0b%d.%sp%d", sign, lead, mant, exponent)
				m, _, err := big.ParseFloat(s, 0, precision, big.ToNearestEven)
				if err != nil {
					panic(err)
				}
				want := m.Text('g', 35)
				want64, acc64 := m.Float64()
				a := uint64(signbit) << 63
				a |= uint64(exp) << 48
				a |= mantBits.a
				b := mantBits.b
				var n string
				switch {
				// Compare floating-point bits of want64, as otherwise +0 == -0
				case math.Float64bits(want64) == math.Float64bits(math.Copysign(0, -1)):
					// -zero
					n = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Copysign(0, -1), acc64: big.%v}, // %s", a, b, want, acc64, s)
				case want64 == math.Inf(+1):
					// +inf
					n = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Inf(+1), acc64: big.%v}, // %s", a, b, want, acc64, s)
				case want64 == math.Inf(-1):
					// -inf
					n = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Inf(-1), acc64: big.%v}, // %s", a, b, want, acc64, s)
				default:
					n = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: %v, acc64: big.%v}, // %s", a, b, want, want64, acc64, s)
				}
				ns = append(ns, n)
			}
		}
	}
	return ns
}

func getDenormalized(mants []MantBits) []string {
	var ds []string
	// denormalized
	//
	// exponent bits: 0x0000
	//
	//    (-1)^signbit * 2^(-14) * 0.mant_2
	const lead = 0
	for signbit := 0; signbit <= 1; signbit++ {
		sign := "+"
		if signbit == 1 {
			sign = "-"
		}
		const exp = 0x0000
		exponent := exp - bias + 1
		// mantissa bits: 112 (48 + 64) bits
		for _, mantBits := range mants {
			mant := fmt.Sprintf("%048b%064b", mantBits.a, mantBits.b)
			s := fmt.Sprintf("%s0b%d.%sp%d", sign, lead, mant, exponent)
			m, _, err := big.ParseFloat(s, 0, precision, big.ToNearestEven)
			if err != nil {
				panic(err)
			}
			want := m.Text('g', 35)
			want64, acc64 := m.Float64()
			a := uint64(signbit) << 63
			a |= uint64(exp) << 48
			a |= mantBits.a
			b := mantBits.b
			var d string
			switch {
			// Compare floating-point bits of want64, as otherwise +0 == -0
			case math.Float64bits(want64) == math.Float64bits(math.Copysign(0, -1)):
				// -zero
				d = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Copysign(0, -1), acc64: big.%v}, // %s", a, b, want, acc64, s)
			case want64 == math.Inf(+1):
				// +inf
				d = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Inf(+1), acc64: big.%v}, // %s", a, b, want, acc64, s)
			case want64 == math.Inf(-1):
				// -inf
				d = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: math.Inf(-1), acc64: big.%v}, // %s", a, b, want, acc64, s)
			default:
				d = fmt.Sprintf("{a: 0x%016X, b: 0x%016X, want: %q, want64: %v, acc64: big.%v}, // %s", a, b, want, want64, acc64, s)
			}
			ds = append(ds, d)
		}
	}
	return ds
}
