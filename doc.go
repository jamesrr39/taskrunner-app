package taskrunner

/*
For running tasks.

Uses a filesystem structure to store job configs, run results and logs.

File structure:

- taskrunner base (default ~/.local/share/taskrunner)
  - jobs/
    - {jobId}/ (numeric autoincrement; 1,2,3 etc)
	  - config.json (configuration file for job)
	  - workspace/ (transitive folder)
	    - script
	  - runs/
	    - {runId}/
		  - stderr.log
		  - stdout.log
		  - summary.json

*/
