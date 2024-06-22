package labels

import (
	"bytes"
	"image"
	"rsherbs/pkg/assets"

	"github.com/go-pdf/fpdf"

	_ "embed"

	"github.com/disintegration/imaging"
)

type PDFImage struct {
	Name  string
	Image image.Image
}

type ItemFunc func(pdf *fpdf.Fpdf, x, y, width, height float64)
type InitFunc func(pdf *fpdf.Fpdf)

type GeneratePDFOptions struct {
	Print    bool
	Margin   float64
	Gap      float64
	Quantity int
	Width    float64
	Height   float64
}

type GeneratePDFWithItemsOptions struct {
	GeneratePDFOptions
	ItemFunc ItemFunc
	InitFunc InitFunc
	Images   []PDFImage
}

func GetCenteredCoordinates(
	x1, y1, x2, y2, width, height float64,
) (float64, float64) {
	return x1 + (x2-x1-width)/2, y1 + (y2-y1-height)/2
}

func GeneratePDFWithItems(
	options GeneratePDFWithItemsOptions,
) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	width, height := pdf.GetPageSize()

	pdf.AddUTF8FontFromBytes("DM Serif Display", "", assets.DMSerifDisplay)
	pdf.AddUTF8FontFromBytes("Domine", "", assets.Domine)

	if options.InitFunc != nil {
		options.InitFunc(pdf)
	}

	if options.Print {
		pdf.SetJavascript("this.print();")
	}

	for _, image := range options.Images {
		var encoded bytes.Buffer
		err := imaging.Encode(&encoded, image.Image, imaging.PNG)
		if err != nil {
			return nil, err
		}

		pdf.RegisterImageOptionsReader(image.Name, fpdf.ImageOptions{
			ImageType: "png",
		}, &encoded)
	}

	pdf.AddPage()

	x, y := options.Margin, options.Margin
	itemWidth := options.Width
	itemHeight := options.Height

	for i := 0; i < options.Quantity; i++ {
		if x+itemWidth+options.Margin > width {
			x = options.Margin
			y += itemHeight + options.Gap
		}
		if y+itemHeight+options.Margin > height {
			pdf.AddPage()
			y = options.Margin
			x = options.Margin
		}
		options.ItemFunc(pdf, x, y, itemWidth, itemHeight)

		x += itemWidth + options.Gap
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type LabelMetadata struct {
	Name  string
	Latin string
}

func FitText(pdf *fpdf.Fpdf, font string, text string, width float64) float64 {
	size := 0.5
	pdf.SetFont(font, "", 1)
	lines := pdf.SplitText(text, width)
	for pdf.GetStringWidth(lines[0]) < width {
		size += 0.5
		pdf.SetFontSize(float64(size))
	}

	return float64(size)
}

func GeneratePDF(
	meta LabelMetadata,
	options GeneratePDFOptions,
) ([]byte, error) {
	bg, err := imaging.Decode(bytes.NewReader(assets.Background))
	if err != nil {
		return nil, err
	}

	images := []PDFImage{
		{
			Name:  "bg",
			Image: bg,
		},
	}

	return GeneratePDFWithItems(GeneratePDFWithItemsOptions{
		GeneratePDFOptions: options,
		Images:             images,
		ItemFunc: func(pdf *fpdf.Fpdf, x, y, width, height float64) {
			pdf.ImageOptions("bg", x, y, width, height, false, fpdf.ImageOptions{}, 0, "")

			headingSize := min(FitText(pdf, "DM Serif Display", meta.Latin, width-width/10), 28)
			pdf.SetFont("DM Serif Display", "", headingSize)
			_, headingHeight := pdf.GetFontSize()

			subheadingSize := min(FitText(pdf, "Domine", meta.Name, width-width/10), 0.7*headingSize, 0.7*28)
			pdf.SetFont("Domine", "", subheadingSize)
			_, subheadingHeight := pdf.GetFontSize()

			pdf.SetFont("DM Serif Display", "", headingSize)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(x, y+height/4-headingHeight/2-headingHeight/4)
			pdf.CellFormat(width, headingHeight, meta.Latin, "", 0, "C", false, 0, "")

			pdf.SetFont("Domine", "", subheadingSize)
			pdf.SetTextColor(120, 124, 130)
			pdf.SetXY(x, y+height/2-subheadingHeight)
			pdf.CellFormat(width, subheadingHeight, meta.Name, "", 0, "C", false, 0, "")
		},
	})
}

func GeneratePDFMultiple(
	meta []LabelMetadata,
	options GeneratePDFOptions,
) ([]byte, error) {
	bg, err := imaging.Decode(bytes.NewReader(assets.Background))
	if err != nil {
		return nil, err
	}

	images := []PDFImage{
		{
			Name:  "bg",
			Image: bg,
		},
	}

	options.Quantity = len(meta)

	index := 0

	return GeneratePDFWithItems(GeneratePDFWithItemsOptions{
		GeneratePDFOptions: options,
		Images:             images,
		ItemFunc: func(pdf *fpdf.Fpdf, x, y, width, height float64) {
			current := meta[index]
			index++

			pdf.ImageOptions("bg", x, y, width, height, false, fpdf.ImageOptions{}, 0, "")

			headingSize := min(FitText(pdf, "DM Serif Display", current.Latin, width-width/10), 28)
			pdf.SetFont("DM Serif Display", "", headingSize)
			_, headingHeight := pdf.GetFontSize()

			subheadingSize := min(FitText(pdf, "Domine", current.Name, width-width/10), 0.7*headingSize, 0.7*28)
			pdf.SetFont("Domine", "", subheadingSize)
			_, subheadingHeight := pdf.GetFontSize()

			pdf.SetFont("DM Serif Display", "", headingSize)
			pdf.SetTextColor(0, 0, 0)
			pdf.SetXY(x, y+height/4-headingHeight/2-headingHeight/4)
			pdf.CellFormat(width, headingHeight, current.Latin, "", 0, "C", false, 0, "")

			pdf.SetFont("Domine", "", subheadingSize)
			pdf.SetTextColor(120, 124, 130)
			pdf.SetXY(x, y+height/2-subheadingHeight)
			pdf.CellFormat(width, subheadingHeight, current.Name, "", 0, "C", false, 0, "")
		},
	})
}
