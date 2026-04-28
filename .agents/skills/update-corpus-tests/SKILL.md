---
name: update-corpus-tests
description: Updating golden files used for tests
---

+ In order to update golden files used for tests run the tests with `--update` flag.
  + example: `go test ./corpus --update`
+ You will get test failures for any file that got updated.
+ Then use git diff on all updated golden files to confirm changes matches with the expectations
