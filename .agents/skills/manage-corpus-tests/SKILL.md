---
name: manage-corpus-tests
description: Creating/updating corpus tests
---

## Test markers
+ corpus tests use the following comments as markers
  + `@output <expected output>`
    - Test harness parses the file top to bottom extracting the expected output and compares it against stdout.
    - Generally it is a good idea to put this right next to the print function call
  + `@error`
    - Test harness validates that each frontend error covers one of these markers
      - For errors that are covered by multiple lines it is sufficient to have one marker in one of those lines
    - IMPORTANT: Test harness doesn't validate error messages
  + `@panic`
    - When there is a runtime panic, test harness validates that the top stack frame location (file:line) matches this annotation

## Updating corpus tests
+ In order to update golden files used for tests, run the tests with `--update` flag.
  + example: `go test ./corpus --update`
+ You will get test failures for any file that got updated.
+ Then use git diff on all updated golden files to confirm changes match with the expectations

## Validating corpus tests
+ It is a good idea to validate output by running corpus files against the java implementation using `bal run $file`
