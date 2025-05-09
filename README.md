# nfu - Nextflow resource usage summarizer

`nfu` is a lightweight command-line tool designed to summarize resource usage of Nextflow pipelines.  
It parses profiling tables (e.g., `execution_trace_*.txt`) and quickly estimates the amount of resources consumed (e.g., total CPU time).

Input files are generated by Nextflow with the [`-with-trace` flag](https://www.nextflow.io/docs/latest/reports.html#trace-file).
