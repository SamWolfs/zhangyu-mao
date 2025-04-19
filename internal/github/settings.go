package github

import (
	"bytes"
	"encoding/json"

	"github.com/google/go-github/v71/github"
	"sigs.k8s.io/yaml"
)

// These structures are based off of those defined in google/go-github with
// fields added - to account for path parameters - and
// fields removed - to remove the possibility of modifying them.

type Protection struct {
	Name string `json:"name"`
	github.ProtectionRequest
}

type Repository struct {
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	Homepage      *string `json:"homepage,omitempty"`
	DefaultBranch *string `json:"default_branch,omitempty"`

	HasIssues   *bool `json:"has_issues,omitempty"`
	HasProjects *bool `json:"has_projects,omitempty"`
	HasWiki     *bool `json:"has_wiki,omitempty"`

	AllowSquashMerge         *bool   `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit         *bool   `json:"allow_merge_commit,omitempty"`
	AllowRebaseMerge         *bool   `json:"allow_rebase_merge,omitempty"`
	AllowUpdateBranch        *bool   `json:"allow_update_branch,omitempty"`
	AllowAutoMerge           *bool   `json:"allow_auto_merge,omitempty"`
	AllowForking             *bool   `json:"allow_forking,omitempty"`
	DeleteBranchOnMerge      *bool   `json:"delete_branch_on_merge,omitempty"`
	SquashMergeCommitTitle   *string `json:"squash_merge_commit_title,omitempty"`
	SquashMergeCommitMessage *string `json:"squash_merge_commit_message,omitempty"`
	MergeCommitTitle         *string `json:"merge_commit_title,omitempty"`
	MergeCommitMessage       *string `json:"merge_commit_message,omitempty"`
}

type service struct {
	Client *github.Client
}

type Settings struct {
	Autolinks   *[]Autolink   `json:"autolinks"`
	Protections *[]Protection `json:"protections"`
	Repository  *Repository   `json:"repository"`
}

func (r *Repository) ToGitHub() *github.Repository {
	return &github.Repository{
		Name: r.Name, Description: r.Description, Homepage: r.Homepage,
		DefaultBranch: r.DefaultBranch, HasIssues: r.HasIssues, HasProjects: r.HasProjects,
		HasWiki: r.HasWiki, AllowSquashMerge: r.AllowSquashMerge, AllowMergeCommit: r.AllowMergeCommit,
		AllowRebaseMerge: r.AllowRebaseMerge, AllowUpdateBranch: r.AllowUpdateBranch,
		AllowAutoMerge: r.AllowAutoMerge, AllowForking: r.AllowForking, DeleteBranchOnMerge: r.DeleteBranchOnMerge,
		SquashMergeCommitTitle: r.SquashMergeCommitTitle, SquashMergeCommitMessage: r.SquashMergeCommitMessage,
		MergeCommitTitle: r.MergeCommitTitle, MergeCommitMessage: r.MergeCommitMessage,
	}
}

func DecodeSettings(yamlSettings []byte) (*Settings, error) {
	jsonSettings, err := yaml.YAMLToJSON(yamlSettings)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(jsonSettings)
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	settings := Settings{}
	err = decoder.Decode(&settings)
	return &settings, err
}
