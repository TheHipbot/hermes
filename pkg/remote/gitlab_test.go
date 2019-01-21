package remote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GitlabRemoteSuite struct {
	suite.Suite
}

func (s *GitlabRemoteSuite) TestGitlabCreator() {
	opts := &DriverOpts{}
	d, err := gitlabCreator(opts)
	s.Nil(err, "Should return without error")
	s.IsType(d, &Gitlab{}, "Driver should be Gitlab type")
	gl := d.(*Gitlab)
	s.Equal(defaultGitlabAPIHost, gl.Host, "Gitlab Driver should be created with default host")
}

func (s *GitlabRemoteSuite) TestGitlabSetHost() {
	d, err := gitlabCreator(&DriverOpts{})
	s.Nil(err, "Creator should not return error")
	gl := d.(*Gitlab)
	d.SetHost("http://test.gitlab.com")
	s.Equal("http://test.gitlab.com", gl.Host, "Host should be set")
}

func (s *GitlabRemoteSuite) TestGitlabSetAuth() {
	d, err := gitlabCreator(&DriverOpts{})
	s.Nil(err, "Creator should not return error")
	testAuth := Auth{
		Token: "1234abc",
	}
	gl := d.(*Gitlab)
	d.Authenticate(testAuth)
	s.Equal(testAuth, gl.Auth, "Auth should be set")
}

func (s *GitlabRemoteSuite) TestGitlabAuthType() {
	d, err := gitlabCreator(&DriverOpts{})
	s.Nil(err, "Creator should not return error")
	s.Equal(authToken, d.AuthType(), "AuthType should be authToken")
}

func (s *GitlabRemoteSuite) TestGetRepos() {
	reqNum := 0
	testToken := "1234abcd"
	testURL := ""
	linkHeaders := []string{
		`<%s/api/v4/projects?membership=true&per_page=20&private_token=%s&page=1>;rel="prev", <%s/api/v4/projects?membership=true&per_page=20&private_token=%s&page=2>;rel="next"`,
		`<%s/api/v4/projects?membership=true&per_page=20&private_token=%s&page=1>; rel="prev", <%s/api/v4/projects?membership=true&per_page=20&private_token=%s&page=1>; rel="first"`,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal(testToken, r.URL.Query()["private_token"][0], "The requests should have the correct access token")
		s.Equal("true", r.URL.Query()["membership"][0], "The requests should have the correct membership value")
		w.Header().Set("Link", fmt.Sprintf(linkHeaders[reqNum], testURL, testToken, testURL, testToken))
		fmt.Fprintln(w, gitlabTestResponses[reqNum])
		reqNum++
	}))
	defer ts.Close()
	testURL = ts.URL

	gl := &Gitlab{
		Auth: Auth{
			Token: testToken,
		},
		Host: ts.URL,
		Opts: &DriverOpts{},
	}
	res, err := gl.GetRepos()
	s.Nil(err, "No error should be returned")
	s.Equal(gitlabTestResult, res, "Results from GetRepos should match the mock responses")
}

func (s *GitlabRemoteSuite) TestGetReposWithMemberOnly() {
	reqNum := 0
	testToken := "1234abcd"
	testURL := ""
	linkHeaders := []string{
		`<%s/api/v4/projects?membership=false&per_page=20&private_token=%s&page=1>;rel="prev", <%s/api/v4/projects?membership=false&per_page=20&private_token=%s&page=2>;rel="next"`,
		`<%s/api/v4/projects?membership=false&per_page=20&private_token=%s&page=1>; rel="prev", <%s/api/v4/projects?membership=false&per_page=20&private_token=%s&page=1>; rel="first"`,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal(testToken, r.URL.Query()["private_token"][0], "The requests should have the correct access token")
		s.Equal("false", r.URL.Query()["membership"][0], "The requests should have the correct membership value")
		w.Header().Set("Link", fmt.Sprintf(linkHeaders[reqNum], testURL, testToken, testURL, testToken))
		fmt.Fprintln(w, gitlabTestResponses[reqNum])
		reqNum++
	}))
	defer ts.Close()
	testURL = ts.URL

	gl := &Gitlab{
		Auth: Auth{
			Token: testToken,
		},
		Host: ts.URL,
		Opts: &DriverOpts{
			AllRepos: true,
		},
	}
	res, err := gl.GetRepos()
	s.Nil(err, "No error should be returned")
	s.Equal(gitlabTestResult, res, "Results from GetRepos should match the mock responses")
}

