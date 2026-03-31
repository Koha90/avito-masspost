package feed

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// File contains a top-level Avito XML feed document.
type File struct {
	XMLName xml.Name `xml:"Ads"`
	Format  string   `xml:"format,attr,omitempty"`
	Version string   `xml:"version,attr,omitempty"`
	Ads     []Ad     `xml:"Ad"`
}

// Ad contains one Avito XML ad entry.
type Ad struct {
	ID              string  `xml:"Id"`
	Title           string  `xml:"Title"`
	Description     string  `xml:"Description"`
	Category        string  `xml:"Category"`
	OperationType   string  `xml:"OperationType,omitempty"`
	Address         string  `xml:"Address,omitempty"`
	Price           int64   `xml:"Price"`
	ContactPhone    string  `xml:"ContactPhone,omitempty"`
	ManagerName     string  `xml:"ManagerName,omitempty"`
	Region          string  `xml:"Region,omitempty"`
	City            string  `xml:"City,omitempty"`
	ConcreteGrade   string  `xml:"ConcreteGrade,omitempty"`
	ConcreteClass   string  `xml:"ConcreteClass,omitempty"`
	Mobility        string  `xml:"Mobility,omitempty"`
	FrostResistance string  `xml:"FrostResistance,omitempty"`
	WaterResistance string  `xml:"WaterResistance,omitempty"`
	MinVolumeM3     float64 `xml:"MinVolumeM3,omitempty"`
	Delivery        string  `xml:"Delivery,omitempty"`
	Images          []Image `xml:"Images>Image,omitempty"`
}

// Image contains one image URL entry.
type Image struct {
	URL string `xml:"Url"`
}

// Write writes a feed file to disk.
func Write(path string, listings []Listing) error {
	file := File{
		Format:  "Avito",
		Version: "1",
		Ads:     make([]Ad, 0, len(listings)),
	}

	for _, listing := range listings {
		ad := Ad{
			ID:              listing.ID,
			Title:           listing.Title,
			Description:     listing.Description,
			Category:        listing.Category,
			OperationType:   listing.OperationType,
			Address:         listing.Address,
			Price:           listing.Price,
			ContactPhone:    listing.ContactPhone,
			ManagerName:     listing.ManagerName,
			Region:          listing.Region,
			City:            listing.City,
			ConcreteGrade:   listing.ConcreteGrade,
			ConcreteClass:   listing.ConcreteClass,
			Mobility:        listing.Mobility,
			FrostResistance: listing.FrostResistance,
			WaterResistance: listing.WaterResistance,
			MinVolumeM3:     listing.MinVolumeM3,
			Delivery:        yesNo(listing.Delivery),
			Images:          images(listing.Images),
		}

		file.Ads = append(file.Ads, ad)
	}

	data, err := xml.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal feed: %w", err)
	}

	data = append([]byte(xml.Header), data...)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create feed directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write feed file: %w", err)
	}

	return nil
}

func images(urls []string) []Image {
	if len(urls) == 0 {
		return nil
	}

	out := make([]Image, 0, len(urls))
	for _, u := range urls {
		out = append(out, Image{URL: u})
	}

	return out
}

func yesNo(v bool) string {
	if v {
		return "Да"
	}
	return "Нет"
}
