package vydra

type Solution struct {
	Tag    string `xml:"tag,attr"`
	Source struct {
		Path string `xml:"path,attr"`
		Type string `xml:"type,attr"`
	} `xml:"source"`
}

type Assets struct {
	Solutions struct {
		Solutions []Solution `xml:"solution"`
	} `xml:"solutions"`
}

type ProblemXML struct {
	Revision  int    `xml:"revision,attr"`
	ShortName string `xml:"short-name,attr"`
	Assets    Assets `xml:"assets"`
}
