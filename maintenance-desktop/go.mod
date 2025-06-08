module hobbyhub-maintenance-desktop

go 1.24.3

require github.com/rivo/tview v0.0.0-20250501113434-0c592cd31026

require (
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/driver/sqlite v1.5.7 // indirect
	gorm.io/gorm v1.30.0 // indirect
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.7.1
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
	hobbyhub-server v0.0.0-00010101000000-000000000000
)

replace hobbyhub-server => ../app-server
