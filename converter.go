package ivgconv

import (
	"encoding/xml"
	"os"
)

type ElementPredicate func(e Element) bool

func ExcludeExact(excl Element) ElementPredicate {
	return func(e Element) bool {
		return e.Equal(excl)
	}
}

// ConverterOptions contains options for the SVG to IconVG converter.
type ConverterOptions struct {
	// OutputSize is the size of the IconVG output image.
	OutputSize float32
	// Excludes is a list of elements to exclude from the IconVG image.
	Excludes []ElementPredicate
}

// Option is a function that configures a ConverterOptions.
type Option func(*ConverterOptions)

// FromFile encodes the SVG file as IconVG.
func FromFile(filepath string, options ...Option) ([]byte, error) {
	// Read the SVG file.
	svgData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// Encode the SVG file content as IconVG.
	return FromContent(svgData, options...)
}

// FromContent encodes the SVG file content as IconVG.
func FromContent(content []byte, options ...Option) ([]byte, error) {
	// Set the default converter options.
	opts := ConverterOptions{
		OutputSize: 48,
		Excludes: []ElementPredicate{
			func(e Element) bool {
				fill := ""
				switch e := e.(type) {
				case Path:
					fill = e.Fill
				case Rect:
					fill = e.Fill
				case Circle:
					fill = e.Fill
				case Ellipse:
					fill = e.Fill
				case Polygon:
					fill = e.Fill
				}
				return fill == "none"
			},
			ExcludeExact(Path{D: "M0 0h24v24H0V0z"}),
			// Matches <path d="M0 0h24v24H0z" fill="none"/>
			// Path{D: "M0 0h24v24H0z", Fill: "none"},
			// Path{D: "M0 0h24v24H0zm0 0h24v24H0z", Fill: "none"},
			// Matches <path d="M0 0H24V24H0z" fill="none"/>
			// Path{D: "M0 0H24V24H0z", Fill: "none"},
			// Rect{X: 0, Y: 0, Width: 24, Height: 24, Fill: "none"},
		},
	}
	// Set the converter options.
	for _, option := range options {
		option(&opts)
	}
	// Parse the SVG file.
	var svg SVG
	svg.ExcludeElements = opts.Excludes
	if err := xml.Unmarshal(content, &svg); err != nil {
		return nil, err
	}
	// Encode the SVG file as IconVG.
	return parseSVG(svg, opts)
}

// WithOutputSize sets the size of the IconVG image.
func WithOutputSize(outputSize float32) Option {
	return func(opts *ConverterOptions) {
		opts.OutputSize = outputSize
	}
}

// WithReplaceExcludedElements sets the list of exact elements to exclude from the IconVG image.
func WithReplaceExcludedElements(excludes []Element) Option {
	return func(opts *ConverterOptions) {
		newRules := make([]ElementPredicate, len(excludes))
		for _, e := range excludes {
			newRules = append(newRules, ExcludeExact(e))
		}
		opts.Excludes = newRules
	}
}

// AddExcludedElements appends the list of exact elements to exclude from the IconVG image.
func AddExcludedElements(excludes []Element) Option {
	return func(opts *ConverterOptions) {
		newRules := make([]ElementPredicate, len(excludes))
		for _, e := range excludes {
			newRules = append(newRules, ExcludeExact(e))
		}
		opts.Excludes = append(opts.Excludes, newRules...)
	}
}
