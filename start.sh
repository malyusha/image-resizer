#!/bin/sh

LOG_LEVEL={"$LOG_LEVEL":"fatal"}

resizer -d "$STORAGE_DIR" -source "$SOURCE_DIR" -presets "$RESIZE_PRESETS_PATH" -addr 0.0.0.0 -log "$LOG_LEVEL"