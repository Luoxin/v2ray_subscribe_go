package country

type CountryInfo struct {
	Code    string `json:"code"`
	EnName  string `json:"en_name"`
	CnName  string `json:"cn_name"`
	Unicode string `json:"unicode"`
	Emoji   string `json:"emoji"`
}

type country struct {
	codeMap, enNameMap, cnNameMap map[string]*CountryInfo
}

var Country = NewCountry()

func NewCountry() *country {
	c := country{
		codeMap:   map[string]*CountryInfo{},
		enNameMap: map[string]*CountryInfo{},
		cnNameMap: map[string]*CountryInfo{},
	}

	countryList.Each(func(info *CountryInfo) {
		c.codeMap[info.Code] = info
		c.enNameMap[info.EnName] = info
		c.cnNameMap[info.CnName] = info
	})

	return &c
}

func (c *country) GetByEnName(name string) *CountryInfo {
	return c.enNameMap[name]
}

func (c *country) GetByCode(code string) *CountryInfo {
	return c.codeMap[code]
}

func (c *country) GetByCnName(name string) *CountryInfo {
	return c.cnNameMap[name]
}
