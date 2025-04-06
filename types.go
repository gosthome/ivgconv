package ivgconv

import (
	"encoding"
	"encoding/xml"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type SVG struct {
	Width   float32
	Height  float32
	ViewBox ViewBox
	Group
}

func (s *SVG) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "width", "height":
			val, err := strconv.ParseFloat(attr.Value, 32)
			if err != nil {
				return err
			}
			if attr.Name.Local == "width" {
				s.Width = float32(val)
			} else {
				s.Height = float32(val)
			}

		case "viewBox":
			err := s.ViewBox.UnmarshalXMLAttr(attr)
			if err != nil {
				return err
			}
		}
	}
	err := s.Group.UnmarshalXML(dec, start)
	if err != nil {
		return err
	}
	if s.ViewBox.Width == 0 && s.ViewBox.Height == 0 {
		s.ViewBox.MinX = 0
		s.ViewBox.MinY = 0
		s.ViewBox.Width = s.Width
		s.ViewBox.Height = s.Height
	}
	// Check if the viewbox has a valid paths or circles.
	if len(s.Elements) == 0 {
		return fmt.Errorf("no path or circle found in the SVG file")
	}

	return nil
}

type ViewBox struct {
	MinX   float32 `xml:"min-x,attr"`
	MinY   float32 `xml:"min-y,attr"`
	Width  float32 `xml:"width,attr"`
	Height float32 `xml:"height,attr"`
}

type Element interface {
	Equal(o Element) bool
}

type Path struct {
	D           string   `xml:"d,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity *float32 `xml:"fill-opacity,attr"`
	Opacity     *float32 `xml:"opacity,attr"`
}

func (p Path) Equal(o Element) bool {
	op, ok := o.(Path)
	if !ok {
		return false
	}
	return p.D == op.D && p.Fill == op.Fill
}

type Circle struct {
	Cx          float32  `xml:"cx,attr"`
	Cy          float32  `xml:"cy,attr"`
	R           float32  `xml:"r,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity *float32 `xml:"fill-opacity,attr"`
	Opacity     *float32 `xml:"opacity,attr"`
}

func (c Circle) Equal(o Element) bool {
	op, ok := o.(Circle)
	if !ok {
		return false
	}
	return c.Cx == op.Cx && c.Cy == op.Cy && c.R == op.R && c.Fill == op.Fill
}

type Ellipse struct {
	Cx          float32  `xml:"cx,attr"`
	Cy          float32  `xml:"cy,attr"`
	Rx          float32  `xml:"rx,attr"`
	Ry          float32  `xml:"ry,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity *float32 `xml:"fill-opacity,attr"`
	Opacity     *float32 `xml:"opacity,attr"`
}

func (c Ellipse) Equal(o Element) bool {
	op, ok := o.(Ellipse)
	if !ok {
		return false
	}
	return c.Cx == op.Cx && c.Cy == op.Cy && c.Rx == op.Rx && c.Ry == op.Ry && c.Fill == op.Fill
}

type Rect struct {
	X           float32  `xml:"x,attr"`
	Y           float32  `xml:"y,attr"`
	Width       float32  `xml:"width,attr"`
	Height      float32  `xml:"height,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity *float32 `xml:"fill-opacity,attr"`
	Opacity     *float32 `xml:"opacity,attr"`
}

func (r Rect) Equal(o Element) bool {
	or, ok := o.(Rect)
	if !ok {
		return false
	}
	return r.X == or.X &&
		r.Y == or.Y &&
		r.Width == or.Width &&
		r.Height == or.Height &&
		r.Fill == or.Fill
}

type Polygon struct {
	Points      PolygonPoints `xml:"points,attr"`
	Fill        string        `xml:"fill,attr"`
	FillOpacity *float32      `xml:"fill-opacity,attr"`
	Opacity     *float32      `xml:"opacity,attr"`
}

func (p Polygon) Equal(o Element) bool {
	op, ok := o.(Polygon)
	if !ok {
		return false
	}
	return slices.EqualFunc(
		p.Points, op.Points,
		func(p1 PolygonPoint, p2 PolygonPoint) bool {
			return p1.X == p2.X && p1.Y == p2.Y
		}) &&
		p.Fill == op.Fill
}

type PolygonPoint struct{ X, Y float32 }

type PolygonPoints []PolygonPoint

var spacelikerepl = []string{"\r", " ", "\n", " ", "\t", " ", "\f", " ", "\v", " "}

// func (p *PolygonPoints) UnmarshalXMLAttr(attr xml.Attr) error {}

