package splitAlign

import (
	"flag"
	"github.com/shenwei356/bio/seq"
	"github.com/shenwei356/bio/seqio/fastx"
)

var (
	align  = flag.String("align", "", "File with full alignment")
	ref    = flag.String("ref", "", "File with reference sequences, how to split")
	outdir = flag.String("out", "", "Output directory")
)

type SeqEntry struct {
	Name string
	Len  int
}

type SeqSplitter []*SeqEntry

func (s SeqSplitter) Cumulate() []int {
	l := len(s)
	r := make([]int, l)

	for i, se := range s {
		if i == 0 {
			r[i] = 1
		} else {
			r[i] = r[i-1] + se.Len - 1
		}
	}

	return r
}

func (s SeqSplitter) Split(seq *fastx.Record) []*fastx.Record {
	return nil
}

func main() {
	flag.Parse()

}
