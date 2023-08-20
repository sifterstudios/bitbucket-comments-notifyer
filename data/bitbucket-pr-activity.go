package data

type Comment struct {
	Properties          CommentProperties `json:"properties"`
	ID                  int               `json:"id"`
	Version             int               `json:"version"`
	Text                string            `json:"text"`
	Author              User              `json:"author"`
	CreatedDate         int64             `json:"createdDate"`
	UpdatedDate         int64             `json:"updatedDate"`
	Comments            []Comment         `json:"comments"`
	Tasks               []Task            `json:"tasks"`
	Severity            string            `json:"severity"`
	State               string            `json:"state"`
	PermittedOperations struct {
		Editable       bool `json:"editable"`
		Transitionable bool `json:"transitionable"`
		Deletable      bool `json:"deletable"`
	} `json:"permittedOperations"`
	ResolvedDate int64 `json:"resolvedDate"`
	Resolver     User  `json:"resolver"`
}

type CommentProperties struct {
	RepositoryID int `json:"repositoryId"`
}

type Task struct {
	Properties          Properties `json:"properties"`
	ID                  int        `json:"id"`
	Version             int        `json:"version"`
	Text                string     `json:"text"`
	Author              User       `json:"author"`
	CreatedDate         int64      `json:"createdDate"`
	UpdatedDate         int64      `json:"updatedDate"`
	Comments            []Comment  `json:"comments"`
	Tasks               []Task     `json:"tasks"`
	Severity            string     `json:"severity"`
	State               string     `json:"state"`
	PermittedOperations struct {
		Editable       bool `json:"editable"`
		Transitionable bool `json:"transitionable"`
		Deletable      bool `json:"deletable"`
	} `json:"permittedOperations"`
}

type Diff struct {
	Source      interface{} `json:"source"`
	Destination struct {
		Components []string `json:"components"`
		Parent     string   `json:"parent"`
		Name       string   `json:"name"`
		Extension  string   `json:"extension"`
		ToString   string   `json:"toString"`
	} `json:"destination"`
	Hunks []struct {
		Context         string `json:"context"`
		SourceLine      int    `json:"sourceLine"`
		SourceSpan      int    `json:"sourceSpan"`
		DestinationLine int    `json:"destinationLine"`
		DestinationSpan int    `json:"destinationSpan"`
		Segments        []struct {
			Type  string `json:"type"`
			Lines []struct {
				Destination int    `json:"destination"`
				Source      int    `json:"source"`
				Line        string `json:"line"`
				Truncated   bool   `json:"truncated"`
			} `json:"lines"`
			Truncated bool `json:"truncated"`
		} `json:"segments"`
		Truncated bool `json:"truncated"`
	} `json:"hunks"`
	Truncated  bool `json:"truncated"`
	Properties struct {
		ToHash   string `json:"toHash"`
		Current  bool   `json:"current"`
		FromHash string `json:"fromHash"`
	} `json:"properties"`
}

type CommentAnchor struct {
	FromHash string `json:"fromHash"`
	ToHash   string `json:"toHash"`
	Line     int    `json:"line"`
	LineType string `json:"lineType"`
	FileType string `json:"fileType"`
	Path     string `json:"path"`
	DiffType string `json:"diffType"`
	Orphaned bool   `json:"orphaned"`
}

type Value struct {
	ID            int           `json:"id"`
	CreatedDate   int64         `json:"createdDate"`
	User          User          `json:"user"`
	Action        string        `json:"action"`
	CommentAction string        `json:"commentAction"`
	Comment       Comment       `json:"comment"`
	CommentAnchor CommentAnchor `json:"commentAnchor"`
	Diff          Diff          `json:"diff"`
}
