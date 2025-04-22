# Contributing to aep-go

I welcome your input! I want to make contributing to aep-go as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development with GitHub
I use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## Making Changes
Pull requests are the best way to propose changes to the codebase. I actively welcome your pull requests:

1. Create a new branch from `master` (e.g., `feature/your-feature-name` or `fix/your-bug-fix`)
2. Make your changes
3. If you've added code that should be tested, add tests
4. If you've changed APIs, update the documentation
5. Ensure the test suite passes
6. Make sure your code lints
7. Push your branch and create a pull request

## Any contributions you make will be under the MIT Software License
In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact me if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/blaberg/aep-go/issues)
I use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/blaberg/aep-go/issues/new); it's that easy!

## Write bug reports with detail, background, and sample code

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Use a Consistent Coding Style

* Use `gofmt` for formatting
* Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* Write tests for new code
* Update documentation for API changes

## Commit Messages

I follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages. This helps maintain a clear and consistent commit history. Please format your commit messages as follows:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Common types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

## License
By contributing, you agree that your contributions will be licensed under its MIT License. 