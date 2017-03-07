package taskrunnerdal

/*
DAL layer of taskrunner-app.

Uses a filesystem structure to store job configs, run results and logs.

File structure:

- taskrunner base (default ~/.local/share/github.com/jamesrr39/taskrunner-app)
  - jobs/
    - {jobId}/ (numeric autoincrement; 1,2,3 etc)
	  - config.json (configuration file for job)
	  - lockfile.txt ("mutex" allowing only one job to be run at once. Contains Pid)
	  - workspace/ (transitive folder - cleaned out at the start of every run)
	    - script
	  - runs/
	    - {runId}/ (numeric autoincrement; 1,2,3 etc)
		  - joboutput.log (mixture of stdout and stderr - each line prefixed with which one it came from)
		  - summary.json

*/
