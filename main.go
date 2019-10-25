package main

import (
	"flag"
	"fmt"
	"github.com/shenwei356/bio/seq"
	"github.com/shenwei356/bio/seqio/fastx"
	"github.com/shenwei356/xopen"
	"io"
	"path"
	"strings"
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

func NewSeqEntry(record *fastx.Record) *SeqEntry {
	return &SeqEntry{
		Name: string(record.Name),
		Len:  record.Seq.Length(),
	}
}

type SeqSplitter []*SeqEntry

func (s SeqSplitter) Cumulate() []int {
	l := len(s)
	r := make([]int, l+1)
	r[0] = 1

	for i, se := range s {
		r[i+1] = r[i] + se.Len - 1
	}

	return r
}

func (s SeqSplitter) Split(old_record *fastx.Record) ([]*fastx.Record, error) {
	r := make([]*fastx.Record, len(s))
	cum := s.Cumulate()

	for i, se := range s {
		name := string(old_record.Name)
		new_name := strings.Join([]string{name, se.Name}, "_")

		rec, err := fastx.NewRecordWithSeq(old_record.ID, []byte(new_name), old_record.Seq.SubSeq(cum[i], cum[i+1]))
		if err != nil {
			return nil, err
		}

		r[i] = rec
	}

	return r, nil
}

func (ss SeqSplitter) SplitAndWrite(outdir string, records <-chan *fastx.Record) error {

	seqdir := make(map[string][]*fastx.Record)

	for r := range records {
		seqs, e := ss.Split(r)
		if e != nil {
			return e
		}

		for i, s := range seqs {
			seqdir[ss[i].Name] = append(seqdir[ss[i].Name], s)
		}
	}

	writeEntity := func(name string, records []*fastx.Record) error {
		filename := name + ".fasta"
		filename = path.Join(outdir, filename)

		f, e := xopen.Wopen(filename)
		if e != nil {
			return e
		}
		defer f.Close()

		for _, r := range records {
			r.FormatToWriter(f, 0)
		}
		return nil
	}

	for _, se := range ss {
		if err := writeEntity(se.Name, seqdir[se.Name]); err != nil {
			return err
		}
	}

	return nil
}

func NewSplitter(filename string) (SeqSplitter, error) {
	reader, err := fastx.NewReader(seq.DNAredundant, filename, "")
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	ss := []*SeqEntry{}

	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ss = append(ss, NewSeqEntry(r))
	}

	return ss, nil
}

func main() {
	flag.Parse()

	ss, err := NewSplitter(*ref)
	if err != nil {
		panic(err)
	}

	fmt.Println("splitter created")

	rchan := make(chan *fastx.Record)

	go func(filename string, out chan<- *fastx.Record) {
		reader, err := fastx.NewReader(seq.DNAredundant, filename, "")
		if err != nil {
			panic(err)
		}

		for {
			r, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			out <- r
		}
		close(rchan)
	}(*align, rchan)

	err = ss.SplitAndWrite(*outdir, rchan)
	if err != nil {
		panic(err)
	}

	fmt.Println("program finished")
}
