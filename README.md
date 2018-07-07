# scanbutton

server for the scan button on hp m1522nf multifunction printers. other devices may work too.
it's inspired by http://bernaerts.dyndns.org/linux/75-debian/266-debian-hp-aio-scan-to-folder-server
and a classic case of not invented here ;)

## example config

	{
		"Address":"http://192.0.2.1/hp/device/notifications.xml",
		"Sleep":"1s",
		"MaxSleep":"64s",
		"Path":"/path/to/scans/",
		"Sane":{
			"Flatbed":
				["-d", "hpaio:/net/HP_LaserJet_M1522nf_MFP?ip=192.0.2.1",
				 "--format", "png",
				 "--mode", "Color",
				 "--resolution", "300",
				 "--batch",
				 "--batch-count", "1",
				 "-l", "0",
				 "-t", "0",
				 "-x", "210",
				 "-y", "296"
				],
			"ADF":
				["-d", "hpaio:/net/HP_LaserJet_M1522nf_MFP?ip=192.0.2.1",
				 "--format", "png",
				 "--mode", "Gray",
				 "--resolution", "300",
				 "--batch",
				 "--source","ADF",
				 "-l", "0",
				 "-t", "0",
				 "-x", "210",
				 "-y", "296"
				]
		}
	}

- Address: url to notifications.xml file
- Sleep: how long to sleep between reading notifications.xml. if unreachable will get doubled in each step.
- MaxSleep: maximum value to sleep for before giving up.
- Path: where to put scanned files
- Sane: scanimage arguments for flatbed and adf scanning
	- Flatbed: flatbed arguments
	- ADF: adf arguments

