package pdf

import (
	"bytes"
	"fmt"
	"jobsearch/internal/parser"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	pageWidth    = 612.0
	pageHeight   = 792.0
	marginLeft   = 46.8
	marginRight  = 46.8
	marginTop    = 43.2
	marginBottom = 43.2
	contentWidth = pageWidth - marginLeft - marginRight

	fontRegular = "eb-garamond-regular"
	fontBold    = "eb-garamond-bold"

	fontSizeName    = 18
	fontSizeSection = 13
	fontSizeBody    = 11
	fontSizeContact = 11

	lineHeight     = 16.0
	sectionSpacing = 8.0
)

type Generator struct {
	pdf    *gopdf.GoPdf
	y      float64
	resume *parser.Resume
}

func New(resume *parser.Resume) *Generator {
	return &Generator{
		pdf:    &gopdf.GoPdf{},
		y:      marginTop,
		resume: resume,
	}
}

func Generate(resume *parser.Resume) ([]byte, error) {
	g := New(resume)

	g.pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: pageWidth, H: pageHeight},
	})

	if err := g.pdf.AddTTFFont(fontRegular, "internal/pdf/fonts/EBGaramond-Regular.ttf"); err != nil {
		return nil, fmt.Errorf("failed to load regular font: %w", err)
	}
	if err := g.pdf.AddTTFFont(fontBold, "internal/pdf/fonts/EBGaramond-Bold.ttf"); err != nil {
		return nil, fmt.Errorf("failed to load bold font: %w", err)
	}

	g.pdf.AddPage()

	g.drawHeader()
	g.drawLine()
	g.drawSummary()
	g.drawLine()
	g.drawExperience()
	g.drawLine()
	g.drawSkills()
	g.drawLine()
	g.drawCerts()
	g.drawLine()
	g.drawEducation()

	var buf bytes.Buffer
	if _, err := g.pdf.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("failed to write pdf: %w", err)
	}
	return buf.Bytes(), nil
}

// --- helpers ---

func (g *Generator) setFont(name string, size int) {
	g.pdf.SetFont(name, "", size)
}

func (g *Generator) drawLine() {
	g.y += 4
	g.pdf.SetLineWidth(0.5)
	g.pdf.Line(marginLeft, g.y, pageWidth-marginRight, g.y)
	g.y += 6
}

func (g *Generator) checkPageBreak(needed float64) {
	if g.y+needed > pageHeight-marginBottom {
		g.pdf.AddPage()
		g.y = marginTop
	}
}

func (g *Generator) drawCenteredText(text string, fontSize int, fontName string) {
	g.setFont(fontName, fontSize)
	w, _ := g.pdf.MeasureTextWidth(text)
	x := (pageWidth - w) / 2
	g.pdf.SetXY(x, g.y)
	g.pdf.Cell(&gopdf.Rect{W: w + 1, H: lineHeight}, text)
	g.y += lineHeight
}

func (g *Generator) drawTextAt(x float64, text string) {
	g.pdf.SetXY(x, g.y)
	g.pdf.Cell(&gopdf.Rect{W: contentWidth, H: lineHeight}, text)
}

func (g *Generator) cleanBullet(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "•")
	text = strings.TrimPrefix(text, "-")
	text = strings.TrimSpace(text)
	return text
}

