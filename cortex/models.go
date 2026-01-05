package cortex

type Cortex struct {
	Openapi string     `yaml:"openapi"`
	Info    CortexInfo `yaml:"info"`
}

type CortexInfo struct {
	Tag   string `yaml:"x-cortex-tag"`
	Title string `yaml:"title"`

	Description    string                 `yaml:"description,omitempty"`
	Type           string                 `yaml:"x-cortex-type,omitempty"`
	Parents        []CortexTag            `yaml:"x-cortex-parents,omitempty"`
	Groups         []string               `yaml:"x-cortex-groups,omitempty"`
	Team           CortexTeam             `yaml:"x-cortex-team,omitempty"`
	Owners         []CortexOwner          `yaml:"x-cortex-owners,omitempty"`
	Slack          CortexSlack            `yaml:"x-cortex-slack,omitempty"`
	Link           []CortexLink           `yaml:"x-cortex-link,omitempty"`
	CustomMetadata map[string]interface{} `yaml:"x-cortex-custom-metadata,omitempty"`
	Git            CortexGit              `yaml:"x-cortex-git,omitempty"`
	Oncall         CortexOncall           `yaml:"x-cortex-oncall,omitempty"`
	Issues         CortexIssues           `yaml:"x-cortex-issues,omitempty"`
	Dependency     CortexDependency       `yaml:"x-cortex-dependency,omitempty"`
	SLOs           CortexSLOs             `yaml:"x-cortex-slos,omitempty"`
	StaticAnalysis CortexStaticAnalysis   `yaml:"x-cortex-static-analysis,omitempty"`
}

type CortexTag struct {
	Tag string `yaml:"tag"`
}

type CortexOwner struct {
	Type        string `yaml:"type"`
	Name        string `yaml:"name,omitempty"`
	Provider    string `yaml:"provider,omitempty"`
	Email       string `yaml:"email,omitempty"`
	Inheritance string `yaml:"inheritance,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type CortexTeam struct {
	Groups  []CortexTeamGroup  `yaml:"groups,omitempty"`
	Members []CortexTeamMember `yaml:"members,omitempty"`
}

type CortexTeamGroup struct {
	Name     string `yaml:"name"`
	Provider string `yaml:"provider"`
}

type CortexTeamMember struct {
	Name                 string `yaml:"name"`
	Email                string `yaml:"email"`
	NotificationsEnabled bool   `yaml:"notificationsEnabled"`
	Role                 string `yaml:"role,omitempty"`
}

type CortexSlack struct {
	Channels []CortexSlackChannel `yaml:"channels"`
}

type CortexSlackChannel struct {
	Name                 string `yaml:"name"`
	NotificationsEnabled bool   `yaml:"notificationsEnabled"`
	Description          string `yaml:"description,omitempty"`
}

type CortexLink struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}
type CortexGit struct {
	Github CortexGithub `yaml:"github"`
}

type CortexGithub struct {
	Repository string `yaml:"repository"`
	BasePath   string `yaml:"basepath,omitempty"`
	Alias      string `yaml:"alias,omitempty"`
}

type CortexOncall struct {
	VictorOps CortexOncallVictorOps `yaml:"victorops"`
}

type CortexOncallVictorOps struct {
	Type string `yaml:"type"`
	ID   string `yaml:"id"`
}

type CortexIssues struct {
	Jira CortexIssuesJira `yaml:"jira"`
}

type CortexIssuesJira struct {
	Projects []string `yaml:"projects"`
}

type CortexDependency struct {
	Cortex []CortexDependencyCortex `yaml:"cortex,omitempty"`
	AWS    CortexDependencyAWS      `yaml:"aws,omitempty"`
}

// UnmarshalYAML implements custom unmarshaling to gracefully handle malformed dependency data
func (cd *CortexDependency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create a temporary type to avoid infinite recursion
	type rawDependency struct {
		Cortex []CortexDependencyCortex `yaml:"cortex,omitempty"`
		AWS    CortexDependencyAWS      `yaml:"aws,omitempty"`
	}

	var raw rawDependency
	if err := unmarshal(&raw); err != nil {
		// Log warning but don't fail - set empty defaults
		// The error will be silently handled, allowing other fields to process
		cd.Cortex = []CortexDependencyCortex{}
		cd.AWS = CortexDependencyAWS{Tags: []Tag{}}
		return nil
	}

	// Successfully unmarshaled, populate the struct
	cd.Cortex = raw.Cortex
	cd.AWS = raw.AWS
	return nil
}

type CortexDependencyCortex struct {
	Tag         string `yaml:"tag"`
	Path        string `yaml:"path,omitempty"`
	Method      string `yaml:"method,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type CortexDependencyAWS struct {
	Tags []Tag `yaml:"tags"`
}

type CortexSLOs struct {
	NewRelic []CortexSLO `yaml:"newrelic"`
}

type CortexSLO struct {
	ID    string `yaml:"id"`
	Alias string `yaml:"alias,omitempty"`
}

type Tag struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type CortexStaticAnalysis struct {
	Sonarqube CortexStaticAnalysisSonarqube `yaml:"sonarqube"`
}

type CortexStaticAnalysisSonarqube struct {
	Project string `yaml:"project"`
	Alias   string `yaml:"alias,omitempty"`
}
