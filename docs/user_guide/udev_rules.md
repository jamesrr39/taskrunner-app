# Udev rules

Udev is a system used by linux to control access to devices.

`taskrunner-app` scans the udev directories to see if there are any rules configured to run using `taskrunner-app`

Currently, you can't add udev rules through the user-interface.
However they are quite straight-forward to do manually.

1. Discover the idVendor and idProduct properties of your device.

Plug your device in
Run `lsusb`. Identify the device you want to run a script when it connects.

2. Open your udev rules configuration

    sudo gedit /etc/udev/rules.d/51-taskrunner-app.rules
	
opens your text editor. Enter this line

	ATTRS{idVendor}=="[id_vendor]", ATTRS{idProduct}=="[id_product]", MODE="0666", RUN+="[/path/to/taskrunner-app] --run-job='[job_name]' --trigger='udev'"

and replace

[id_vendor] with the idVendor property from the lsusb output
[id_product] with the idProduct property from the lsusb output
[/path/to/taskrunner-app] with the full path to the taskrunner-app application
[job_name] is the name of the job! Bear in mind, if you change the job name, this isn't going to work anymore.

an example would be

	ATTRS{idVendor}=="0400", ATTRS{idProduct}=="0400", MODE="0666", RUN+="/opt/taskrunner-app --run-job='backup phone' --trigger='udev'"



3. Back in the terminal, run

	sudo service udev restart
	
so that the udev daemon loads the configuration