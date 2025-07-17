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

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as contributors and maintainers pledge to making participation in our project and our community a harassment-free experience for everyone.

### Our Standards

Examples of behavior that contributes to creating a positive environment include:

- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable behavior and are expected to take appropriate and fair corrective action in response to any instances of unacceptable behavior.