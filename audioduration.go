// mautrix-imessage - A Matrix-iMessage puppeting bridge.
// Copyright (C) 2026 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// probeAudioDurationMS returns the duration of the given audio data in
// milliseconds.
//
// Voice messages need this because `duration` is a required field of the
// org.matrix.msc1767.audio content block for ruma-based clients such as
// Element X. Sending the block without it makes the entire event fail to
// deserialize with "missing field `duration`", and the message renders as an
// "Unsupported event" placeholder rather than degrading to a playable file.
//
// ffprobe comes from the same package as ffmpeg, which the bridge already
// requires to convert CAF to Ogg/Opus, so this adds no new dependency.
func probeAudioDurationMS(ctx context.Context, data []byte, extension string) (int, error) {
	tmp, err := os.CreateTemp("", "mautrix-imessage-probe-*"+extension)
	if err != nil {
		return 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err = tmp.Write(data); err != nil {
		_ = tmp.Close()
		return 0, fmt.Errorf("failed to write temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return 0, fmt.Errorf("failed to close temp file: %w", err)
	}

	out, err := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		tmp.Name(),
	).Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	// ffprobe prints "N/A" when it can't measure the stream.
	raw := strings.TrimSpace(string(out))
	seconds, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("unparseable ffprobe duration %q", raw)
	}
	if seconds <= 0 || math.IsInf(seconds, 0) || math.IsNaN(seconds) {
		return 0, fmt.Errorf("implausible ffprobe duration %q", raw)
	}

	return int(math.Round(seconds * 1000)), nil
}
