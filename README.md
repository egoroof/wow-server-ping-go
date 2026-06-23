# wow-server-ping

Ping tool for WoW 335a servers.

![console usage](./images/console.png) 

It can work as a Prometheus metrics exporter and display graphics in Grafana:

![grafana usage](./images/grafana.png)

## Usage

For Windows you can find binary in the [Release page](https://github.com/egoroof/wow-server-ping/releases/latest). For other OS you will need golang to compile source code (`go build wow-ping.go`).

Server configs are in `servers` folder. You can change them or add another file in same format.

Start `wow-ping` with `-servers` option to choose servers config. For example `wow-ping -servers x1` will load servers config from `servers/x1.json` file.

Ping tool will start collecting statistics and will print it to console periodically.

You can confugure many options. To see all available options run `wow-ping --help`

## Antivirus reaction

Some antivirus software can detect malware (false positive) in downloaded Windows release and block download. You can add an exception and try to download it again. This tool doesn't have any malware. You can check source code and compile it yourself with golang. Also you can scan it with VirusTotal.
