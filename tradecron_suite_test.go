// Copyright 2021-2024
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tradecron_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/penny-vault/tradecron"
	"github.com/rs/zerolog/log"
)

func TestPortfolio(t *testing.T) {
	log.Logger = log.Output(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tradecron Suite")
}

func GetTimezone() *time.Location {
	tz, err := time.LoadLocation("America/New_York") // New York is the reference time
	if err != nil {
		log.Panic().Err(err).Msg("could not load timezone")
	}
	return tz
}

func LoadMarketHolidays() {
	holidayList := []string{
		"2010-01-01", "2010-01-18", "2010-02-15", "2010-04-02", "2010-05-31", "2010-07-05", "2010-09-06",
		"2010-11-25", "2010-12-24", "2011-01-17", "2011-02-21", "2011-04-22", "2011-05-30", "2011-07-04",
		"2011-09-05", "2011-11-24", "2011-12-26", "2012-01-02", "2012-01-16", "2012-02-20", "2012-04-06",
		"2012-05-28", "2012-07-04", "2012-09-03", "2012-10-29", "2012-10-30", "2012-11-22", "2012-12-25",
		"2013-01-01", "2013-01-21", "2013-02-18", "2013-03-29", "2013-05-27", "2013-07-04", "2013-09-02",
		"2013-11-28", "2013-12-25", "2014-01-01", "2014-01-20", "2014-02-17", "2014-04-18", "2014-05-26",
		"2014-07-04", "2014-09-01", "2014-11-27", "2014-12-25", "2015-01-01", "2015-01-19", "2015-02-16",
		"2015-04-03", "2015-05-25", "2015-07-03", "2015-09-07", "2015-11-26", "2015-12-25", "2016-01-01",
		"2016-01-18", "2016-02-15", "2016-03-25", "2016-05-30", "2016-07-04", "2016-09-05", "2016-11-24",
		"2016-12-26", "2017-01-02", "2017-01-16", "2017-02-20", "2017-04-14", "2017-05-29", "2017-07-04",
		"2017-09-04", "2017-11-23", "2017-12-25", "2018-01-01", "2018-01-15", "2018-02-19", "2018-03-30",
		"2018-05-28", "2018-07-04", "2018-09-03", "2018-11-22", "2018-12-05", "2018-12-25", "2019-01-01",
		"2019-01-21", "2019-02-18", "2019-04-19", "2019-05-27", "2019-07-04", "2019-09-02", "2019-11-28",
		"2019-12-25", "2020-01-01", "2020-01-20", "2020-02-17", "2020-04-10", "2020-05-25", "2020-07-03",
		"2020-09-07", "2020-11-26", "2020-12-25", "2021-01-01", "2021-01-18", "2021-02-15", "2021-04-02",
		"2021-05-31", "2021-07-05", "2021-09-06", "2021-11-25", "2021-12-24", "2022-01-17", "2022-02-21",
		"2022-04-15", "2022-05-30", "2022-06-20", "2022-07-04", "2022-09-05", "2022-11-24", "2022-12-26",
	}

	holidays := make(map[int64]int, len(holidayList))
	for _, dtStr := range holidayList {
		dt, err := time.Parse("2006-01-02", dtStr)
		Expect(err).To(BeNil())
		dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, GetTimezone())
		holidays[dt.Unix()] = 0
	}

	// add early close
	dt, err := time.Parse("2006-01-02", "2022-11-25")
	Expect(err).To(BeNil())
	dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, GetTimezone())
	holidays[dt.Unix()] = 1300

	tradecron.SetMarketHolidays(holidays)
}
