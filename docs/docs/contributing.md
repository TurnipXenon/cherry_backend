# Contributing to Cherry Backend

Thank you for considering contributing to Cherry Backend! This document outlines the guidelines for contributing to this project.

## Commit Message Convention

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for our commit messages. This leads to more readable messages that are easy to follow when looking through the project history.

### Commit Message Format

Each commit message consists of a **header**, a **body**, and a **footer**. The header has a special format that includes a **type**, a **scope**, and a **subject**:

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

The **header** is mandatory and the **scope** of the header is optional.

### Type

Must be one of the following:

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools and libraries such as documentation generation
- **ci**: Changes to our CI configuration files and scripts
- **build**: Changes that affect the build system or external dependencies
- **revert**: Reverts a previous commit

### Scope

The scope could be anything specifying the place of the commit change. For example `api`, `model`, `controller`, etc.

### Subject

The subject contains a succinct description of the change:

- Use the imperative, present tense: "change" not "changed" nor "changes"
- Don't capitalize the first letter
- No dot (.) at the end

### Body

The body should include the motivation for the change and contrast this with previous behavior.

### Footer

The footer should contain any information about **Breaking Changes** and is also the place to reference GitHub issues that this commit **Closes**.

### Examples

```
feat(auth): add ability to login with Google

Implement Google OAuth login flow using the Google API client library.

Closes #123
```

```
fix(api): prevent race condition in request handler

The request handler was not properly synchronized, which could lead to data corruption.

Fixes #456
```

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations, and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent.
4. The Pull Request will be merged once you have the sign-off of at least one other developer, or if you do not have permission to do that, you may request the reviewer to merge it for you.


## MKDocs

You don't really need to run through this to make edits to MKDocs, but if you want to see the layout and what it looks
like served, check this guide.

This assumes that you have **Python** installed locally.

### MKDocs: Setup

We use Material for MkDocs theme and several extensions. To set up your local environment, install the dependencies from the requirements file:

```shell
# From the project root
pip install -r docs/requirements.txt
```

This ensures you're using the same versions of packages as the CI/CD pipeline and other contributors.

#### Updating Dependencies

If you need to add or update dependencies for the documentation:

1. Update the `docs/requirements.txt` file with the new package or version
2. Test locally to ensure everything works as expected
3. Commit the changes to the requirements file

We pin specific versions in the requirements file to ensure consistent builds across different environments.

### MKDocs: Commands

When entering these commands, go to `/docs` instead of being in the project's root folder `/`.

* `mkdocs serve` - Start the live-reloading docs server.
* `mkdocs build` - Build the documentation site.
* `mkdocs -h` - Print help message and exit.

### MKDocs: Ideal workflow

1. Make changes
2. See changes made using `mkdocs serve`
3. Commit and push your changes to the main branch

### Material for MkDocs Features

Our documentation uses Material for MkDocs with several features enabled:

#### Theme Features

- **Navigation Instant**: Provides faster page loading
- **Navigation Tracking**: Automatically focuses the current section in the navigation
- **Navigation Expand**: Expands all collapsible sections in the navigation
- **Navigation Indexes**: Adds index pages to sections
- **Content Code Copy**: Adds a copy button to code blocks

#### Markdown Extensions

We've enabled several markdown extensions that you can use in your documentation:

- **Admonitions**: Create call-out blocks for notes, warnings, etc.
  ```markdown
  !!! note "Optional Title"
      This is a note admonition.
  ```

- **Code Highlighting**: Syntax highlighting for code blocks
  ```markdown
  ```python
  def hello_world():
      print("Hello, world!")
  ```
  ```

- **Tabbed Content**: Create tabbed content sections
  ```markdown
  === "Tab 1"
      Content of tab 1

  === "Tab 2"
      Content of tab 2
  ```

- **Tables**: Create tables using markdown syntax
  ```markdown
  | Header 1 | Header 2 |
  | -------- | -------- |
  | Cell 1   | Cell 2   |
  ```

For more details on how to use these features, refer to the [Material for MkDocs documentation](https://squidfunk.github.io/mkdocs-material/).

### MKDocs: Automated Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the main branch. This is handled by a GitHub Actions workflow that:

1. Builds the MkDocs site
2. Deploys it to the gh-pages branch
3. Makes it available at the project's GitHub Pages URL

You don't need to manually build or deploy the documentation - just push your changes to the main branch, and the workflow will handle the rest.
