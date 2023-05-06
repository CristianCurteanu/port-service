package ports

import "go.mongodb.org/mongo-driver/bson"

type Port struct {
	PortCode    string    `bson:"port_code,omitempty"`
	Name        string    `bson:"name,omitempty"`
	City        string    `bson:"city,omitempty"`
	Country     string    `bson:"country,omitempty"`
	Code        string    `bson:"code,omitempty"`
	Alias       []string  `bson:"alias,omitempty"`
	Regions     []string  `bson:"regions,omitempty"`
	Coordinates []float64 `bson:"coordinates,omitempty"`
	Province    string    `bson:"province,omitempty"`
	Timezone    string    `bson:"timezone,omitempty"`
	Unlocs      []string  `bson:"unlocs,omitempty"`
}

func (p Port) AsBson() bson.M {
	return bson.M{
		"name":        p.Name,
		"city":        p.City,
		"country":     p.Country,
		"code":        p.Code,
		"alias":       p.Alias,
		"regions":     p.Regions,
		"coordinates": p.Coordinates,
		"province":    p.Province,
		"timezone":    p.Timezone,
		"unlocs":      p.Unlocs,
	}
}
