module github.com/dolanor/hal

go 1.14

require (
	github.com/BurntSushi/xgb v0.0.0-20200324125942-20f126ea2843
	github.com/BurntSushi/xgbutil v0.0.0-20190907113008-ad855c713046
	github.com/aquilax/go-wakatime v0.1.1
	gopkg.in/ini.v1 v1.55.0
)

replace (
	github.com/aquilax/go-wakatime v0.1.1 => /home/dolanor/go/src/github.com/aquilax/go-wakatime
	github.com/coreos/go-systemd/v22 v22.0.1-0.20200316104309-cb8b64719ae3 => /home/dolanor/go/src/github.com/coreos/go-systemd
)
