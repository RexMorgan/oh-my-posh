//go:build darwin

package segments

import "strings"

func (a *AppleMusic) Enabled() bool {
	// Batching commands to reduce latency. Each individual call to `osascript` creates additional delays.
	// Using '|' as a delimiter in the batched command since it's unlikely
	// to appear in track or artist names, making it safe for splitting the output
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

	batchedOutput := a.runAppleScriptCommand(batchedCommand)

	outputStrings := strings.SplitN(batchedOutput, "|", 4)
	if outputStrings[0] == "false" || outputStrings[0] == "" || len(outputStrings) != 4 {
		a.Status = stopped
		return false
	}

	a.Status = outputStrings[1]

	// Check if running
	if a.Status == "" {
		a.Status = stopped
		return false
	}

	if a.Status == stopped {
		return false
	}

	a.Artist = outputStrings[2]
	a.Track = outputStrings[3]
	a.resolveIcon()

	return true
}

func (a *AppleMusic) runAppleScriptCommand(command string) string {
	val, _ := a.env.RunCommand("osascript", "-e", command)
	return val
}
