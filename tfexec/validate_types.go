package tfexec

// TODO: move these types to terraform-json

type Validation struct {
	Valid        bool `json:"valid"`
	ErrorCount   int  `json:"error_count"`
	WarningCount int  `json:"warning_count"`

	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Severity string `json:"severity,omitempty"`
	Summary  string `json:"summary,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Range    *Range `json:"range,omitempty"`
}

type Range struct {
	Filename string `json:"filename"`
	Start    Pos    `json:"start"`
	End      Pos    `json:"end"`
}

type Pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	Byte   int `json:"byte"`
}
