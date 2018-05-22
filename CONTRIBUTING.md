# Contribute to Hermes

## Reporting Issues

Before reporting an issue, please search the open issues for the project [here](https://github.com/TheHipbot/hermes/issues) to ensure duplicate issues aren't created. Also, for bugs try updating the project to see if it has not already been fixed before submitting.

## Conventions

Branches should include the name of the associated github issue. 

## Development Workflow

### Setup

1. Create a fork of hermes and ensure the master branch on the fork is up to date with the upstream
2. Clone down your fork and the hermes project itself
3. Create a branch in your fork and push it up to your remote
4. Go into your local copy of the hermes repository (the main project not your fork) and run 
    ```bash
    git remote add fork https://github.com/<your user name>/hermes
    git fetch fork
    git checkout -t fork/<name of your branch>
    ```

5. Now you should be in your branch but in the hermes repository (this is necessary so that the go imports will work correctly). Run `dep ensure` or `make ensure` to install all vendored dependencies

### Workflow

1. Make your changes along with tests. You can run tests either with `go test ./...` or you could use the make target by running `make test`.
2. Before committing, always make sure that all tests pass, run `dep ensure` or `make ensure` to ensure all dependencies are captured, and run `go fmt ./...`
3. Commit using the `-s` to sign your commit
4. When your pull request is ready to be merged, first sqash any commits made on your branch so the PR contains a single commit then submit