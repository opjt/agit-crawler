package bootstrap

import (
	"agit-crawler/app/cron"
	"agit-crawler/app/lib"
	"agit-crawler/app/pkg"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	cron.Module,
	pkg.Module,
)
