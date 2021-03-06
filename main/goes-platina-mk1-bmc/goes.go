// Copyright © 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"github.com/platinasystems/go/goes"
	"github.com/platinasystems/go/goes/cmd/bang"
	"github.com/platinasystems/go/goes/cmd/boot"
	"github.com/platinasystems/go/goes/cmd/cat"
	"github.com/platinasystems/go/goes/cmd/cd"
	"github.com/platinasystems/go/goes/cmd/chmod"
	"github.com/platinasystems/go/goes/cmd/cli"
	"github.com/platinasystems/go/goes/cmd/cmdline"
	"github.com/platinasystems/go/goes/cmd/cp"
	"github.com/platinasystems/go/goes/cmd/daemons"
	"github.com/platinasystems/go/goes/cmd/dmesg"
	"github.com/platinasystems/go/goes/cmd/echo"
	"github.com/platinasystems/go/goes/cmd/eeprom"
	"github.com/platinasystems/go/goes/cmd/env"
	"github.com/platinasystems/go/goes/cmd/exec"
	"github.com/platinasystems/go/goes/cmd/exit"
	"github.com/platinasystems/go/goes/cmd/export"
	"github.com/platinasystems/go/goes/cmd/fantrayd"
	"github.com/platinasystems/go/goes/cmd/femtocom"
	"github.com/platinasystems/go/goes/cmd/fspd"
	"github.com/platinasystems/go/goes/cmd/gpio"
	"github.com/platinasystems/go/goes/cmd/hdel"
	"github.com/platinasystems/go/goes/cmd/hdelta"
	"github.com/platinasystems/go/goes/cmd/helpers"
	"github.com/platinasystems/go/goes/cmd/hexists"
	"github.com/platinasystems/go/goes/cmd/hget"
	"github.com/platinasystems/go/goes/cmd/hgetall"
	"github.com/platinasystems/go/goes/cmd/hkeys"
	"github.com/platinasystems/go/goes/cmd/hset"
	"github.com/platinasystems/go/goes/cmd/i2c"
	"github.com/platinasystems/go/goes/cmd/i2cd"
	"github.com/platinasystems/go/goes/cmd/iminfo"
	"github.com/platinasystems/go/goes/cmd/imx6d"
	"github.com/platinasystems/go/goes/cmd/insmod"
	"github.com/platinasystems/go/goes/cmd/install"
	"github.com/platinasystems/go/goes/cmd/ip"
	"github.com/platinasystems/go/goes/cmd/kexec"
	"github.com/platinasystems/go/goes/cmd/keys"
	"github.com/platinasystems/go/goes/cmd/kill"
	"github.com/platinasystems/go/goes/cmd/ln"
	"github.com/platinasystems/go/goes/cmd/log"
	"github.com/platinasystems/go/goes/cmd/ls"
	"github.com/platinasystems/go/goes/cmd/lsmod"
	"github.com/platinasystems/go/goes/cmd/mkdir"
	"github.com/platinasystems/go/goes/cmd/mknod"
	"github.com/platinasystems/go/goes/cmd/mount"
	"github.com/platinasystems/go/goes/cmd/nlcounters"
	"github.com/platinasystems/go/goes/cmd/nld"
	"github.com/platinasystems/go/goes/cmd/ping"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/bmc/diag"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/bmc/ledgpiod"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/bmc/ucd9090d"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/bmc/upgrade"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/bmc/upgraded"
	"github.com/platinasystems/go/goes/cmd/platina/mk1/toggle"
	"github.com/platinasystems/go/goes/cmd/ps"
	"github.com/platinasystems/go/goes/cmd/pwd"
	"github.com/platinasystems/go/goes/cmd/reboot"
	"github.com/platinasystems/go/goes/cmd/redisd"
	"github.com/platinasystems/go/goes/cmd/reload"
	"github.com/platinasystems/go/goes/cmd/restart"
	"github.com/platinasystems/go/goes/cmd/rm"
	"github.com/platinasystems/go/goes/cmd/rmmod"
	"github.com/platinasystems/go/goes/cmd/show_commands"
	"github.com/platinasystems/go/goes/cmd/show_packages"
	"github.com/platinasystems/go/goes/cmd/slashinit"
	"github.com/platinasystems/go/goes/cmd/sleep"
	"github.com/platinasystems/go/goes/cmd/source"
	"github.com/platinasystems/go/goes/cmd/start"
	"github.com/platinasystems/go/goes/cmd/stop"
	"github.com/platinasystems/go/goes/cmd/stty"
	"github.com/platinasystems/go/goes/cmd/subscribe"
	"github.com/platinasystems/go/goes/cmd/sync"
	"github.com/platinasystems/go/goes/cmd/telnetd"
	"github.com/platinasystems/go/goes/cmd/umount"
	"github.com/platinasystems/go/goes/cmd/uninstall"
	"github.com/platinasystems/go/goes/cmd/uptimed"
	"github.com/platinasystems/go/goes/cmd/w83795d"
	"github.com/platinasystems/go/goes/cmd/watchdog"
	"github.com/platinasystems/go/goes/cmd/wget"
	"github.com/platinasystems/go/goes/lang"
)

const (
	Name    = "goes-platina-mk1-bmc"
	Apropos = "platina's mk1 baseboard management controller"
)

func Goes() *goes.Goes {
	g := goes.New(Name, "",
		lang.Alt{
			lang.EnUS: Apropos,
		},
		lang.Alt{})
	g.Plot(helpers.New()...)
	g.Plot(cli.New()...)
	g.Plot(bang.New(),
		boot.New(),
		cat.New(),
		cd.New(),
		chmod.New(),
		cmdline.New(),
		cp.New(),
		daemons.New(),
		diag.New(),
		dmesg.New(),
		echo.New(),
		eeprom.New(),
		env.New(),
		exec.New(),
		exit.New(),
		export.New(),
		fantrayd.New(),
		femtocom.New(),
		fspd.New(),
		gpio.New(),
		hdel.New(),
		hdelta.New(),
		hexists.New(),
		hget.New(),
		hgetall.New(),
		hkeys.New(),
		hset.New(),
		i2c.New(),
		i2cd.New(),
		iminfo.New(),
		imx6d.New(),
		insmod.New(),
		install.New(),
		ip.New(),
		kexec.New(),
		keys.New(),
		kill.New(),
		ledgpiod.New(),
		ln.New(),
		log.New(),
		ls.New(),
		lsmod.New(),
		mkdir.New(),
		mknod.New(),
		mount.New(),
		nlcounters.New(),
		nld.New(),
		ping.New(),
		ps.New(),
		pwd.New(),
		reboot.New(),
		redisd.New(),
		reload.New(),
		restart.New(),
		rm.New(),
		rmmod.New(),
		show_commands.New(),
		show_packages.New(""),
		show_packages.New("show-packages"),
		show_packages.New("license"),
		show_packages.New("version"),
		slashinit.New(),
		sleep.New(),
		upgrade.New(),
		upgraded.New(),
		source.New(),
		start.New(),
		stop.New(),
		stty.New(),
		subscribe.New(),
		sync.New(),
		telnetd.New(),
		toggle.New(),
		ucd9090d.New(),
		umount.New(),
		uninstall.New(),
		uptimed.New(),
		w83795d.New(),
		watchdog.New(),
		wget.New(),
	)
	return g
}
