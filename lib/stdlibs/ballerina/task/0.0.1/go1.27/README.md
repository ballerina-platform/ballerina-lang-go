# Ballerina Task Library

## Overview

The `ballerina/task` library schedules Ballerina jobs to run once at a given time or repeatedly at a fixed interval, manages scheduled jobs (pause, resume, unschedule, list), runs task services on a `task:Listener`, and coordinates warm-backup task nodes through a database. The Go Native Interpreter supports job scheduling and management, the worker-pool configuration, and the listener lifecycle; attaching services to a listener and database task coordination are not yet supported.

## Key Functionalities

- Schedule a one-time job at a `time:Civil` time (`scheduleOneTimeJob`)
- Schedule a recurring job with an interval, maximum count, start and end time (`scheduleJobRecurByFrequency`)
- Error policies (`LOG_AND_TERMINATE`, `LOG_AND_CONTINUE`, `TERMINATE`, `CONTINUE`) and waiting policies (`WAIT`, `IGNORE`, `LOG_AND_IGNORE`) following the Quartz scheduler semantics jBallerina uses
- Pause, resume and unschedule jobs individually or all at once (`pauseJob`, `resumeJob`, `pauseAllJobs`, `resumeAllJobs`, `unscheduleJob`)
- List the running jobs (`getRunningJobs`)
- Configure the worker pool and the misfire waiting time (`configureWorkerPool`)
- Convert a `time:Civil` to epoch milliseconds (`getTimeInMillies`)
- Create, start and stop a `task:Listener`

## Examples

```ballerina
import ballerina/io;
import ballerina/lang.runtime;
import ballerina/task;

isolated int fired = 0;

isolated class CounterJob {
    *task:Job;

    public isolated function execute() {
        lock {
            fired += 1;
        }
    }
}

public function main() returns error? {
    task:JobId id = check task:scheduleJobRecurByFrequency(new CounterJob(), 0.1, maxCount = 3);
    runtime:sleep(0.5);
    lock {
        io:println(fired); // 3
    }
    io:println(task:getRunningJobs().length()); // 0

    task:Error? e = task:unscheduleJob(id);
    if e is task:Error {
        io:println(e.message()); // Invalid job id: <id>
    }
}
```

## Go Native Interpreter Support Status

This library is currently being migrated to Go to support the Ballerina Native Interpreter. The table below outlines the current support level for various features of this library in the Go implementation.

Support Levels:

- **Supported**: Fully implemented and tested in the Go version.
- **Partially Supported**: Implemented but lacking some edge cases, options, or sub-features. (See comments).
- **Not Yet Supported**: Planned for migration, but not yet implemented.
- **Cannot Support**: Cannot be implemented in the Go version due to technical limitations or architectural differences. (See comments).

| Feature/API | Support Status | Comments / Limitations |
|---|---|---|
| Module error type | Supported | `task:Error` is a `distinct error` |
| Job abstraction | Supported | `task:Job`, `task:JobId` |
| One-time job scheduling | Supported | `scheduleOneTimeJob`; a panic in `execute` follows the `LOG_AND_CONTINUE` policy |
| Recurring job scheduling | Supported | `scheduleJobRecurByFrequency` with `interval`, `maxCount`, `startTime`, `endTime` |
| Error policies | Supported | A panic in `Job.execute` is passed to the policy |
| Waiting policies | Supported | `WAIT` fires every missed run; `IGNORE` and `LOG_AND_IGNORE` follow Quartz's smart misfire policy |
| Worker pool configuration | Supported | `configureWorkerPool`; as in jBallerina, each run is handed to its own strand, so the worker count is validated but does not bound concurrent executions |
| Pause and resume a job | Supported | `pauseJob`, `resumeJob` |
| Pause and resume all jobs | Supported | `pauseAllJobs`, `resumeAllJobs` |
| Unschedule a job | Supported | `unscheduleJob`; the scheduler shuts down when the last job is unscheduled |
| Running jobs listing | Supported | `getRunningJobs` |
| Civil time to milliseconds | Supported | `getTimeInMillies`; the `Couldn't convert given time to milli seconds` message ends with the error text of `ballerina/time`'s `utcFromCivil`, which differs from jBallerina |
| Global scheduler configuration variables | Partially Supported | `globalSchedulerWorkerCount` and `globalSchedulerWaitingTime` are plain public module variables because `configurable` variables are not yet supported by the interpreter; jBallerina's natives never read them either |
| Listener lifecycle | Supported | `task:Listener` `init`, `start`, `gracefulStop`, `immediateStop` |
| Attaching services to a listener | Not Yet Supported | `task:Service`, `Listener.attach`, `Listener.detach` are implemented natively, but user code cannot yet create a `task:Service` value: the interpreter rejects `service ... on taskListener` bodies because `execute` is a non-public member of another module, and `service object {}` constructors are not implemented |
| Retry configuration | Not Yet Supported | `task:RetryConfiguration` is implemented natively for listener services and blocked on attaching services |
| Task coordination | Not Yet Supported | `task:WarmBackupConfig`, `task:MysqlConfig`, `task:PostgresqlConfig` are declared; starting a listener with a `warmBackupConfig` and attached services returns a `task:Error` because database coordination needs SQL drivers that are not available |

### Notable Behavioural Changes

- **Errors print without the native error form.** jBallerina creates scheduler errors in Java, so `io:println(err)` shows `error("Invalid job id: 5")`; the Go-native version creates them in Ballerina, so the same error prints as `error Error ("Invalid job id: 5")` — the message, the `task:Error` type and `is` checks are identical.
- **Pausing or resuming all jobs before any job was scheduled returns an error.** jBallerina panics with a `NullPointerException` because no scheduler exists yet; the Go-native version returns a `task:Error` with the message `The Scheduler has been shutdown.`.
- **An out-of-range job id returns an error.** jBallerina panics with an integer-overflow `ArithmeticException` when a `JobId` outside the 32-bit range is paused, resumed or unscheduled; the Go-native version returns `Invalid job id: <id>`.
- **Job ids never collide.** jBallerina draws a random id below 1,000,000 and scheduling fails if it is already taken; the Go-native version draws again until the id is unused.
- **Running jobs are listed in ascending id order.** jBallerina returns them in hash-map order; the Go-native version sorts them — the order was never specified.
- **The skipped-trigger log uses Go's zone abbreviation.** The `LOG_AND_IGNORE` message prints the trigger start time in Java's `Date.toString` layout; the Go-native version uses the same layout with the zone abbreviation from the Go time zone database, which can differ (for example `+0530` instead of `IST`).
- **A panic in a listener service is reported as a single line.** jBallerina prints an uncaught-exception stack trace from the virtual thread running the service; the Go-native version recovers the panic and prints `error: <message>` to stderr.
- **Non-isolated jobs run concurrently.** jBallerina runs a job whose class or `execute` method is not isolated under the runtime's lock for non-isolated strands; the Go-native version runs every execution on its own strand, like `http` resources — use isolated jobs and `lock` statements for shared state.
