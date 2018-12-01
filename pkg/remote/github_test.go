package remote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GitHubRemoteSuite struct {
	suite.Suite
}

func (s *GitHubRemoteSuite) TestGitHubCreator() {
	d, err := githubCreator()
	s.Nil(err, "Should return without error")
	s.IsType(d, &GitHub{}, "Driver should be GitHub type")
	gh := d.(*GitHub)
	s.Equal(defaultGitHubAPIHost, gh.Host, "GitHub Driver should be created with default host")
}

func (s *GitHubRemoteSuite) TestSetHost() {
	d, err := githubCreator()
	s.Nil(err, "Creator should not return error")
	gh := d.(*GitHub)
	d.SetHost("http://test.github.com")
	s.Equal("http://test.github.com", gh.Host, "Host should be set")
}

func (s *GitHubRemoteSuite) TestSetAuth() {
	d, err := githubCreator()
	s.Nil(err, "Creator should not return error")
	testAuth := Auth{
		Token: "1234abc",
	}
	gh := d.(*GitHub)
	d.Authenticate(testAuth)
	s.Equal(testAuth, gh.Auth, "Auth should be set")
}

func (s *GitHubRemoteSuite) TestAuthType() {
	d, err := githubCreator()
	s.Nil(err, "Creator should not return error")
	s.Equal(authToken, d.AuthType(), "AuthType should be authToken")
}

func (s *GitHubRemoteSuite) TestGetRepos() {
	reqNum := 0
	testToken := "1234abcd"
	testURL := ""
	linkHeaders := []string{
		`<%s/user/repos?access_token=%s&page=2>; rel="next", <%s/user/repos?access_token=%s&page=2>; rel="last"`,
		`<%s/user/repos?access_token=%s&page=1>; rel="prev", <%s/user/repos?access_token=%s&page=1>; rel="first"`,
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal(testToken, r.URL.Query()["access_token"][0], "The requests should have the correct access token")
		w.Header().Set("Link", fmt.Sprintf(linkHeaders[reqNum], testURL, testToken, testURL, testToken))
		fmt.Fprintln(w, testResponses[reqNum])
		reqNum++
	}))
	defer ts.Close()
	testURL = ts.URL

	gh := &GitHub{
		Auth: Auth{
			Token: testToken,
		},
		Host: ts.URL,
	}
	res, err := gh.GetRepos()
	s.Nil(err, "No error should be returned")
	s.Equal(testResult, res, "Results from GetRepos should match the mock responses")
}

func TestGitHubRemoteSuite(t *testing.T) {
	suite.Run(t, new(GitHubRemoteSuite))
}

