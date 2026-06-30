# wow-server-ping

| **🇬🇧 English** | [🇷🇺 Русский](README.ru.md) |
| :-: | :-: |

Ping tool for WoW 335a servers.

![console usage](./images/console.png) 

It can work as a Prometheus metrics exporter and display graphics in Grafana:

![grafana usage](./images/grafana.png)

## Usage

For Windows you can find builds on the [Release page](https://github.com/egoroof/wow-server-ping/releases/latest). Open an issue if you need another OS builds.

Server configs are in `servers` folder. You can change them or add another file in same format.

Start `wow-ping.exe` with `-servers` option to choose servers config. For example `wow-ping.exe -servers x1` will load servers config from `servers/x1.json` file.

Ping tool will start collecting statistics and will print it to console periodically.

Available settings:

| Flag | Default | Description |
|---|---|---|
| `-servers` | `logon.wowcircle.me` | Servers config from `servers` folder |
| `-port` | - | Listen port for Prometheus metrics |
| `-timeout` | `1s` | Ping timeout |
| `-interval` | `500ms` | Sleep time between requests |
| `-stats-interval` | `30s` | How often stats should be printed to console |
| `-stats` | - | How many stats to display before exit |
| `-filter` | - | Regexp for filter servers by name |

## Antivirus reaction

Some antivirus software can detect malware (false positive) in downloaded Windows release and block download. You can add an exception and try to download it again. This tool doesn't have any malware. You can check source code and compile it yourself with golang. Also you can scan it with VirusTotal.
