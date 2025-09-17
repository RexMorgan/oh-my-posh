package segments

type AppleMusic struct {
	Base

	MusicPlayer
}

func (a *AppleMusic) Template() string {
	return " {{ .Icon }}{{ if ne .Status \"stopped\" }}{{ .Artist }} - {{ .Track }}{{ end }} "
}

func (a *AppleMusic) resolveIcon() {
	switch a.Status {
	case stopped:
		// in this case, no artist or track info
		a.Icon = a.props.GetString(StoppedIcon, "\uF04D ")
	case paused:
		a.Icon = a.props.GetString(PausedIcon, "\uF8E3 ")
	case playing:
		a.Icon = a.props.GetString(PlayingIcon, "\uE602 ")
	}
}
