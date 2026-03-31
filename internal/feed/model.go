// Package feed builds Avito autoload feeds.
package feed

// Listing describes one Avito listing in internal application form.
type Listing struct {
	ID              string
	Title           string
	Description     string
	Category        string
	OperationType   string
	Address         string
	Price           int64
	ContactPhone    string
	ManagerName     string
	City            string
	Region          string
	Images          []string
	ConcreteGrade   string
	ConcreteClass   string
	Mobility        string
	FrostResistance string
	WaterResistance string
	MinVolumeM3     float64
	Delivery        bool
}
