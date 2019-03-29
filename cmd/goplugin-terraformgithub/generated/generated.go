// Code generated by Lyra DO NOT EDIT.

// This code is generated on a per-Provider basis using "tf-gen"
// Long term our hope is to remove this generation step and adopt dynamic approach

package github

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/servicesdk/service"
	"github.com/lyraproj/terraform-bridge/pkg/bridge"
)

type (
	Branch_protection struct {
		Branch_protection_id          *string `lyra:"tf-gen.ignore"`
		Branch                        string
		Enforce_admins                bool `puppet:"value=>false"`
		Etag                          *string
		Repository                    string
		Required_pull_request_reviews *map[string]interface{}   `puppet:"type => Optional[Struct[Optional[dismiss_stale_reviews]=>Boolean,Optional[dismissal_teams]=>Array[String],Optional[dismissal_users]=>Array[String],Optional[require_code_owner_reviews]=>Boolean]]"`
		Required_status_checks        *map[string]interface{}   `puppet:"type => Optional[Struct[Optional[contexts]=>Array[String],Optional[strict]=>Boolean]]"`
		Restrictions                  *map[string][]interface{} `puppet:"type => Optional[Struct[Optional[teams]=>Array[String],Optional[users]=>Array[String]]]"`
	}

	Issue_label struct {
		Issue_label_id *string `lyra:"tf-gen.ignore"`
		Color          string
		Description    *string
		Etag           *string
		Name           string
		Repository     string
		Url            *string
	}

	Membership struct {
		Membership_id *string `lyra:"tf-gen.ignore"`
		Etag          *string
		Role          string `puppet:"value=>'member'"`
		Username      string
	}

	Organization_project struct {
		Organization_project_id *string `lyra:"tf-gen.ignore"`
		Body                    *string
		Etag                    *string
		Name                    string
		Url                     *string
	}

	Organization_webhook struct {
		Organization_webhook_id *string            `lyra:"tf-gen.ignore"`
		Active                  bool               `puppet:"value=>true"`
		Configuration           *map[string]string `puppet:"type => Optional[Struct[Optional[content_type]=>String,Optional[insecure_ssl]=>String,Optional[secret]=>String,url=>String]]"`
		Etag                    *string
		Events                  []string
		Name                    string
		Url                     *string
	}

	Project_column struct {
		Project_column_id *string `lyra:"tf-gen.ignore"`
		Etag              *string
		Name              string
		Project_id        string
	}

	Repository struct {
		Repository_id      *string `lyra:"tf-gen.ignore"`
		Allow_merge_commit bool    `puppet:"value=>true"`
		Allow_rebase_merge bool    `puppet:"value=>true"`
		Allow_squash_merge bool    `puppet:"value=>true"`
		Archived           bool    `puppet:"value=>false"`
		Auto_init          *bool
		Default_branch     *string
		Description        *string
		Etag               *string
		Full_name          *string
		Git_clone_url      *string
		Gitignore_template *string
		Has_downloads      *bool
		Has_issues         *bool
		Has_projects       *bool
		Has_wiki           *bool
		Homepage_url       *string
		Html_url           *string
		Http_clone_url     *string
		License_template   *string
		Name               string
		Private            *bool
		Ssh_clone_url      *string
		Svn_url            *string
		Topics             *[]string
	}

	Repository_collaborator struct {
		Repository_collaborator_id *string `lyra:"tf-gen.ignore"`
		Permission                 string  `puppet:"value=>'push'"`
		Repository                 string
		Username                   string
	}

	Repository_deploy_key struct {
		Repository_deploy_key_id *string `lyra:"tf-gen.ignore"`
		Etag                     *string
		Key                      string
		Read_only                bool `puppet:"value=>true"`
		Repository               string
		Title                    string
	}

	Repository_project struct {
		Repository_project_id *string `lyra:"tf-gen.ignore"`
		Body                  *string
		Etag                  *string
		Name                  string
		Repository            string
		Url                   *string
	}

	Repository_webhook struct {
		Repository_webhook_id *string            `lyra:"tf-gen.ignore"`
		Active                bool               `puppet:"value=>true"`
		Configuration         *map[string]string `puppet:"type => Optional[Struct[Optional[content_type]=>String,Optional[insecure_ssl]=>String,Optional[secret]=>String,url=>String]]"`
		Etag                  *string
		Events                []string
		Name                  string
		Repository            string
		Url                   *string
	}

	Team struct {
		Team_id        *string `lyra:"tf-gen.ignore"`
		Description    *string
		Etag           *string
		Ldap_dn        *string
		Name           string
		Parent_team_id *int64
		Privacy        string `puppet:"value=>'secret'"`
		Slug           *string
	}

	Team_membership struct {
		Team_membership_id *string `lyra:"tf-gen.ignore"`
		Etag               *string
		Role               string `puppet:"value=>'member'"`
		Team_id            string
		Username           string
	}

	Team_repository struct {
		Team_repository_id *string `lyra:"tf-gen.ignore"`
		Etag               *string
		Permission         string `puppet:"value=>'pull'"`
		Repository         string
		Team_id            string
	}

	User_gpg_key struct {
		User_gpg_key_id    *string `lyra:"tf-gen.ignore"`
		Armored_public_key string
		Etag               *string
		Key_id             *string
	}

	User_ssh_key struct {
		User_ssh_key_id *string `lyra:"tf-gen.ignore"`
		Etag            *string
		Key             string
		Title           string
		Url             *string
	}
)

