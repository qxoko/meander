/*
	Meander
	A portable Fountain utility for production writing
	Copyright (C) 2022-2023 Harley Denham

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

const final_draft_extension = ".fdx"

// we discard _a lot_ of crufty data that's
// not useful to us in Fountain, but try to
// safely apply stuff we can understand

type XML_FinalDraft struct {
	XMLName xml.Name         `xml:"FinalDraft"`
	Content []*XML_Paragraph `xml:"Content>Paragraph"`
	Title   []*XML_Paragraph `xml:"TitlePage>Content>Paragraph"`

	// potentially need to look at capturing
	// the HeaderAndFooter attributes for
	// more information, such as visibility
	// and starting pages.
	Header []*XML_Paragraph `xml:"HeaderAndFooter>Header>Paragraph"`
	Footer []*XML_Paragraph `xml:"HeaderAndFooter>Footer>Paragraph"`
}

type XML_Paragraph struct {
	XMLName   xml.Name     `xml:"Paragraph"`
	Type      string       `xml:"Type,attr"`
	Number    string       `xml:"Number,attr"`
	Alignment string       `xml:"Alignment,attr"`
	Chunks    []*XML_Chunk `xml:"Text"`
}

// the entire screenplay is actually stored in a mixed-type
// array because XML is a garbage format designed by criminals

// to combat this **we literally find/replace** the lesser-used
// types into the major type in the input stream, then capture
// everything in this XML_Chunk, which contains all the known
// attributes that we need to collect from these different types

// because these attributes are guaranteed to only be populated
// in relevant cases (because they were originally distinct
// types), we can later re-identify them by checking
// the "signature" of which attributes are populated

// we do this because XML is a garbage format designed by
// criminals, of course, but also because the Go XML parser
// will only really reflect into interfaces and it's nasty

type XML_Chunk struct {
	XMLName xml.Name `xml:"Text"`
	Style   string   `xml:"Style,attr"` // <Text>         attribute
	Label   string   `xml:"Type,attr"`  // <DynamicLabel> attribute
	Text    string   `xml:",chardata"`  // <Text>         content
}

func parse_final_draft_xml(source_file string) (*XML_FinalDraft, bool) {
	byte_stream, err := os.ReadFile(source_file)
	if err != nil {
		panic(err)
	}

	// in order to capture <DynamicLabels> as part of <Text> arrays
	// we just... swap the literal text like a absolute maniac.
	// because the attribute signatures are different, we can separately
	// identify them within the struct -- see above
	byte_stream = bytes.ReplaceAll(byte_stream, []byte("DynamicLabel"), []byte("Text"))

	document := &XML_FinalDraft{}

	err = xml.Unmarshal(byte_stream, document)

	if err != nil {
		panic(err)
		return nil, false
	}

	return document, true
}

func convert_final_draft(config *config) {
	data, ok := parse_final_draft_xml(fix_path(config.source_file))

	if !ok {
		eprintln("failed to parse Final Draft document")
		return
	}

	buffer := strings.Builder{}
	buffer.Grow(len(data.Content) * 128)

	// title page
	{
		// because Final Draft title pages are manually placed
		// we simply assign the "central" items to the "title"
		// section, and all others to the "notes" section
		// which leaves it compatible with most parsers
		title_buffer := strings.Builder{}
		title_buffer.Grow(len(data.Title) * 128)
		title_buffer.WriteString("title:")

		info_buffer := strings.Builder{}
		info_buffer.Grow(len(data.Title) * 128)
		info_buffer.WriteString("notes:")

		for _, paragraph := range data.Title {
			has_text := false

			// go's XML will fill the array with empty items
			// so we have to check if the paragraph actually
			// has any data in it

			// encoding/xml's "omitempty" does **nothing**
			// to assist, so we're just not using it.

			// big_fish.fdx (https://fountain.io) has
			// about 50 title page "paragraphs" in it
			// that are totally empty other than random
			// stylings, so this (https://xkcd.com/2109)
			// is hilariously relevant
			for _, chunk := range paragraph.Chunks {
				if len(chunk.Label) != 0 || len(chunk.Text) != 0 {
					has_text = true
					break
				}
			}

			if !has_text {
				continue
			}

			// if centered, we assume "title"
			if paragraph.Alignment == "Center" {
				title_buffer.WriteString("\n\t")

				for _, chunk := range paragraph.Chunks {
					if chunk.Text != "" {
						title_buffer.WriteString(strings.TrimSpace(chunk.Text))
					}
				}

			// ...otherwise assign to "info"
			} else {
				info_buffer.WriteString("\n\t")

				for _, chunk := range paragraph.Chunks {
					if chunk.Text != "" {
						info_buffer.WriteString(strings.TrimSpace(chunk.Text))
					}
				}
			}
		}

		buffer.WriteString(title_buffer.String())
		buffer.WriteRune('\n')

		buffer.WriteString(info_buffer.String())
		buffer.WriteRune('\n')
	}

	// header + footer
	{
		for _, paragraph := range data.Header {
			has_text := false

			for _, chunk := range paragraph.Chunks {
				if len(consume_whitespace(chunk.Label)) != 0 || len(consume_whitespace(chunk.Text)) != 0 {
					has_text = true
					break
				}
			}

			if has_text {
				buffer.WriteString("\n{{header:")
				buffer.WriteString(write_chunks(paragraph.Chunks, false))
				buffer.WriteString("}}")
			}
		}

		for _, paragraph := range data.Footer {
			has_text := false

			for _, chunk := range paragraph.Chunks {
				if len(consume_whitespace(chunk.Label)) != 0 || len(consume_whitespace(chunk.Text)) != 0 {
					has_text = true
					break
				}
			}

			if has_text {
				buffer.WriteString("\n{{footer:")
				buffer.WriteString(write_chunks(paragraph.Chunks, false))
				buffer.WriteString("}}")
			}
		}
	}

	// base content
	{
		for _, paragraph := range data.Content {
			// same as before, ignore empty paragraphs
			has_text := false

			for _, chunk := range paragraph.Chunks {
				if len(chunk.Label) != 0 || len(chunk.Text) != 0 {
					has_text = true
					break
				}
			}

			if !has_text {
				continue
			}

			// handle individual cases that need it differently
			// this is **not exhaustive** and not fully tested
			switch paragraph.Type {
			case "Scene Heading":
				buffer.WriteString("\n\n")
				text := write_chunks(paragraph.Chunks, true)

				// @todo write a dedicated scene validator
				n := strings.IndexRune(text, '.')

				// force scenes if we know Fountain wouldn't identify them
				if n < 0 || !valid_scene[strings.ToLower(clean_string(text[:n]))] {
					buffer.WriteRune('.')
				}

				buffer.WriteString(text)

				// add the scene number if it's encoded
				if paragraph.Number != "" {
					buffer.WriteString(fmt.Sprintf(" #%s#", paragraph.Number))
				}

			case "Character":
				buffer.WriteString("\n\n")
				text := write_chunks(paragraph.Chunks, true)

				// force characters if we know Fountain wouldn't identify them
				if !is_valid_character(text) {
					buffer.WriteRune('@')
				}

				buffer.WriteString(text)

			case "Dialogue", "Parenthetical":
				buffer.WriteRune('\n') // no space between char + dialogue
				buffer.WriteString(write_chunks(paragraph.Chunks, false))

			default:
				buffer.WriteString("\n\n")
				buffer.WriteString(write_chunks(paragraph.Chunks, false))
			}
		}
	}

	ok = write_file(fix_path(config.output_file), []byte(buffer.String()))
	if !ok {
		eprintln("failed to write", config.output_file)
	}
}

// lookup table for FD > Fountain markers
var final_draft_styles = map[string]string{
	"Italic":    "*",
	"Bold":      "**",
	"Underline": "_",
}

// loop through chunks of text, and based on their styles
// (of which we do not know all of the available ones yet)
// write in the relevant Fountain ones as safely as possible
func write_chunks(input []*XML_Chunk, force_uppercase bool) string {
	buffer := strings.Builder{}
	buffer.Grow(len(input) * 128)

	for _, chunk := range input {
		if force_uppercase || strings.Contains(chunk.Style, "AllCaps") {
			chunk.Text = strings.ToUpper(chunk.Text)
		}

		opening := ""
		closing := ""

		if len(chunk.Style) != 0 {
			styles := strings.Split(chunk.Style, "+")

			for _, s := range styles {
				if x, ok := final_draft_styles[s]; ok {
					opening = opening + x
					closing = x + closing
				}
				// println("debug: missed a final draft thing", s)
			}
		}

		if len(chunk.Label) != 0 {
			switch chunk.Label {
			case "Page #":
				buffer.WriteString("%p")
			case "Last Revised":
				buffer.WriteString("{{timestamp}}")
			}
			continue
		}

		buffer.WriteString(opening)
		buffer.WriteString(chunk.Text)
		buffer.WriteString(closing)
	}

	return buffer.String()
}