var (
	testResult = []map[string]string{
		{
			"url":  "https://github.com/carsdotcom/beacon",
			"name": "github.com/carsdotcom/beacon",
		},
		{
			"url":  "https://github.com/carsdotcom/bitcar",
			"name": "github.com/carsdotcom/bitcar",
		},
		{
			"url":  "https://github.com/carsdotcom/cars.com-sellandtrade-application",
			"name": "github.com/carsdotcom/cars.com-sellandtrade-application",
		},
		{
			"url":  "https://github.com/TheHipbot/blog-resources",
			"name": "github.com/TheHipbot/blog-resources",
		},
		{
			"url":  "https://github.com/TheHipbot/causal-relations",
			"name": "github.com/TheHipbot/causal-relations",
		},
	}
	testResponses = []string{
		`[
		{
		  "id": 72761897,
		  "node_id": "MDEwOlJlcG9zaXRvcnk3Mjc2MTg5Nw==",
		  "name": "beacon",
		  "full_name": "carsdotcom/beacon",
		  "owner": {
			"login": "carsdotcom",
			"id": 8259181,
			"node_id": "MDEyOk9yZ2FuaXphdGlvbjgyNTkxODE=",
			"avatar_url": "https://avatars2.githubusercontent.com/u/8259181?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/carsdotcom",
			"html_url": "https://github.com/carsdotcom",
			"followers_url": "https://api.github.com/users/carsdotcom/followers",
			"following_url": "https://api.github.com/users/carsdotcom/following{/other_user}",
			"gists_url": "https://api.github.com/users/carsdotcom/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/carsdotcom/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/carsdotcom/subscriptions",
			"organizations_url": "https://api.github.com/users/carsdotcom/orgs",
			"repos_url": "https://api.github.com/users/carsdotcom/repos",
			"events_url": "https://api.github.com/users/carsdotcom/events{/privacy}",
			"received_events_url": "https://api.github.com/users/carsdotcom/received_events",
			"type": "Organization",
			"site_admin": false
		  },
		  "private": false,
		  "html_url": "https://github.com/carsdotcom/beacon",
		  "description": null,
		  "fork": false,
		  "url": "https://api.github.com/repos/carsdotcom/beacon",
		  "forks_url": "https://api.github.com/repos/carsdotcom/beacon/forks",
		  "keys_url": "https://api.github.com/repos/carsdotcom/beacon/keys{/key_id}",
		  "collaborators_url": "https://api.github.com/repos/carsdotcom/beacon/collaborators{/collaborator}",
		  "teams_url": "https://api.github.com/repos/carsdotcom/beacon/teams",
		  "hooks_url": "https://api.github.com/repos/carsdotcom/beacon/hooks",
		  "issue_events_url": "https://api.github.com/repos/carsdotcom/beacon/issues/events{/number}",
		  "events_url": "https://api.github.com/repos/carsdotcom/beacon/events",
		  "assignees_url": "https://api.github.com/repos/carsdotcom/beacon/assignees{/user}",
		  "branches_url": "https://api.github.com/repos/carsdotcom/beacon/branches{/branch}",
		  "tags_url": "https://api.github.com/repos/carsdotcom/beacon/tags",
		  "blobs_url": "https://api.github.com/repos/carsdotcom/beacon/git/blobs{/sha}",
		  "git_tags_url": "https://api.github.com/repos/carsdotcom/beacon/git/tags{/sha}",
		  "git_refs_url": "https://api.github.com/repos/carsdotcom/beacon/git/refs{/sha}",
		  "trees_url": "https://api.github.com/repos/carsdotcom/beacon/git/trees{/sha}",
		  "statuses_url": "https://api.github.com/repos/carsdotcom/beacon/statuses/{sha}",
		  "languages_url": "https://api.github.com/repos/carsdotcom/beacon/languages",
		  "stargazers_url": "https://api.github.com/repos/carsdotcom/beacon/stargazers",
		  "contributors_url": "https://api.github.com/repos/carsdotcom/beacon/contributors",
		  "subscribers_url": "https://api.github.com/repos/carsdotcom/beacon/subscribers",
		  "subscription_url": "https://api.github.com/repos/carsdotcom/beacon/subscription",
		  "commits_url": "https://api.github.com/repos/carsdotcom/beacon/commits{/sha}",
		  "git_commits_url": "https://api.github.com/repos/carsdotcom/beacon/git/commits{/sha}",
		  "comments_url": "https://api.github.com/repos/carsdotcom/beacon/comments{/number}",
		  "issue_comment_url": "https://api.github.com/repos/carsdotcom/beacon/issues/comments{/number}",
		  "contents_url": "https://api.github.com/repos/carsdotcom/beacon/contents/{+path}",
		  "compare_url": "https://api.github.com/repos/carsdotcom/beacon/compare/{base}...{head}",
		  "merges_url": "https://api.github.com/repos/carsdotcom/beacon/merges",
		  "archive_url": "https://api.github.com/repos/carsdotcom/beacon/{archive_format}{/ref}",
		  "downloads_url": "https://api.github.com/repos/carsdotcom/beacon/downloads",
		  "issues_url": "https://api.github.com/repos/carsdotcom/beacon/issues{/number}",
		  "pulls_url": "https://api.github.com/repos/carsdotcom/beacon/pulls{/number}",
		  "milestones_url": "https://api.github.com/repos/carsdotcom/beacon/milestones{/number}",
		  "notifications_url": "https://api.github.com/repos/carsdotcom/beacon/notifications{?since,all,participating}",
		  "labels_url": "https://api.github.com/repos/carsdotcom/beacon/labels{/name}",
		  "releases_url": "https://api.github.com/repos/carsdotcom/beacon/releases{/id}",
		  "deployments_url": "https://api.github.com/repos/carsdotcom/beacon/deployments",
		  "created_at": "2016-11-03T15:56:26Z",
		  "updated_at": "2016-11-03T16:11:05Z",
		  "pushed_at": "2016-11-03T16:11:04Z",
		  "git_url": "git://github.com/carsdotcom/beacon.git",
		  "ssh_url": "git@github.com:carsdotcom/beacon.git",
		  "clone_url": "https://github.com/carsdotcom/beacon.git",
		  "svn_url": "https://github.com/carsdotcom/beacon",
		  "homepage": null,
		  "size": 7,
		  "stargazers_count": 0,
		  "watchers_count": 0,
		  "language": "Arduino",
		  "has_issues": true,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": false,
		  "forks_count": 0,
		  "mirror_url": null,
		  "archived": false,
		  "open_issues_count": 0,
		  "license": {
			"key": "apache-2.0",
			"name": "Apache License 2.0",
			"spdx_id": "Apache-2.0",
			"url": "https://api.github.com/licenses/apache-2.0",
			"node_id": "MDc6TGljZW5zZTI="
		  },
		  "forks": 0,
		  "open_issues": 0,
		  "watchers": 0,
		  "default_branch": "master",
		  "permissions": {
			"admin": true,
			"push": true,
			"pull": true
		  }
		},
		{
		  "id": 81501065,
		  "node_id": "MDEwOlJlcG9zaXRvcnk4MTUwMTA2NQ==",
		  "name": "bitcar",
		  "full_name": "carsdotcom/bitcar",
		  "owner": {
			"login": "carsdotcom",
			"id": 8259181,
			"node_id": "MDEyOk9yZ2FuaXphdGlvbjgyNTkxODE=",
			"avatar_url": "https://avatars2.githubusercontent.com/u/8259181?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/carsdotcom",
			"html_url": "https://github.com/carsdotcom",
			"followers_url": "https://api.github.com/users/carsdotcom/followers",
			"following_url": "https://api.github.com/users/carsdotcom/following{/other_user}",
			"gists_url": "https://api.github.com/users/carsdotcom/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/carsdotcom/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/carsdotcom/subscriptions",
			"organizations_url": "https://api.github.com/users/carsdotcom/orgs",
			"repos_url": "https://api.github.com/users/carsdotcom/repos",
			"events_url": "https://api.github.com/users/carsdotcom/events{/privacy}",
			"received_events_url": "https://api.github.com/users/carsdotcom/received_events",
			"type": "Organization",
			"site_admin": false
		  },
		  "private": false,
		  "html_url": "https://github.com/carsdotcom/bitcar",
		  "description": "seemlessly jump between repos from the command line",
		  "fork": false,
		  "url": "https://api.github.com/repos/carsdotcom/bitcar",
		  "forks_url": "https://api.github.com/repos/carsdotcom/bitcar/forks",
		  "keys_url": "https://api.github.com/repos/carsdotcom/bitcar/keys{/key_id}",
		  "collaborators_url": "https://api.github.com/repos/carsdotcom/bitcar/collaborators{/collaborator}",
		  "teams_url": "https://api.github.com/repos/carsdotcom/bitcar/teams",
		  "hooks_url": "https://api.github.com/repos/carsdotcom/bitcar/hooks",
		  "issue_events_url": "https://api.github.com/repos/carsdotcom/bitcar/issues/events{/number}",
		  "events_url": "https://api.github.com/repos/carsdotcom/bitcar/events",
		  "assignees_url": "https://api.github.com/repos/carsdotcom/bitcar/assignees{/user}",
		  "branches_url": "https://api.github.com/repos/carsdotcom/bitcar/branches{/branch}",
		  "tags_url": "https://api.github.com/repos/carsdotcom/bitcar/tags",
		  "blobs_url": "https://api.github.com/repos/carsdotcom/bitcar/git/blobs{/sha}",
		  "git_tags_url": "https://api.github.com/repos/carsdotcom/bitcar/git/tags{/sha}",
		  "git_refs_url": "https://api.github.com/repos/carsdotcom/bitcar/git/refs{/sha}",
		  "trees_url": "https://api.github.com/repos/carsdotcom/bitcar/git/trees{/sha}",
		  "statuses_url": "https://api.github.com/repos/carsdotcom/bitcar/statuses/{sha}",
		  "languages_url": "https://api.github.com/repos/carsdotcom/bitcar/languages",
		  "stargazers_url": "https://api.github.com/repos/carsdotcom/bitcar/stargazers",
		  "contributors_url": "https://api.github.com/repos/carsdotcom/bitcar/contributors",
		  "subscribers_url": "https://api.github.com/repos/carsdotcom/bitcar/subscribers",
		  "subscription_url": "https://api.github.com/repos/carsdotcom/bitcar/subscription",
		  "commits_url": "https://api.github.com/repos/carsdotcom/bitcar/commits{/sha}",
		  "git_commits_url": "https://api.github.com/repos/carsdotcom/bitcar/git/commits{/sha}",
		  "comments_url": "https://api.github.com/repos/carsdotcom/bitcar/comments{/number}",
		  "issue_comment_url": "https://api.github.com/repos/carsdotcom/bitcar/issues/comments{/number}",
		  "contents_url": "https://api.github.com/repos/carsdotcom/bitcar/contents/{+path}",
		  "compare_url": "https://api.github.com/repos/carsdotcom/bitcar/compare/{base}...{head}",
		  "merges_url": "https://api.github.com/repos/carsdotcom/bitcar/merges",
		  "archive_url": "https://api.github.com/repos/carsdotcom/bitcar/{archive_format}{/ref}",
		  "downloads_url": "https://api.github.com/repos/carsdotcom/bitcar/downloads",
		  "issues_url": "https://api.github.com/repos/carsdotcom/bitcar/issues{/number}",
		  "pulls_url": "https://api.github.com/repos/carsdotcom/bitcar/pulls{/number}",
		  "milestones_url": "https://api.github.com/repos/carsdotcom/bitcar/milestones{/number}",
		  "notifications_url": "https://api.github.com/repos/carsdotcom/bitcar/notifications{?since,all,participating}",
		  "labels_url": "https://api.github.com/repos/carsdotcom/bitcar/labels{/name}",
		  "releases_url": "https://api.github.com/repos/carsdotcom/bitcar/releases{/id}",
		  "deployments_url": "https://api.github.com/repos/carsdotcom/bitcar/deployments",
		  "created_at": "2017-02-09T22:24:31Z",
		  "updated_at": "2018-01-04T03:06:59Z",
		  "pushed_at": "2017-11-28T04:03:37Z",
		  "git_url": "git://github.com/carsdotcom/bitcar.git",
		  "ssh_url": "git@github.com:carsdotcom/bitcar.git",
		  "clone_url": "https://github.com/carsdotcom/bitcar.git",
		  "svn_url": "https://github.com/carsdotcom/bitcar",
		  "homepage": null,
		  "size": 1411,
		  "stargazers_count": 11,
		  "watchers_count": 11,
		  "language": "JavaScript",
		  "has_issues": true,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": false,
		  "forks_count": 3,
		  "mirror_url": null,
		  "archived": false,
		  "open_issues_count": 2,
		  "license": {
			"key": "apache-2.0",
			"name": "Apache License 2.0",
			"spdx_id": "Apache-2.0",
			"url": "https://api.github.com/licenses/apache-2.0",
			"node_id": "MDc6TGljZW5zZTI="
		  },
		  "forks": 3,
		  "open_issues": 2,
		  "watchers": 11,
		  "default_branch": "master",
		  "permissions": {
			"admin": true,
			"push": true,
			"pull": true
		  }
		},
		{
		  "id": 57904843,
		  "node_id": "MDEwOlJlcG9zaXRvcnk1NzkwNDg0Mw==",
		  "name": "cars.com-sellandtrade-application",
		  "full_name": "carsdotcom/cars.com-sellandtrade-application",
		  "owner": {
			"login": "carsdotcom",
			"id": 8259181,
			"node_id": "MDEyOk9yZ2FuaXphdGlvbjgyNTkxODE=",
			"avatar_url": "https://avatars2.githubusercontent.com/u/8259181?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/carsdotcom",
			"html_url": "https://github.com/carsdotcom",
			"followers_url": "https://api.github.com/users/carsdotcom/followers",
			"following_url": "https://api.github.com/users/carsdotcom/following{/other_user}",
			"gists_url": "https://api.github.com/users/carsdotcom/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/carsdotcom/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/carsdotcom/subscriptions",
			"organizations_url": "https://api.github.com/users/carsdotcom/orgs",
			"repos_url": "https://api.github.com/users/carsdotcom/repos",
			"events_url": "https://api.github.com/users/carsdotcom/events{/privacy}",
			"received_events_url": "https://api.github.com/users/carsdotcom/received_events",
			"type": "Organization",
			"site_admin": false
		  },
		  "private": false,
		  "html_url": "https://github.com/carsdotcom/cars.com-sellandtrade-application",
		  "description": null,
		  "fork": true,
		  "url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application",
		  "forks_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/forks",
		  "keys_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/keys{/key_id}",
		  "collaborators_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/collaborators{/collaborator}",
		  "teams_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/teams",
		  "hooks_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/hooks",
		  "issue_events_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/issues/events{/number}",
		  "events_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/events",
		  "assignees_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/assignees{/user}",
		  "branches_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/branches{/branch}",
		  "tags_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/tags",
		  "blobs_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/git/blobs{/sha}",
		  "git_tags_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/git/tags{/sha}",
		  "git_refs_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/git/refs{/sha}",
		  "trees_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/git/trees{/sha}",
		  "statuses_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/statuses/{sha}",
		  "languages_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/languages",
		  "stargazers_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/stargazers",
		  "contributors_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/contributors",
		  "subscribers_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/subscribers",
		  "subscription_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/subscription",
		  "commits_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/commits{/sha}",
		  "git_commits_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/git/commits{/sha}",
		  "comments_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/comments{/number}",
		  "issue_comment_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/issues/comments{/number}",
		  "contents_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/contents/{+path}",
		  "compare_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/compare/{base}...{head}",
		  "merges_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/merges",
		  "archive_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/{archive_format}{/ref}",
		  "downloads_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/downloads",
		  "issues_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/issues{/number}",
		  "pulls_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/pulls{/number}",
		  "milestones_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/milestones{/number}",
		  "notifications_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/notifications{?since,all,participating}",
		  "labels_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/labels{/name}",
		  "releases_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/releases{/id}",
		  "deployments_url": "https://api.github.com/repos/carsdotcom/cars.com-sellandtrade-application/deployments",
		  "created_at": "2016-05-02T16:44:11Z",
		  "updated_at": "2016-05-02T16:44:14Z",
		  "pushed_at": "2016-04-15T13:11:53Z",
		  "git_url": "git://github.com/carsdotcom/cars.com-sellandtrade-application.git",
		  "ssh_url": "git@github.com:carsdotcom/cars.com-sellandtrade-application.git",
		  "clone_url": "https://github.com/carsdotcom/cars.com-sellandtrade-application.git",
		  "svn_url": "https://github.com/carsdotcom/cars.com-sellandtrade-application",
		  "homepage": null,
		  "size": 102916,
		  "stargazers_count": 0,
		  "watchers_count": 0,
		  "language": "Java",
		  "has_issues": false,
		  "has_projects": true,
		  "has_downloads": true,
		  "has_wiki": true,
		  "has_pages": false,
		  "forks_count": 0,
		  "mirror_url": null,
		  "archived": false,
		  "open_issues_count": 0,
		  "license": null,
		  "forks": 0,
		  "open_issues": 0,
		  "watchers": 0,
		  "default_branch": "master",
		  "permissions": {
			"admin": true,
			"push": true,
			"pull": true
		  }
		}
	]`,
		`[
	{
		"id": 91115144,
		"node_id": "MDEwOlJlcG9zaXRvcnk5MTExNTE0NA==",
		"name": "blog-resources",
		"full_name": "TheHipbot/blog-resources",
		"owner": {
		"login": "TheHipbot",
		"id": 1820334,
		"node_id": "MDQ6VXNlcjE4MjAzMzQ=",
		"avatar_url": "https://avatars2.githubusercontent.com/u/1820334?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/TheHipbot",
		"html_url": "https://github.com/TheHipbot",
		"followers_url": "https://api.github.com/users/TheHipbot/followers",
		"following_url": "https://api.github.com/users/TheHipbot/following{/other_user}",
		"gists_url": "https://api.github.com/users/TheHipbot/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/TheHipbot/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/TheHipbot/subscriptions",
		"organizations_url": "https://api.github.com/users/TheHipbot/orgs",
		"repos_url": "https://api.github.com/users/TheHipbot/repos",
		"events_url": "https://api.github.com/users/TheHipbot/events{/privacy}",
		"received_events_url": "https://api.github.com/users/TheHipbot/received_events",
		"type": "User",
		"site_admin": false
		},
		"private": false,
		"html_url": "https://github.com/TheHipbot/blog-resources",
		"description": "Content and resources for blog posts I have made",
		"fork": false,
		"url": "https://api.github.com/repos/TheHipbot/blog-resources",
		"forks_url": "https://api.github.com/repos/TheHipbot/blog-resources/forks",
		"keys_url": "https://api.github.com/repos/TheHipbot/blog-resources/keys{/key_id}",
		"collaborators_url": "https://api.github.com/repos/TheHipbot/blog-resources/collaborators{/collaborator}",
		"teams_url": "https://api.github.com/repos/TheHipbot/blog-resources/teams",
		"hooks_url": "https://api.github.com/repos/TheHipbot/blog-resources/hooks",
		"issue_events_url": "https://api.github.com/repos/TheHipbot/blog-resources/issues/events{/number}",
		"events_url": "https://api.github.com/repos/TheHipbot/blog-resources/events",
		"assignees_url": "https://api.github.com/repos/TheHipbot/blog-resources/assignees{/user}",
		"branches_url": "https://api.github.com/repos/TheHipbot/blog-resources/branches{/branch}",
		"tags_url": "https://api.github.com/repos/TheHipbot/blog-resources/tags",
		"blobs_url": "https://api.github.com/repos/TheHipbot/blog-resources/git/blobs{/sha}",
		"git_tags_url": "https://api.github.com/repos/TheHipbot/blog-resources/git/tags{/sha}",
		"git_refs_url": "https://api.github.com/repos/TheHipbot/blog-resources/git/refs{/sha}",
		"trees_url": "https://api.github.com/repos/TheHipbot/blog-resources/git/trees{/sha}",
		"statuses_url": "https://api.github.com/repos/TheHipbot/blog-resources/statuses/{sha}",
		"languages_url": "https://api.github.com/repos/TheHipbot/blog-resources/languages",
		"stargazers_url": "https://api.github.com/repos/TheHipbot/blog-resources/stargazers",
		"contributors_url": "https://api.github.com/repos/TheHipbot/blog-resources/contributors",
		"subscribers_url": "https://api.github.com/repos/TheHipbot/blog-resources/subscribers",
		"subscription_url": "https://api.github.com/repos/TheHipbot/blog-resources/subscription",
		"commits_url": "https://api.github.com/repos/TheHipbot/blog-resources/commits{/sha}",
		"git_commits_url": "https://api.github.com/repos/TheHipbot/blog-resources/git/commits{/sha}",
		"comments_url": "https://api.github.com/repos/TheHipbot/blog-resources/comments{/number}",
		"issue_comment_url": "https://api.github.com/repos/TheHipbot/blog-resources/issues/comments{/number}",
		"contents_url": "https://api.github.com/repos/TheHipbot/blog-resources/contents/{+path}",
		"compare_url": "https://api.github.com/repos/TheHipbot/blog-resources/compare/{base}...{head}",
		"merges_url": "https://api.github.com/repos/TheHipbot/blog-resources/merges",
		"archive_url": "https://api.github.com/repos/TheHipbot/blog-resources/{archive_format}{/ref}",
		"downloads_url": "https://api.github.com/repos/TheHipbot/blog-resources/downloads",
		"issues_url": "https://api.github.com/repos/TheHipbot/blog-resources/issues{/number}",
		"pulls_url": "https://api.github.com/repos/TheHipbot/blog-resources/pulls{/number}",
		"milestones_url": "https://api.github.com/repos/TheHipbot/blog-resources/milestones{/number}",
		"notifications_url": "https://api.github.com/repos/TheHipbot/blog-resources/notifications{?since,all,participating}",
		"labels_url": "https://api.github.com/repos/TheHipbot/blog-resources/labels{/name}",
		"releases_url": "https://api.github.com/repos/TheHipbot/blog-resources/releases{/id}",
		"deployments_url": "https://api.github.com/repos/TheHipbot/blog-resources/deployments",
		"created_at": "2017-05-12T17:36:09Z",
		"updated_at": "2017-05-12T17:36:09Z",
		"pushed_at": "2017-05-12T17:41:27Z",
		"git_url": "git://github.com/TheHipbot/blog-resources.git",
		"ssh_url": "git@github.com:TheHipbot/blog-resources.git",
		"clone_url": "https://github.com/TheHipbot/blog-resources.git",
		"svn_url": "https://github.com/TheHipbot/blog-resources",
		"homepage": null,
		"size": 0,
		"stargazers_count": 0,
		"watchers_count": 0,
		"language": null,
		"has_issues": true,
		"has_projects": true,
		"has_downloads": true,
		"has_wiki": true,
		"has_pages": false,
		"forks_count": 0,
		"mirror_url": null,
		"archived": false,
		"open_issues_count": 0,
		"license": null,
		"forks": 0,
		"open_issues": 0,
		"watchers": 0,
		"default_branch": "master",
		"permissions": {
		"admin": true,
		"push": true,
		"pull": true
		}
	},
	{
		"id": 46675428,
		"node_id": "MDEwOlJlcG9zaXRvcnk0NjY3NTQyOA==",
		"name": "causal-relations",
		"full_name": "TheHipbot/causal-relations",
		"owner": {
		"login": "TheHipbot",
		"id": 1820334,
		"node_id": "MDQ6VXNlcjE4MjAzMzQ=",
		"avatar_url": "https://avatars2.githubusercontent.com/u/1820334?v=4",
		"gravatar_id": "",
		"url": "https://api.github.com/users/TheHipbot",
		"html_url": "https://github.com/TheHipbot",
		"followers_url": "https://api.github.com/users/TheHipbot/followers",
		"following_url": "https://api.github.com/users/TheHipbot/following{/other_user}",
		"gists_url": "https://api.github.com/users/TheHipbot/gists{/gist_id}",
		"starred_url": "https://api.github.com/users/TheHipbot/starred{/owner}{/repo}",
		"subscriptions_url": "https://api.github.com/users/TheHipbot/subscriptions",
		"organizations_url": "https://api.github.com/users/TheHipbot/orgs",
		"repos_url": "https://api.github.com/users/TheHipbot/repos",
		"events_url": "https://api.github.com/users/TheHipbot/events{/privacy}",
		"received_events_url": "https://api.github.com/users/TheHipbot/received_events",
		"type": "User",
		"site_admin": false
		},
		"private": false,
		"html_url": "https://github.com/TheHipbot/causal-relations",
		"description": "Causal Relations - An open source blog on all things computer science",
		"fork": false,
		"url": "https://api.github.com/repos/TheHipbot/causal-relations",
		"forks_url": "https://api.github.com/repos/TheHipbot/causal-relations/forks",
		"keys_url": "https://api.github.com/repos/TheHipbot/causal-relations/keys{/key_id}",
		"collaborators_url": "https://api.github.com/repos/TheHipbot/causal-relations/collaborators{/collaborator}",
		"teams_url": "https://api.github.com/repos/TheHipbot/causal-relations/teams",
		"hooks_url": "https://api.github.com/repos/TheHipbot/causal-relations/hooks",
		"issue_events_url": "https://api.github.com/repos/TheHipbot/causal-relations/issues/events{/number}",
		"events_url": "https://api.github.com/repos/TheHipbot/causal-relations/events",
		"assignees_url": "https://api.github.com/repos/TheHipbot/causal-relations/assignees{/user}",
		"branches_url": "https://api.github.com/repos/TheHipbot/causal-relations/branches{/branch}",
		"tags_url": "https://api.github.com/repos/TheHipbot/causal-relations/tags",
		"blobs_url": "https://api.github.com/repos/TheHipbot/causal-relations/git/blobs{/sha}",
		"git_tags_url": "https://api.github.com/repos/TheHipbot/causal-relations/git/tags{/sha}",
		"git_refs_url": "https://api.github.com/repos/TheHipbot/causal-relations/git/refs{/sha}",
		"trees_url": "https://api.github.com/repos/TheHipbot/causal-relations/git/trees{/sha}",
		"statuses_url": "https://api.github.com/repos/TheHipbot/causal-relations/statuses/{sha}",
		"languages_url": "https://api.github.com/repos/TheHipbot/causal-relations/languages",
		"stargazers_url": "https://api.github.com/repos/TheHipbot/causal-relations/stargazers",
		"contributors_url": "https://api.github.com/repos/TheHipbot/causal-relations/contributors",
		"subscribers_url": "https://api.github.com/repos/TheHipbot/causal-relations/subscribers",
		"subscription_url": "https://api.github.com/repos/TheHipbot/causal-relations/subscription",
		"commits_url": "https://api.github.com/repos/TheHipbot/causal-relations/commits{/sha}",
		"git_commits_url": "https://api.github.com/repos/TheHipbot/causal-relations/git/commits{/sha}",
		"comments_url": "https://api.github.com/repos/TheHipbot/causal-relations/comments{/number}",
		"issue_comment_url": "https://api.github.com/repos/TheHipbot/causal-relations/issues/comments{/number}",
		"contents_url": "https://api.github.com/repos/TheHipbot/causal-relations/contents/{+path}",
		"compare_url": "https://api.github.com/repos/TheHipbot/causal-relations/compare/{base}...{head}",
		"merges_url": "https://api.github.com/repos/TheHipbot/causal-relations/merges",
		"archive_url": "https://api.github.com/repos/TheHipbot/causal-relations/{archive_format}{/ref}",
		"downloads_url": "https://api.github.com/repos/TheHipbot/causal-relations/downloads",
		"issues_url": "https://api.github.com/repos/TheHipbot/causal-relations/issues{/number}",
		"pulls_url": "https://api.github.com/repos/TheHipbot/causal-relations/pulls{/number}",
		"milestones_url": "https://api.github.com/repos/TheHipbot/causal-relations/milestones{/number}",
		"notifications_url": "https://api.github.com/repos/TheHipbot/causal-relations/notifications{?since,all,participating}",
		"labels_url": "https://api.github.com/repos/TheHipbot/causal-relations/labels{/name}",
		"releases_url": "https://api.github.com/repos/TheHipbot/causal-relations/releases{/id}",
		"deployments_url": "https://api.github.com/repos/TheHipbot/causal-relations/deployments",
		"created_at": "2015-11-22T19:11:12Z",
		"updated_at": "2016-02-09T01:48:38Z",
		"pushed_at": "2016-02-09T04:27:49Z",
		"git_url": "git://github.com/TheHipbot/causal-relations.git",
		"ssh_url": "git@github.com:TheHipbot/causal-relations.git",
		"clone_url": "https://github.com/TheHipbot/causal-relations.git",
		"svn_url": "https://github.com/TheHipbot/causal-relations",
		"homepage": null,
		"size": 7,
		"stargazers_count": 0,
		"watchers_count": 0,
		"language": "HTML",
		"has_issues": true,
		"has_projects": true,
		"has_downloads": true,
		"has_wiki": true,
		"has_pages": false,
		"forks_count": 0,
		"mirror_url": null,
		"archived": false,
		"open_issues_count": 0,
		"license": null,
		"forks": 0,
		"open_issues": 0,
		"watchers": 0,
		"default_branch": "master",
		"permissions": {
		"admin": true,
		"push": true,
		"pull": true
		}
	}
]`,
	}
)
