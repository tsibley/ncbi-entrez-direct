// ===========================================================================
//
//                            PUBLIC DOMAIN NOTICE
//            National Center for Biotechnology Information (NCBI)
//
//  This software/database is a "United States Government Work" under the
//  terms of the United States Copyright Act. It was written as part of
//  the author's official duties as a United States Government employee and
//  thus cannot be copyrighted. This software/database is freely available
//  to the public for use. The National Library of Medicine and the U.S.
//  Government do not place any restriction on its use or reproduction.
//  We would, however, appreciate having the NCBI and the author cited in
//  any work or product based on this material.
//
//  Although all reasonable efforts have been taken to ensure the accuracy
//  and reliability of the software and data, the NLM and the U.S.
//  Government do not and cannot warrant the performance or results that
//  may be obtained by using this software or data. The NLM and the U.S.
//  Government disclaim all warranties, express or implied, including
//  warranties of performance, merchantability or fitness for any particular
//  purpose.
//
// ===========================================================================
//
// File Name:  s2p.go
//
// Author:  Jonathan Kans
//
// ==========================================================================

package main

// demonstrates simple build with Go modules and scripted cross-compilation,
// added to provide codeathon participants with a minimal Go language template
// and demonstration of graphic components for displaying pangenome networks

/*
cd edirect
go mod init edirect
go mod tidy

go build xtract.go common.go
go build rchive.go common.go

go build j2x.go
go build t2x.go

echo "darwin amd64 Darwin linux amd64 Linux windows 386 CYGWIN_NT linux arm ARM" |
xargs -n 3 sh -c 'env GOOS="$0" GOARCH="$1" go build -o symbols."$2" s2p.go'
*/

import (
	"flag"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"math/rand"
	"os"
	"time"
)

// global variables for command-line arguments
var (
	pngPath string
	numObj  int
	reseed  bool
)

func DrawPicture() {

	const (
		picWidth  = 900
		picHeight = 600
		maxEdge   = 50
	)

	colorTriplets := []float64{
		// permutations 4 values 3 at a time, minor manual tweaks, removals, and rearrangements
		 64, 128, 192, 255, 192,  64,  64, 128, 255,  64, 192, 128, 255, 128,  64,
		192,  64, 128,  64, 255, 128, 128,  64, 192,  64, 255, 192, 255,  64, 192,
		128,  64, 255, 128, 192,  64, 128, 192, 255, 169, 196, 181, 192,  64, 255,
		192, 128,  64, 255,  64, 128, 192, 128, 255, 214, 198, 222, 255, 128, 192,
		255, 192, 128,
	}

	colorIndex := 0

	nextColor := func() float64 {

		colorIndex = colorIndex % len(colorTriplets)
		val := colorTriplets[colorIndex]
		colorIndex++

		return val
	}

	dc := gg.NewContext(picWidth, picHeight)

	font, err := truetype.Parse(gobold.TTF)
	if err != nil {
		os.Exit(1)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 16})

	dc.SetFontFace(face)
	dc.SetRGB(127, 0, 127)
	dc.DrawStringAnchored("Unknown Miró or Random Symbols?", picWidth/2, picHeight-30, 0.5, 0.5)

	if reseed {
		rand.Seed(int64(time.Now().Nanosecond() % (1e9 - 1)))
	}

	randRange := func(low, high int) float64 {
		return float64(rand.Float32()*float32(high-low) + float32(low))
	}

	nextPoint := func() (float64, float64) {
		return randRange(maxEdge, picWidth-2*maxEdge), randRange(maxEdge, picHeight-2*maxEdge)
	}

	drawObject := func(i int) {

		dc.Push()
		defer dc.Pop()

		r, g, b := nextColor(), nextColor(), nextColor()
		dc.SetRGB(r, g, b)

		dc.SetLineWidth(randRange(2, 5))

		x, y := nextPoint()
		w, h := randRange(5, maxEdge), randRange(5, maxEdge)

		switch i % 7 {
		case 0, 3:
			dc.DrawRectangle(x, y, w, h)
		case 1, 5:
			dc.SetRGBA(r, g, b, randRange(1, 9)/10)
			dc.DrawCircle(x, y, w)
		case 2, 4:
			dc.RotateAbout(gg.Radians(randRange(0, 45)), x, y)
			dc.DrawEllipse(x, y, w, h)
		case 6:
			p, q := nextPoint()
			dc.DrawLine(x, y, p, q)
			dc.Stroke()
			return
		default:
		}

		switch i % 5 {
		case 0, 3:
			dc.Fill()
		case 1, 2, 4:
			dc.Stroke()
		default:
		}
	}

	for i := 0; i < numObj; i++ {
		drawObject(i)
	}

	dc.SavePNG(pngPath)
}

// init functions automatically run before program starts
func init() {
	flag.StringVar(&pngPath, "o", "random.png", "png output file")
	flag.IntVar(&numObj, "n", 30, "number of objects")
	flag.BoolVar(&reseed, "s", false, "initialize with random seed")
}

func main() {

	// process command-line arguments
	flag.Parse()

	DrawPicture()
}
