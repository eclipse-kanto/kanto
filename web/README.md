![Kanto logo](https://github.com/eclipse-kanto/kanto/raw/master/logo/kanto.svg)

# Eclipse Kanto Site and Documentation

This is the repository used to build and publish the official Eclipse Kanto [website](https://eclipse.org/kanto/).

## Quick start

We use [Hugo](https://gohugo.io/) and the [Docsy theme](https://docsy.dev/)
to build and render the site. You will need the “extended”
Sass/SCSS version of the Hugo binary to work on this site. We recommend
to use Hugo 0.83.1 or higher.

Steps needed to have this working locally and work on it:

- Follow the [Install Hugo](https://docsy.dev/docs/getting-started/#install-hugo) instructions from Docsy
- Clone this repository
- Run `git submodule update --init --recursive`
- Run `cd site`
- Run `hugo server`
