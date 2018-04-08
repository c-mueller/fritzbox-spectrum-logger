#!/usr/bin/env fish
set revision (git rev-parse HEAD)
set build_no (git rev-list --all --count)
set build_time (date)
set version (git describe --tags --exact-match; or git symbolic-ref -q --short HEAD)

if test -n "$DRONE_PREV_BUILD_NUMBER"
  set build_no (math $DRONE_PREV_BUILD_NUMBER + 1)
end
if test -n "$DRONE_COMMIT_BRANCH"
  set version $DRONE_COMMIT_BRANCH
end

go build -v -ldflags "-X main.version=$version -X main.revision=$revision -X main.buildNumber=$build_no -X \"main.buildTimestamp=$build_time\""