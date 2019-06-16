# containerenv

Docker-powered environment bundling platform.

## What is this?

Using the power of Docker we can now bundle our developer environments into containers! This is particularly useful when we want to run another distro
ontop of a specially tuned host.

## Is this secure?

Sortof. Depending on the options you provide it's less secure. Using xorg in the container is currently not very secure due to it's need for root access (limitation of Xorg). So it's important to note that this isn't sandbox friendly or non-privileged user friendly


## Usage

Download a [Release](releases).

Run `containerenv init` to create an environment configuration file. Keep in mind that you will need push access to Docker or some other registry to make this
useful.

Provide that config to `containerenv create <env.yaml>`. This will create the container on your local host.

Exec into the container `containerenv exec <name of env>`. You're now able to use whatever options you turned on!

## FAQ

### How does X11 forwarding work?

If it's detected that X11 is not currently running then your container is configured to run X11 inside of the container. If it's detected that X11 is currently running,
then it is configured to X11 forward, which allows applications to run on host X, but not run window managers.

Either case, the env var `X11_CONFIG` is set to `HOST` or `CONTAINER` respective to both options.

### How does pulseaudio work>

By default, only pulseaudio forwarding currently works.

Either case, the env var `PULSEAUDIO_CONFIG` is set to `HOST` or `CONTAINER` respective to both options.



## License

MIT