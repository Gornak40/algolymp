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
	Description string  `xml:"description,attr"`
	Method      string  `xml:"method,attr"`
	Sample      bool    `xml:"sample,attr"`
	Cmd         string  `xml:"cmd,attr"`
	Verdict     string  `xml:"verdict,attr"`
	Group       string  `xml:"group,attr"`
	Points      float32 `xml:"points,attr"`
}

type Group struct {
	FeedbackPolicy string  `xml:"feedback-policy,attr"`
	PointsPolicy   string  `xml:"points-policy,attr"`
	Name           string  `xml:"name,attr"`
	Points         float32 `xml:"points,attr"`
	Dependencies   struct {
		Dependencies []struct {
			Group string `xml:"group,attr"`
		} `xml:"dependency"`
	} `xml:"dependencies"`
}

type TestSet struct {
	Name              string `xml:"name,attr"`
	TimeLimit         int    `xml:"time-limit"`
	MemoryLimit       int    `xml:"memory-limit"`
	TestCount         int    `xml:"test-count"`
	InputPathPattern  string `xml:"input-path-pattern"`
	OutputPathPattern string `xml:"output-path-pattern"`
	AnswerPathPattern string `xml:"answer-path-pattern"`
	Tests             struct {
		Tests []Test `xml:"test"`
	} `xml:"tests"`
	Groups struct {
		Groups []Group `xml:"group"`
	} `xml:"groups"`
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

type Judging struct {
	CPUName    string    `xml:"cpu-name,attr"`
	CPUSpeed   int       `xml:"cpu-speed,attr"`
	InputFile  string    `xml:"input-file,attr"`
	OutputFile string    `xml:"output-file,attr"`
	TestSets   []TestSet `xml:"testset"`
}

type ProblemXML struct {
	Revision   int    `xml:"revision,attr"`
	ShortName  string `xml:"short-name,attr"`
	Assets     Assets `xml:"assets"`
	Files      Files  `xml:"files"`
	Statements struct {
		Statements []Statement `xml:"statement"`
	} `xml:"statements"`
	Judging Judging `xml:"judging"`
	Tags    struct {
		Tags []Tag `xml:"tag"`
	} `xml:"tags"`
}
