package containerenv

// ConfigFileV1 is a config object for an environment
type ConfigFileV1 struct {
	Version int `yaml:"version"`

	Environment struct {
		// Name is the name of this environment
		Name string `yaml:"name"`

		// Base is the base image to use
		Base string `yaml:"base"`

		// Username is the username to use for this image
		Username string `yaml:"username"`

		// Options contains toggleable features
		Options struct {
			// PulseAudio enables pulseaudio features
			PulseAudio bool `yaml:"pulseaudio"`

			// X11 enables xorg support
			X11 bool `yaml:"x11"`

			// SystemD enables systemd support
			SystemD bool `yaml:"systemd"`
		}
	}
}

// Environment is a containerized user environment that should be run
type Environment struct {
	// Name of the environment
	Name string `json:"name"`

	// Username is the user we should run as in this container
	Username string `json:"username"`

	// SystemD toggles support for systemd, defaults to true
	SystemD bool `json:"systemd"`

	// Image is the Docker image used to run this environment.
	Image string `json:"image"`

	// PulseAudio configures the pulseaudio integration
	PulseAudio PulseAudioSettings `json:"pulseaudio"`

	// X11 enables suppot for X11. This requires X11 to be running on the host
	X11 X11Settings `json:"x11"`
}

// PulseAudioSettings configures pulseaudio
type PulseAudioSettings struct {
	// Host provides host pulseaudio access
	Host bool `json:"host"`

	// Containerized specifies we are going to be running pulseaudio in the container
	Containerized bool `json:"containerized"`
}

// X11Settings configures X11
type X11Settings struct {
	// Host provides host x11 access
	Host bool `json:"host"`

	// Containerized specifies we are going to be running xorg-server in the container
	Containerized bool `json:"containerized"`
}
