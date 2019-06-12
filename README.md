# containerenv

Docker-powered environment bundling platform.

## What is this?

Using the power of Docker we can now bundle our developer environments into containers! This is particularly useful when we want to run another distro
ontop of a specially tuned host.

## Is this secure?

Sortof. Depending on the options you provide it's less secure. Using xorg in the container is currently not very secure due to it's need for root access (limitation of Xorg). So it's important to note that this isn't sandbox friendly or non-privileged user friendly


## Usage

Download a [Release](#releases). Then look at the configuration manifests located at `./contrib/designs/v1.yaml` for an idea of what to provide to `create`.

## License

MIT