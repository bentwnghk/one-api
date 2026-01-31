package cron

import (
	"github.com/spf13/viper"
	"one-api/common/config"
	"one-api/common/logger"
	"one-api/common/scheduler"
	"one-api/model"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func InitCron() {
	if !config.IsMasterNode {
		logger.SysLog("Cron is disabled on slave node")
		return
	}

	// 添加每日統計任務
	err := scheduler.Manager.AddJob(
		"update_daily_statistics",
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 0, 30),
			)),
		gocron.NewTask(func() {
			model.UpdateStatistics(model.StatisticsUpdateTypeYesterday)
			logger.SysLog("更新昨日統計數據")
		}),
	)
	if err != nil {
		logger.SysError("Cron job error: " + err.Error())
		return
	}

	if config.UserInvoiceMonth {
		// 每月一號早上四點生成上個月的賬單數據
		err = scheduler.Manager.AddJob(
			"generate_statistics_month",
			gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(4, 0, 0))),
			gocron.NewTask(func() {
				err := model.InsertStatisticsMonth()
				if err != nil {
					logger.SysError("Generate statistics month data error:" + err.Error())
				}
			}),
		)
	}

	// 每十分鐘更新一次統計數據
	err = scheduler.Manager.AddJob(
		"update_statistics",
		gocron.DurationJob(10*time.Minute),
		gocron.NewTask(func() {
			model.UpdateStatistics(model.StatisticsUpdateTypeToDay)
			logger.SysLog("10分鐘統計數據")
		}),
	)

	// 開啟自動更新 並且設置了有效自動更新時間 同時自動更新模式不是system 則會從服務器拉取最新價格表
	autoPriceUpdatesInterval := viper.GetInt("auto_price_updates_interval")
	autoPriceUpdates := viper.GetBool("auto_price_updates")
	autoPriceUpdatesMode := viper.GetString("auto_price_updates_mode")

	if autoPriceUpdates &&
		autoPriceUpdatesInterval > 0 &&
		(autoPriceUpdatesMode == string(model.PriceUpdateModeAdd) ||
			autoPriceUpdatesMode == string(model.PriceUpdateModeOverwrite) ||
			autoPriceUpdatesMode == string(model.PriceUpdateModeUpdate)) {
		// 指定時間週期更新價格表
		err := scheduler.Manager.AddJob(
			"update_pricing_by_service",
			gocron.DurationJob(time.Duration(autoPriceUpdatesInterval)*time.Minute),
			gocron.NewTask(func() {
				err := model.UpdatePriceByPriceService()
				if err != nil {
					logger.SysError("Update Price Error: " + err.Error())
					return
				}
				logger.SysLog("Update Price Done")
			}),
		)
		if err != nil {
			logger.SysError("Cron job error: " + err.Error())
			return
		}
	}

	if err != nil {
		logger.SysError("Cron job error: " + err.Error())
		return
	}
}