func (g *Generator) wrapText(text string, width float64) []string {
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, " .", ".")
	text = strings.TrimSpace(text)

	words := strings.Split(text, " ")
	var lines []string
	current := ""

	for _, word := range words {
		if word == "" {
			continue
		}
		test := current
		if test != "" {
			test += " "
		}
		test += word

		w, _ := g.pdf.MeasureTextWidth(test)
		if w > width && current != "" {
			lines = append(lines, current)
			current = word
		} else {
			current = test
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func (g *Generator) drawWrappedText(text string, x, width float64) {
	lines := g.wrapText(text, width)
	for _, line := range lines {
		g.checkPageBreak(lineHeight)
		g.pdf.SetXY(x, g.y)
		g.pdf.Cell(&gopdf.Rect{W: width, H: lineHeight}, line)
		g.y += lineHeight
	}
}

func (g *Generator) drawSectionHeader(title string) {
	g.checkPageBreak(lineHeight + sectionSpacing)
	g.y += 2
	g.setFont(fontBold, fontSizeSection)
	w, _ := g.pdf.MeasureTextWidth(title)
	x := (pageWidth - w) / 2
	g.pdf.SetXY(x, g.y)
	g.pdf.Cell(&gopdf.Rect{W: w + 1, H: lineHeight}, title)
	g.y += lineHeight
	g.y += 3
}

// --- sections ---

func (g *Generator) drawHeader() {
	g.drawCenteredText(g.resume.Name, fontSizeName, fontBold)
	g.y += 8

	g.setFont(fontRegular, fontSizeContact)
	w, _ := g.pdf.MeasureTextWidth(g.resume.Contact)
	x := (pageWidth - w) / 2
	g.pdf.SetXY(x, g.y)
	g.pdf.Cell(&gopdf.Rect{W: w + 1, H: lineHeight}, g.resume.Contact)
	g.y += lineHeight + 2
}

func (g *Generator) drawSummary() {
	g.drawSectionHeader("SUMMARY")
	g.setFont(fontRegular, fontSizeBody)
	g.drawWrappedText(g.resume.Summary, marginLeft, contentWidth)
}

func (g *Generator) drawExperience() {
	g.drawSectionHeader("EXPERIENCE")

	for _, job := range g.resume.Experience {
		g.checkPageBreak(lineHeight * 3)

		g.setFont(fontBold, fontSizeBody)
		jobHeader := job.Title + "  |  " + job.Company + "  |  " + job.Dates
		g.drawTextAt(marginLeft, jobHeader)
		g.y += lineHeight

		if job.Intro != "" {
			g.setFont(fontRegular, fontSizeBody)
			g.drawWrappedText(job.Intro, marginLeft, contentWidth)
		}

		for _, section := range job.Sections {
			if section.Header != "" {
				g.checkPageBreak(lineHeight + 4)
				g.y += 2
				g.setFont(fontBold, fontSizeBody)
				g.drawTextAt(marginLeft+10, section.Header)
				g.y += lineHeight
			}

			g.setFont(fontRegular, fontSizeBody)
			for _, bullet := range section.Bullets {
				bulletIndent := marginLeft + 18.0
				bulletWidth := contentWidth - 18.0

				cleaned := g.cleanBullet(bullet)
				lines := g.wrapText(cleaned, bulletWidth)

				for i, line := range lines {
					g.checkPageBreak(lineHeight)
					if i == 0 {
						g.pdf.SetXY(marginLeft+7, g.y)
						g.pdf.Cell(&gopdf.Rect{W: 11, H: lineHeight}, "•")
						g.pdf.SetXY(bulletIndent, g.y)
					} else {
						g.pdf.SetXY(bulletIndent, g.y)
					}
					g.pdf.Cell(&gopdf.Rect{W: bulletWidth, H: lineHeight}, line)
					g.y += lineHeight
				}
			}
		}
		g.y += 6
	}
}

func (g *Generator) drawSkills() {
	g.drawSectionHeader("TECHNICAL SKILLS")

	g.setFont(fontBold, fontSizeBody)
	for _, skill := range g.resume.Skills {
		g.checkPageBreak(lineHeight * 2)

		g.setFont(fontBold, fontSizeBody)
		labelText := skill.Label + ": "
		labelWidth, _ := g.pdf.MeasureTextWidth(labelText)
		g.pdf.SetXY(marginLeft, g.y)
		g.pdf.Cell(&gopdf.Rect{W: labelWidth, H: lineHeight}, labelText)

		g.setFont(fontRegular, fontSizeBody)
		remainingWidth := contentWidth - labelWidth
		valueLines := g.wrapText(skill.Value, remainingWidth)

		for i, line := range valueLines {
			if i == 0 {
				g.pdf.SetXY(marginLeft+labelWidth, g.y)
				g.pdf.Cell(&gopdf.Rect{W: remainingWidth, H: lineHeight}, line)
				g.y += lineHeight
			} else {
				g.checkPageBreak(lineHeight)
				g.pdf.SetXY(marginLeft, g.y) 
				g.pdf.Cell(&gopdf.Rect{W: contentWidth, H: lineHeight}, line) 
				g.y += lineHeight
			}
		}		

	}
}

func (g *Generator) drawCerts() {
	g.drawSectionHeader("CERTIFICATIONS")
	g.setFont(fontRegular, fontSizeBody)
	g.drawWrappedText(g.resume.Certs, marginLeft, contentWidth)
}

func (g *Generator) drawEducation() {
	g.drawSectionHeader("EDUCATION")
	g.setFont(fontRegular, fontSizeBody)
	for _, edu := range g.resume.Education {
		g.checkPageBreak(lineHeight)
		g.drawTextAt(marginLeft, edu)
		g.y += lineHeight
	}
}
