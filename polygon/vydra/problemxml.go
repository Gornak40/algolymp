package vydra

type File struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type Source struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type Executable struct {
	Source Source `xml:"source"`
}

type Solution struct {
	Tag    string `xml:"tag,attr"`
	Source Source `xml:"source"`
}

type Statement struct {
	Charset  string `xml:"charset,attr"`
	Language string `xml:"language,attr"`
	Path     string `xml:"path,attr"`
	Type     string `xml:"type,attr"`
}

type Test struct {
	Method  string `xml:"method,attr"`
	Sample  bool   `xml:"sample,attr"`
	Cmd     string `xml:"cmd,attr"`
	Verdict string `xml:"verdict,attr"`
}

type TestSet struct {
	TestCount         int    `xml:"test-count"`
	InputPathPattern  string `xml:"input-path-pattern"`
	OutputPathPattern string `xml:"output-path-pattern"`
	AnswerPathPattern string `xml:"answer-path-pattern"`
	Tests             []Test `xml:"tests"`
}

type Validator struct {
	Source  Source  `xml:"source"`
	TestSet TestSet `xml:"testset"`
}

type Checker struct {
	Name    string  `xml:"name,attr"`
	Type    string  `xml:"type,attr"`
	Source  Source  `xml:"source"`
	TestSet TestSet `xml:"testset"`
}

type Assets struct {
	Solutions struct {
		Solutions []Solution `xml:"solution"`
	} `xml:"solutions"`
	Validators struct {
		Validator *Validator `xml:"validator"`
	} `xml:"validators"`
	Checker *Checker `xml:"checker"`
}

type Files struct {
	Resources struct {
		Files []File `xml:"file"`
	} `xml:"resources"`
	Executables struct {
		Executables []Executable `xml:"executable"`
	} `xml:"executables"`
}

type Tag struct {
	Value string `xml:"value,attr"`
}

type ProblemXML struct {
	Revision   int    `xml:"revision,attr"`
	ShortName  string `xml:"short-name,attr"`
	Assets     Assets `xml:"assets"`
	Files      Files  `xml:"files"`
	Statements struct {
		Statements []Statement `xml:"statement"`
	} `xml:"statements"`
	Tags struct {
		Tags []Tag `xml:"tag"`
	} `xml:"tags"`
}
