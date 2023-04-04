package radiofreq

type Spectrum int
type Frequency float64

const (
	S_5G Spectrum = iota
	S_LTE
)

func (s Spectrum) String() string {
	switch s {
	case S_5G:
		return "5G"
	case S_LTE:
		return "LTE"
	default:
		return "unknown"
	}
}

// Band defines a single band within a spectrum
type Band struct {
	Spec      Spectrum
	Shortname string
	Frequency float64
}

// Spectrum returns a given Band's spectrum name
func (b *Band) Spectrum() string {
	return b.Spec.String()
}

// Bands contains a list of bands the software may encounter
type Bands struct {
	bands []Band
}

// FrequencyFromShortname when given an appropriate shortname will return the
// associated freuency in GHz, else 0
func (b *Bands) FrequencyFromShortname(shortname string) float64 {
	for _, band := range b.bands {
		if band.Shortname == shortname {
			return band.Frequency
		}
	}
	return 0
}

// BandFromShortname when given an appropriate shortname will return the band
// object associated with it, else nil
func (b *Bands) BandFromShortname(shortname string) *Band {
	for _, band := range b.bands {
		if band.Shortname == shortname {
			return &band
		}
	}
	return nil
}

// Map returns a map[string]float64 of all bands
func (b *Bands) Map() map[string]float64 {
	ret := make(map[string]float64, len(b.bands))
	for _, band := range b.bands {
		ret[band.Shortname] = band.Frequency
	}
	return ret
}

// BandMap is a mapped list of band names and their associated
// Spectrum and Frequency in GHz
var BandMap = &Bands{
	bands: []Band{
		{
			Spec:      S_5G,
			Shortname: "n71",
			Frequency: 0.6,
		},
		{
			Spec:      S_5G,
			Shortname: "n41",
			Frequency: 2.5,
		},
		{
			Spec:      S_5G,
			Shortname: "n2",
			Frequency: 3.4,
		},
		{
			Spec:      S_5G,
			Shortname: "n77",
			Frequency: 3.7,
		},
		{
			Spec:      S_5G,
			Shortname: "n258",
			Frequency: 24,
		},
		{
			Spec:      S_5G,
			Shortname: "n261",
			Frequency: 39,
		},
		{
			Spec:      S_5G,
			Shortname: "n262",
			Frequency: 47,
		},
		{
			Spec:      S_LTE,
			Shortname: "B71",
			Frequency: 0.6,
		},
		{
			Spec:      S_LTE,
			Shortname: "B12",
			Frequency: 0.7,
		},
		{
			Spec:      S_LTE,
			Shortname: "B5",
			Frequency: 0.85,
		},
		{
			Spec:      S_LTE,
			Shortname: "B4",
			Frequency: 1.7,
		},
		{
			Spec:      S_LTE,
			Shortname: "B66",
			Frequency: 2.1,
		},
		{
			Spec:      S_LTE,
			Shortname: "B2",
			Frequency: 1.9,
		},
	},
}
