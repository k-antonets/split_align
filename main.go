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

func main() {
	flag.Parse()

}
