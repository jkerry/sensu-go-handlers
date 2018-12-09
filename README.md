# sensu_gcp_pubsub_handler

sensu_gcp_pubsub_handler is a metrics handler for sensu go that manages sending metric data to a Google Cloud Platform PubSub topic. From there one can marshal the data to any number of GCP services using dataflow jobs.

sensu_gcp_pubsub_handler uses the cloud.google.com/go/pubsub pubsub library. This library assumes that you have a GCP application credentials file on disk with the environment variable GOOGLE_APPLICATION_CREDENTIALS set pointed to it's location OR your instance is running in a GCP compute environment where the application permissions can be inferred from the instance's default credentials.

```bash
Usage of ./gcp_pubsub_handler:
  -alsologtostderr
        log to standard error as well as files
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -logtostderr
        log to standard error instead of files
  -project_id string
        the project id for the GCP PubSub topic.
  -stderrthreshold value
        logs at or above this threshold go to stderr
  -test_permissions
        set to test pubsub publish permissions.
  -topic string
        the project id for the GCP PubSub topic.
  -v value
        log level for V logs
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging
```