// UnmarshalText implements encoding.TextUnmarshaler.
func (p *PolygonPoints) UnmarshalText(btext []byte) error {
	*p = make(PolygonPoints, 0)
	for s := range strings.FieldsSeq(string(btext)) {
		x, y, ok := strings.Cut(s, ",")
		if !ok {
			return errors.New("invalid coordinate pair")
		}
		xf, err := strconv.ParseFloat(x, 32)
		if err != nil {
			return err
		}
		yf, err := strconv.ParseFloat(y, 32)
		if err != nil {
			return err
		}
		*p = append(*p, struct {
			X float32
			Y float32
		}{X: float32(xf), Y: float32(yf)})
	}
	return nil
}

var _ encoding.TextUnmarshaler = (*PolygonPoints)(nil)

type Group struct {
	Elements []Element

	ExcludeElements []ElementPredicate
}

func (g *Group) exclude(e Element) bool {
	return slices.ContainsFunc(g.ExcludeElements, func(excl ElementPredicate) bool {
		return excl(e)
	})
}

func (g *Group) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	for {
		t, _ := dec.Token()
		if t == nil {
			break
		}

		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "path":
				var p Path
				err = dec.DecodeElement(&p, &t)
				if err != nil {
					return err
				}
				p.D = strings.Replace(p.D, ",", " ", -1)
				if g.exclude(p) {
					continue
				}
				g.Elements = append(g.Elements, p)
			case "circle":
				var c Circle
				err = dec.DecodeElement(&c, &t)
				if err != nil {
					return err
				}
				if g.exclude(c) {
					continue
				}
				g.Elements = append(g.Elements, c)
			case "rect":
				var r Rect
				err = dec.DecodeElement(&r, &t)
				if err != nil {
					return err
				}
				if g.exclude(r) {
					continue
				}
				g.Elements = append(g.Elements, Path{
					D:           fmt.Sprintf("M%[1]f %[2]f h %[3]f v %[4]f h -%[3]f Z", r.X, r.Y, r.Width, r.Height),
					Fill:        r.Fill,
					FillOpacity: r.FillOpacity,
					Opacity:     r.Opacity,
				})
			case "polygon", "polyline":
				var pol Polygon
				err = dec.DecodeElement(&pol, &t)
				if err != nil {
					return err
				}
				if g.exclude(pol) {
					continue
				}
				b := strings.Builder{}
				for i, p := range pol.Points {
					if i == 0 {
						b.WriteString("M")
					} else {
						b.WriteString(" L")
					}
					fmt.Fprintf(&b, "%f %f", p.X, p.Y)
				}
				if t.Name.Local != "polyline" {
					b.WriteString(" Z")
				}
				g.Elements = append(g.Elements, Path{
					D:           b.String(),
					Fill:        pol.Fill,
					FillOpacity: pol.FillOpacity,
					Opacity:     pol.Opacity,
				})

			case "ellipse":
				var e Ellipse
				err = dec.DecodeElement(&e, &t)
				if err != nil {
					return err
				}
				if g.exclude(e) {
					continue
				}

				g.Elements = append(g.Elements, Path{
					D:           fmt.Sprintf("M%f %fa%f %f 0 1 0 %f 0a%f %f 0 1 0 -%f 0 Z", e.Cx-e.Rx, e.Cy, e.Rx, e.Ry, e.Rx*2, e.Rx, e.Ry, e.Rx*2),
					Fill:        e.Fill,
					FillOpacity: e.FillOpacity,
					Opacity:     e.Opacity,
				})

			case "g":
				var gr Group
				gr.ExcludeElements = g.ExcludeElements
				err = dec.DecodeElement(&gr, &t)
				if err != nil {
					return err
				}
				g.Elements = append(g.Elements, gr.Elements...)

			case "title": // ignored elements
				err = dec.DecodeElement(&struct{}{}, &t)
				if err != nil {
					return err
				}
				continue
			default:
				return fmt.Errorf("unexpected XML element: %s", t.Name.Local)
			}
		case xml.EndElement:
			if t.Name.Local != start.Name.Local {
				return errors.New("unexpected end element: " + t.Name.Local)
			}
		}
	}
	return nil
}

var _ xml.Unmarshaler = (*Group)(nil)

// UnmarshalXMLAttr implements the xml.UnmarshalerAttr interface.
func (vb *ViewBox) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Name.Local != "viewBox" {
		return nil
	}
	if _, err := fmt.Sscanf(attr.Value, "%f %f %f %f", &vb.MinX, &vb.MinY, &vb.Width, &vb.Height); err != nil {
		return err
	}
	return nil
}
