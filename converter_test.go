package ivgconv

import (
	"bytes"
	"fmt"

	"github.com/gio-eui/ivgconv/testdata"
	"github.com/google/gopacket/bytediff"
	"golang.org/x/exp/shiny/iconvg"

	"testing"
)

func TestFromFile(t *testing.T) {
	// Encode the SVG file as IconVG.

	type tcase struct {
		fn       string
		expected []byte
	}
	cases := []tcase{
		{
			fn:       "testdata/close.svg",
			expected: testdata.Close,
		},
		{
			fn:       "testdata/StarHalf.svg",
			expected: testdata.StarHalf,
		},
		{
			fn:       "testdata/abtesting.svg",
			expected: testdata.Abtesting,
		},
		{
			fn:       "testdata/i123.svg",
			expected: testdata.I123,
		},
		{
			fn:       "testdata/account_balance.svg",
			expected: testdata.Account_balance,
		},
		{
			fn:       "testdata/edit_off.svg",
			expected: testdata.Edit_off,
		},
		{
			fn:       "testdata/python.svg",
			expected: testdata.Python,
		},
	}

	for _, tc := range cases {
		t.Run(tc.fn, func(t *testing.T) {
			ivgData, err := FromFile(tc.fn)
			if err != nil {
				t.Fatal(err)
			}

			// Check that the IconVG data matches the expected output.

			show(t, ivgData)
			if !bytes.Equal(tc.expected, ivgData) {
				// disasm(t, tc.expected)
				disasm(t, ivgData)
				t.Fatal(bytediff.BashOutput.String(bytediff.Diff(tc.expected, ivgData)))
			}
		})
	}

}

func show(t *testing.T, ivgData []byte) {
	buf := &bytes.Buffer{}
	for _, b := range ivgData {
		fmt.Fprintf(buf, "0x%02x, ", b)
	}
	t.Log(buf.String())
}

type tD struct {
	t *testing.T
}

func (td *tD) Reset(m iconvg.Metadata) {
	td.t.Logf("Reset %#v", m)
}

func (td *tD) SetCSel(cSel uint8) {
	td.t.Log("SetCSel ", cSel)
}
func (td *tD) SetNSel(nSel uint8) {
	td.t.Log("SetNSel ", nSel)
}
func (td *tD) SetCReg(adj uint8, incr bool, c iconvg.Color) {
	td.t.Log("SetCReg ", adj, incr, c)
}
func (td *tD) SetNReg(adj uint8, incr bool, f float32) {
	td.t.Log("SetNReg ", adj, incr, f)
}
func (td *tD) SetLOD(lod0, lod1 float32) {
	td.t.Log("SetLOD ", lod0, lod1)
}

func (td *tD) StartPath(adj uint8, x, y float32) {
	td.t.Log("StartPath ", adj, x, y)
}
func (td *tD) ClosePathEndPath() {
	td.t.Log("ClosePathEndPath")
}
func (td *tD) ClosePathAbsMoveTo(x, y float32) {
	td.t.Log("ClosePathAbsMoveTo ", x, y)
}
func (td *tD) ClosePathRelMoveTo(x, y float32) {
	td.t.Log("ClosePathRelMoveTo ", x, y)
}

func (td *tD) AbsHLineTo(x float32) {
	td.t.Log("AbsHLineTo ", x)
}
func (td *tD) RelHLineTo(x float32) {
	td.t.Log("RelHLineTo ", x)
}
func (td *tD) AbsVLineTo(y float32) {
	td.t.Log("AbsVLineTo ", y)
}
func (td *tD) RelVLineTo(y float32) {
	td.t.Log("RelVLineTo ", y)
}
func (td *tD) AbsLineTo(x, y float32) {
	td.t.Log("AbsLineTo ", x, y)
}
func (td *tD) RelLineTo(x, y float32) {
	td.t.Log("RelLineTo ", x, y)
}
func (td *tD) AbsSmoothQuadTo(x, y float32) {
	td.t.Log("AbsSmoothQuadTo ", x, y)
}
func (td *tD) RelSmoothQuadTo(x, y float32) {
	td.t.Log("RelSmoothQuadTo ", x, y)
}
func (td *tD) AbsQuadTo(x1, y1, x, y float32) {
	td.t.Log("AbsQuadTo ", x1, y1, x, y)
}
func (td *tD) RelQuadTo(x1, y1, x, y float32) {
	td.t.Log("RelQuadTo ", x1, y1, x, y)
}
func (td *tD) AbsSmoothCubeTo(x2, y2, x, y float32) {
	td.t.Log("AbsSmoothCubeTo ", x2, y2, x, y)
}
func (td *tD) RelSmoothCubeTo(x2, y2, x, y float32) {
	td.t.Log("RelSmoothCubeTo ", x2, y2, x, y)
}
func (td *tD) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	td.t.Log("AbsCubeTo ", x1, y1, x2, y2, x, y)
}
func (td *tD) RelCubeTo(x1, y1, x2, y2, x, y float32) {
	td.t.Log("RelCubeTo ", x1, y1, x2, y2, x, y)
}
func (td *tD) AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	td.t.Log("AbsArcTo ", rx, ry, xAxisRotation, largeArc, sweep, x, y)
}
func (td *tD) RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	td.t.Log("RelArcTo ", rx, ry, xAxisRotation, largeArc, sweep, x, y)
}

func disasm(t *testing.T, ivgData []byte) {
	err := iconvg.Decode(&tD{t}, ivgData, nil)
	if err != nil {
		t.Log("error", err)
	}
}
