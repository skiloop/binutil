package main

import (
	"flag"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/skiloop/gbinutil/binfile"
	"os"
)

type ListCmd struct {
	Input  string `arg:"" help:"input file name"`
	Offset int64  `arg:"" optional:"" help:"start document position" default:"0"`
}

type ReadCmd struct {
	Input  string `arg:"" help:"input file name"`
	Offset int64  `arg:"" optional:"" help:"start position" default:"0"`
}
type CountCmd struct {
	Input  string `arg:"" help:"input file name"`
	Offset int64  `arg:"" optional:"" help:"start position" default:"0"`
}

var client struct {
	CompressType string   `help:"compression type, options are gzip, bz2 and zip, default is gzip" enum:"gzip,bz2,zip" default:"gzip"`
	Verbose      bool     `short:"v" help:"verbose" default:"false"`
	List         ListCmd  `cmd:"" aliases:"l,ls" help:"List documents from position."`
	Read         ReadCmd  `cmd:"" aliases:"r,ra" help:"Read document file in bin file at position"`
	Count        CountCmd `cmd:"" aliases:"c" help:"count document file in bin file from position"`
}

func listDocs() {
	ct, ok := binfile.CompressTypes[client.CompressType]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown compression type %s\n", client.CompressType)
		return
	}

	br := binfile.NewBinReader(client.List.Input, ct)
	if br != nil {
		br.List(client.List.Offset, nil)
		return
	}
	fmt.Fprintf(os.Stderr, "file not found: %s\n", client.List.Input)
}

func readDoc() {
	ct, ok := binfile.CompressTypes[client.CompressType]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown compression type %s\n", client.CompressType)
		return
	}
	br := binfile.NewBinReader(client.Read.Input, ct)
	if br == nil {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", client.Read.Input)
		return
	}
	doc, err := br.ReadAt(client.Read.Offset, true)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(doc.Content)
		fmt.Println(doc.Key)
	}
}

func countDocs() {
	ct, ok := binfile.CompressTypes[client.CompressType]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown compression type %s\n", client.CompressType)
		return
	}
	br := binfile.NewBinReader(client.Count.Input, ct)
	if br == nil {
		fmt.Fprintf(os.Stderr, "file not found: %s\n", client.Count.Input)
		return
	}
	count, err := br.Count(client.Count.Offset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file read error: %v\n", err)
	} else {
		fmt.Printf("%d", count)
	}
}

func main() {
	ctx := kong.Parse(&client)
	binfile.Verbose = client.Verbose
	switch ctx.Command() {
	case "list <input>", "list <input> <offset>":
		listDocs()
		break
	case "read <input>", "read <input> <offset>":
		readDoc()
		break
	case "count <input>", "count <input> <offset>":
		countDocs()
		break
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", ctx.Command())
		flag.Usage()
	}
}
