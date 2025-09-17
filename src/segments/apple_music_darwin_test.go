//go:build darwin

package segments

import (
	"errors"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
)

func TestAppleMusicDarwinEnabledAndSpotifyPlaying(t *testing.T) {
	cases := []struct {
		Error       error
		BatchedCase string
		Expected    string
		Enabled     bool
	}{
		{BatchedCase: "false|||", Expected: "", Enabled: false},
		{BatchedCase: "false||", Expected: "", Error: errors.New("oops"), Enabled: false},
		{BatchedCase: "true|playing|Candlemass|Spellbreaker", Expected: "\ue602 Candlemass - Spellbreaker", Enabled: true},
		{BatchedCase: "true|paused|Candlemass|Spellbreaker", Expected: "\uF8E3 Candlemass - Spellbreaker", Enabled: true},
	}
	batchedCommand := `
	if application "Music" is running then
		tell application "Music"
			set playerState to player state as string
			set artistName to ""
			set trackName to ""
			if playerState is not "stopped" then
				set artistName to artist of current track as string
				set trackName to name of current track as string
			end if
			return "true|" & playerState & "|" & artistName & "|" & trackName
		end tell
	else
		return "false|||"
	end if
	`
	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("RunCommand", "osascript", []string{"-e", batchedCommand}).Return(tc.BatchedCase, tc.Error)

		a := &AppleMusic{}
		a.Init(properties.Map{}, env)

		assert.Equal(t, tc.Enabled, a.Enabled())
		assert.Equal(t, tc.Expected, renderTemplate(env, a.Template(), a))
	}
}