func TestGitlabRemoteSuite(t *testing.T) {
	suite.Run(t, new(GitlabRemoteSuite))
}

var (
	gitlabTestResult = []map[string]string{
		{
			"url":       "https://gitlab.thehipbot.com/jchamber/base-images",
			"name":      "gitlab.thehipbot.com/jchamber/base-images",
			"ssh_url":   "ssh://git@gitlab.thehipbot.com:3389/jchamber/base-images.git",
			"clone_url": "https://gitlab.thehipbot.com/jchamber/base-images.git",
		},
		{
			"url":       "https://gitlab.thehipbot.com/jchamber/test-gitlab",
			"name":      "gitlab.thehipbot.com/jchamber/test-gitlab",
			"ssh_url":   "ssh://git@gitlab.thehipbot.com:3389/jchamber/test-gitlab.git",
			"clone_url": "https://gitlab.thehipbot.com/jchamber/test-gitlab.git",
		},
		{
			"url":       "https://gitlab.thehipbot.com/poit-deploy-scripts/api-test",
			"name":      "gitlab.thehipbot.com/poit-deploy-scripts/api-test",
			"ssh_url":   "ssh://git@gitlab.thehipbot.com:3389/poit-deploy-scripts/api-test.git",
			"clone_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/api-test.git",
		},
		{
			"url":       "https://gitlab.thehipbot.com/HDMCS/artisan_task_generator",
			"name":      "gitlab.thehipbot.com/HDMCS/artisan_task_generator",
			"ssh_url":   "ssh://git@gitlab.thehipbot.com:3389/HDMCS/artisan_task_generator.git",
			"clone_url": "https://gitlab.thehipbot.com/HDMCS/artisan_task_generator.git",
		},
		{
			"url":       "https://gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy",
			"name":      "gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy",
			"ssh_url":   "ssh://git@gitlab.thehipbot.com:3389/poit-deploy-scripts/datastoreproxy.git",
			"clone_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy.git",
		},
	}
	gitlabTestResponses = []string{
		`[
			{
				"id": 520,
				"description": "Docker Base Images",
				"name": "base-images",
				"name_with_namespace": "Chambers, Jeremy / base-images",
				"path": "base-images",
				"path_with_namespace": "jchamber/base-images",
				"created_at": "2018-11-29T11:01:22.044-06:00",
				"default_branch": "master",
				"tag_list": [],
				"ssh_url_to_repo": "ssh://git@gitlab.thehipbot.com:3389/jchamber/base-images.git",
				"http_url_to_repo": "https://gitlab.thehipbot.com/jchamber/base-images.git",
				"web_url": "https://gitlab.thehipbot.com/jchamber/base-images",
				"readme_url": "https://gitlab.thehipbot.com/jchamber/base-images/blob/master/README.md",
				"avatar_url": null,
				"star_count": 0,
				"forks_count": 0,
				"last_activity_at": "2018-11-29T16:17:29.399-06:00",
				"namespace": {
					"id": 528,
					"name": "jchamber",
					"path": "jchamber",
					"kind": "user",
					"full_path": "jchamber",
					"parent_id": null
				},
				"_links": {
					"self": "https://gitlab.thehipbot.com/api/v4/projects/520",
					"issues": "https://gitlab.thehipbot.com/api/v4/projects/520/issues",
					"merge_requests": "https://gitlab.thehipbot.com/api/v4/projects/520/merge_requests",
					"repo_branches": "https://gitlab.thehipbot.com/api/v4/projects/520/repository/branches",
					"labels": "https://gitlab.thehipbot.com/api/v4/projects/520/labels",
					"events": "https://gitlab.thehipbot.com/api/v4/projects/520/events",
					"members": "https://gitlab.thehipbot.com/api/v4/projects/520/members"
				},
				"archived": false,
				"visibility": "internal",
				"owner": {
					"id": 431,
					"name": "Chambers, Jeremy",
					"username": "jchamber",
					"state": "active",
					"avatar_url": "https://secure.gravatar.com/avatar/726e4c20a30fd1e77a7b587ee35b9722?s=80&d=identicon",
					"web_url": "https://gitlab.thehipbot.com/jchamber"
				},
				"resolve_outdated_diff_discussions": false,
				"container_registry_enabled": true,
				"issues_enabled": true,
				"merge_requests_enabled": true,
				"wiki_enabled": true,
				"jobs_enabled": true,
				"snippets_enabled": true,
				"shared_runners_enabled": true,
				"lfs_enabled": true,
				"creator_id": 431,
				"import_status": "none",
				"open_issues_count": 0,
				"public_jobs": true,
				"ci_config_path": null,
				"shared_with_groups": [],
				"only_allow_merge_if_pipeline_succeeds": false,
				"request_access_enabled": false,
				"only_allow_merge_if_all_discussions_are_resolved": false,
				"printing_merge_request_link_enabled": true,
				"merge_method": "merge",
				"permissions": {
					"project_access": {
					"access_level": 40,
					"notification_level": 3
					},
					"group_access": null
				},
				"approvals_before_merge": 0,
				"mirror": false,
				"external_authorization_classification_label": null
			},
			{
				"id": 515,
				"description": "",
				"name": "test-gitlab",
				"name_with_namespace": "Chambers, Jeremy / test-gitlab",
				"path": "test-gitlab",
				"path_with_namespace": "jchamber/test-gitlab",
				"created_at": "2018-11-28T16:05:20.396-06:00",
				"default_branch": null,
				"tag_list": [],
				"ssh_url_to_repo": "ssh://git@gitlab.thehipbot.com:3389/jchamber/test-gitlab.git",
				"http_url_to_repo": "https://gitlab.thehipbot.com/jchamber/test-gitlab.git",
				"web_url": "https://gitlab.thehipbot.com/jchamber/test-gitlab",
				"readme_url": null,
				"avatar_url": null,
				"star_count": 0,
				"forks_count": 0,
				"last_activity_at": "2018-11-28T16:05:20.396-06:00",
				"namespace": {
					"id": 528,
					"name": "jchamber",
					"path": "jchamber",
					"kind": "user",
					"full_path": "jchamber",
					"parent_id": null
				},
				"_links": {
					"self": "https://gitlab.thehipbot.com/api/v4/projects/515",
					"issues": "https://gitlab.thehipbot.com/api/v4/projects/515/issues",
					"merge_requests": "https://gitlab.thehipbot.com/api/v4/projects/515/merge_requests",
					"repo_branches": "https://gitlab.thehipbot.com/api/v4/projects/515/repository/branches",
					"labels": "https://gitlab.thehipbot.com/api/v4/projects/515/labels",
					"events": "https://gitlab.thehipbot.com/api/v4/projects/515/events",
					"members": "https://gitlab.thehipbot.com/api/v4/projects/515/members"
				},
				"archived": false,
				"visibility": "internal",
				"owner": {
					"id": 431,
					"name": "Chambers, Jeremy",
					"username": "jchamber",
					"state": "active",
					"avatar_url": "https://secure.gravatar.com/avatar/726e4c20a30fd1e77a7b587ee35b9722?s=80&d=identicon",
					"web_url": "https://gitlab.thehipbot.com/jchamber"
				},
				"resolve_outdated_diff_discussions": false,
				"container_registry_enabled": true,
				"issues_enabled": true,
				"merge_requests_enabled": true,
				"wiki_enabled": true,
				"jobs_enabled": true,
				"snippets_enabled": true,
				"shared_runners_enabled": true,
				"lfs_enabled": true,
				"creator_id": 431,
				"import_status": "none",
				"open_issues_count": 0,
				"public_jobs": true,
				"ci_config_path": null,
				"shared_with_groups": [],
				"only_allow_merge_if_pipeline_succeeds": false,
				"request_access_enabled": false,
				"only_allow_merge_if_all_discussions_are_resolved": false,
				"printing_merge_request_link_enabled": true,
				"merge_method": "merge",
				"permissions": {
					"project_access": {
					"access_level": 40,
					"notification_level": 3
					},
					"group_access": null
				},
				"approvals_before_merge": 0,
				"mirror": false,
				"external_authorization_classification_label": null
			}, {
				"id": 480,
				"description": "— Test project for integration with the GitLab API",
				"name": "api-test",
				"name_with_namespace": "HAD POIT Scripts / api-test",
				"path": "api-test",
				"path_with_namespace": "poit-deploy-scripts/api-test",
				"created_at": "2018-11-15T10:24:27.715-06:00",
				"default_branch": "master",
				"tag_list": [],
				"ssh_url_to_repo": "ssh://git@gitlab.thehipbot.com:3389/poit-deploy-scripts/api-test.git",
				"http_url_to_repo": "https://gitlab.thehipbot.com/poit-deploy-scripts/api-test.git",
				"web_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/api-test",
				"readme_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/api-test/blob/master/README.md",
				"avatar_url": null,
				"star_count": 0,
				"forks_count": 0,
				"last_activity_at": "2018-11-30T15:16:28.890-06:00",
				"namespace": {
				  "id": 184,
				  "name": "HAD POIT Scripts",
				  "path": "poit-deploy-scripts",
				  "kind": "group",
				  "full_path": "poit-deploy-scripts",
				  "parent_id": null
				},
				"_links": {
				  "self": "https://gitlab.thehipbot.com/api/v4/projects/480",
				  "issues": "https://gitlab.thehipbot.com/api/v4/projects/480/issues",
				  "merge_requests": "https://gitlab.thehipbot.com/api/v4/projects/480/merge_requests",
				  "repo_branches": "https://gitlab.thehipbot.com/api/v4/projects/480/repository/branches",
				  "labels": "https://gitlab.thehipbot.com/api/v4/projects/480/labels",
				  "events": "https://gitlab.thehipbot.com/api/v4/projects/480/events",
				  "members": "https://gitlab.thehipbot.com/api/v4/projects/480/members"
				},
				"archived": false,
				"visibility": "internal",
				"resolve_outdated_diff_discussions": false,
				"container_registry_enabled": true,
				"issues_enabled": true,
				"merge_requests_enabled": true,
				"wiki_enabled": true,
				"jobs_enabled": true,
				"snippets_enabled": true,
				"shared_runners_enabled": true,
				"lfs_enabled": true,
				"creator_id": 194,
				"import_status": "none",
				"open_issues_count": 0,
				"public_jobs": true,
				"ci_config_path": null,
				"shared_with_groups": [],
				"only_allow_merge_if_pipeline_succeeds": false,
				"request_access_enabled": false,
				"only_allow_merge_if_all_discussions_are_resolved": false,
				"printing_merge_request_link_enabled": true,
				"merge_method": "merge",
				"permissions": {
				  "project_access": null,
				  "group_access": {
					"access_level": 50,
					"notification_level": 3
				  }
				},
				"approvals_before_merge": 0,
				"mirror": false,
				"external_authorization_classification_label": null
			  }
		]`, `[
			  {
				"id": 414,
				"description": "",
				"name": "artisan task generator",
				"name_with_namespace": "HD Map Content Services / artisan task generator",
				"path": "artisan_task_generator",
				"path_with_namespace": "HDMCS/artisan_task_generator",
				"created_at": "2018-10-29T15:43:49.275-05:00",
				"default_branch": "master",
				"tag_list": [],
				"ssh_url_to_repo": "ssh://git@gitlab.thehipbot.com:3389/HDMCS/artisan_task_generator.git",
				"http_url_to_repo": "https://gitlab.thehipbot.com/HDMCS/artisan_task_generator.git",
				"web_url": "https://gitlab.thehipbot.com/HDMCS/artisan_task_generator",
				"readme_url": null,
				"avatar_url": null,
				"star_count": 0,
				"forks_count": 0,
				"last_activity_at": "2018-11-30T12:11:51.713-06:00",
				"namespace": {
				  "id": 183,
				  "name": "HD Map Content Services",
				  "path": "HDMCS",
				  "kind": "group",
				  "full_path": "HDMCS",
				  "parent_id": null
				},
				"_links": {
				  "self": "https://gitlab.thehipbot.com/api/v4/projects/414",
				  "issues": "https://gitlab.thehipbot.com/api/v4/projects/414/issues",
				  "merge_requests": "https://gitlab.thehipbot.com/api/v4/projects/414/merge_requests",
				  "repo_branches": "https://gitlab.thehipbot.com/api/v4/projects/414/repository/branches",
				  "labels": "https://gitlab.thehipbot.com/api/v4/projects/414/labels",
				  "events": "https://gitlab.thehipbot.com/api/v4/projects/414/events",
				  "members": "https://gitlab.thehipbot.com/api/v4/projects/414/members"
				},
				"archived": false,
				"visibility": "internal",
				"resolve_outdated_diff_discussions": false,
				"container_registry_enabled": true,
				"issues_enabled": true,
				"merge_requests_enabled": true,
				"wiki_enabled": true,
				"jobs_enabled": true,
				"snippets_enabled": true,
				"shared_runners_enabled": true,
				"lfs_enabled": true,
				"creator_id": 137,
				"import_status": "finished",
				"open_issues_count": 0,
				"public_jobs": true,
				"ci_config_path": "",
				"shared_with_groups": [],
				"only_allow_merge_if_pipeline_succeeds": false,
				"request_access_enabled": false,
				"only_allow_merge_if_all_discussions_are_resolved": false,
				"printing_merge_request_link_enabled": true,
				"merge_method": "merge",
				"permissions": {
				  "project_access": null,
				  "group_access": {
					"access_level": 50,
					"notification_level": 3
				  }
				},
				"approvals_before_merge": 0,
				"mirror": false,
				"external_authorization_classification_label": null
			},
			{
				"id": 362,
				"description": "",
				"name": "DataStoreProxy",
				"name_with_namespace": "HAD POIT Scripts / DataStoreProxy",
				"path": "datastoreproxy",
				"path_with_namespace": "poit-deploy-scripts/datastoreproxy",
				"created_at": "2018-10-11T15:55:07.954-05:00",
				"default_branch": "master",
				"tag_list": [],
				"ssh_url_to_repo": "ssh://git@gitlab.thehipbot.com:3389/poit-deploy-scripts/datastoreproxy.git",
				"http_url_to_repo": "https://gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy.git",
				"web_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy",
				"readme_url": "https://gitlab.thehipbot.com/poit-deploy-scripts/datastoreproxy/blob/master/README.md",
				"avatar_url": null,
				"star_count": 0,
				"forks_count": 0,
				"last_activity_at": "2018-12-03T09:59:42.388-06:00",
				"namespace": {
				  "id": 184,
				  "name": "HAD POIT Scripts",
				  "path": "poit-deploy-scripts",
				  "kind": "group",
				  "full_path": "poit-deploy-scripts",
				  "parent_id": null
				},
				"_links": {
				  "self": "https://gitlab.thehipbot.com/api/v4/projects/362",
				  "issues": "https://gitlab.thehipbot.com/api/v4/projects/362/issues",
				  "merge_requests": "https://gitlab.thehipbot.com/api/v4/projects/362/merge_requests",
				  "repo_branches": "https://gitlab.thehipbot.com/api/v4/projects/362/repository/branches",
				  "labels": "https://gitlab.thehipbot.com/api/v4/projects/362/labels",
				  "events": "https://gitlab.thehipbot.com/api/v4/projects/362/events",
				  "members": "https://gitlab.thehipbot.com/api/v4/projects/362/members"
				},
				"archived": false,
				"visibility": "internal",
				"resolve_outdated_diff_discussions": false,
				"container_registry_enabled": true,
				"issues_enabled": true,
				"merge_requests_enabled": true,
				"wiki_enabled": true,
				"jobs_enabled": true,
				"snippets_enabled": true,
				"shared_runners_enabled": true,
				"lfs_enabled": true,
				"creator_id": 139,
				"import_status": "none",
				"open_issues_count": 0,
				"public_jobs": true,
				"ci_config_path": null,
				"shared_with_groups": [],
				"only_allow_merge_if_pipeline_succeeds": false,
				"request_access_enabled": false,
				"only_allow_merge_if_all_discussions_are_resolved": false,
				"printing_merge_request_link_enabled": true,
				"merge_method": "merge",
				"permissions": {
				  "project_access": null,
				  "group_access": {
					"access_level": 50,
					"notification_level": 3
				  }
				},
				"approvals_before_merge": 0,
				"mirror": false,
				"external_authorization_classification_label": null
			  }
]`,
	}
)