func Initialize(sb *service.Builder, p *schema.Provider) {
	// Generic handler API
	sb.RegisterAPI("TerraformGitHub::GenericHandler", bridge.NewTFHandler(nil, nil, "", ""))

	// Registration of resource types with handler
	var rt px.Type
	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Branch_protection{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("branch_protection_id", "etag")
		b.ImmutableAttributes("branch", "repository")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Branch_protectionHandler", bridge.NewTFHandler(p, rt, "branch_protection_id", "github_branch_protection"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Issue_label{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("issue_label_id", "etag", "url")
		b.ImmutableAttributes("repository")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Issue_labelHandler", bridge.NewTFHandler(p, rt, "issue_label_id", "github_issue_label"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Membership{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("membership_id", "etag")
		b.ImmutableAttributes("username")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::MembershipHandler", bridge.NewTFHandler(p, rt, "membership_id", "github_membership"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Organization_project{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("organization_project_id", "etag", "url")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Organization_projectHandler", bridge.NewTFHandler(p, rt, "organization_project_id", "github_organization_project"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Organization_webhook{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("organization_webhook_id", "etag", "url")
		b.ImmutableAttributes("name")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Organization_webhookHandler", bridge.NewTFHandler(p, rt, "organization_webhook_id", "github_organization_webhook"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Project_column{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("project_column_id", "etag")
		b.ImmutableAttributes("project_id")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Project_columnHandler", bridge.NewTFHandler(p, rt, "project_column_id", "github_project_column"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Repository{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("repository_id", "default_branch", "etag", "full_name", "git_clone_url", "html_url", "http_clone_url", "ssh_clone_url", "svn_url")
		b.ImmutableAttributes("auto_init", "gitignore_template", "license_template", "name")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::RepositoryHandler", bridge.NewTFHandler(p, rt, "repository_id", "github_repository"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Repository_collaborator{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("repository_collaborator_id")
		b.ImmutableAttributes("permission", "repository", "username")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Repository_collaboratorHandler", bridge.NewTFHandler(p, rt, "repository_collaborator_id", "github_repository_collaborator"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Repository_deploy_key{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("repository_deploy_key_id", "etag")
		b.ImmutableAttributes("key", "read_only", "repository", "title")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Repository_deploy_keyHandler", bridge.NewTFHandler(p, rt, "repository_deploy_key_id", "github_repository_deploy_key"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Repository_project{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("repository_project_id", "etag", "url")
		b.ImmutableAttributes("repository")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Repository_projectHandler", bridge.NewTFHandler(p, rt, "repository_project_id", "github_repository_project"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Repository_webhook{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("repository_webhook_id", "etag", "url")
		b.ImmutableAttributes("name", "repository")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Repository_webhookHandler", bridge.NewTFHandler(p, rt, "repository_webhook_id", "github_repository_webhook"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Team{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("team_id", "etag", "slug")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::TeamHandler", bridge.NewTFHandler(p, rt, "team_id", "github_team"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Team_membership{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("team_membership_id", "etag")
		b.ImmutableAttributes("team_id", "username")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Team_membershipHandler", bridge.NewTFHandler(p, rt, "team_membership_id", "github_team_membership"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&Team_repository{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("team_repository_id", "etag")
		b.ImmutableAttributes("repository", "team_id")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::Team_repositoryHandler", bridge.NewTFHandler(p, rt, "team_repository_id", "github_team_repository"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&User_gpg_key{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("user_gpg_key_id", "etag", "key_id")
		b.ImmutableAttributes("armored_public_key")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::User_gpg_keyHandler", bridge.NewTFHandler(p, rt, "user_gpg_key_id", "github_user_gpg_key"), rt)

	rt = sb.RegisterTypes("TerraformGitHub", sb.BuildResource(&User_ssh_key{}, func(b service.ResourceTypeBuilder) {
		b.ProvidedAttributes("user_ssh_key_id", "etag", "url")
		b.ImmutableAttributes("key", "title")
	}))[0]
	sb.RegisterHandler("TerraformGitHub::User_ssh_keyHandler", bridge.NewTFHandler(p, rt, "user_ssh_key_id", "github_user_ssh_key"), rt)

}
