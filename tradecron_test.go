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
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/penny-vault/tradecron"
)

var _ = Describe("Tradecron", func() {
	BeforeEach(func() {
		LoadMarketHolidays()
	})

	DescribeTable("when parsing tradecron spec",
		func(spec string, hours tradecron.MarketHours, expectedTimeSpec string, expectedTimeFlag string, expectedDateFlag string, expectedError error) {
			cron, err := tradecron.New(spec, hours)
			if expectedError == nil {
				Expect(err).To(BeNil())
				Expect(cron.ScheduleString).To(Equal(spec))
				Expect(cron.TimeSpec).To(Equal(expectedTimeSpec))
				Expect(cron.TimeFlag).To(Equal(expectedTimeFlag))
				Expect(cron.DateFlag).To(Equal(expectedDateFlag))
			} else {
				Expect(err).To(Equal(expectedError))
			}
		},
		Entry("Daily every 5 minutes, regular hours", "*/5 * * * *", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes brief form, regular hours", "*/5", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes 2 of 5 fields specified, regular hours", "*/5 *", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes 3 of 5 fields specified, regular hours", "*/5 * *", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes 4 of 5 fields specified, regular hours", "*/5 * * *", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes trailing whitespace, regular hours", "*/5 ", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Daily every 5 minutes leading whitespace, regular hours", " */5", tradecron.RegularHours, "*/5 * * * *", "", "", nil),
		Entry("Malformed timespec with invalid characters", "$/5 * * * *", tradecron.RegularHours, "", "", "", errors.New("failed to parse int from $: strconv.Atoi: parsing \"$\": invalid syntax")),
		Entry("Malformed timespec with too many fields", "*/5 * * * * *", tradecron.RegularHours, "*/5 * * * *", "", "", errors.New("expected exactly 5 fields, found 6: [*/5 * * * * *]")),
		Entry("Daily 5 minutes after market open, regular hours", "@open 5 0 * * *", tradecron.RegularHours, "35 9 * * *", "@open", "", nil),
		Entry("At market open, regular hours", "@open", tradecron.RegularHours, "30 9 * * *", "@open", "", nil),
		Entry("5 min after market open brief form, regular hours", "@open 5", tradecron.RegularHours, "35 9 * * *", "@open", "", nil),
		Entry("Daily 5 minutes before market open, regular hours", "@open -5 0 * * *", tradecron.RegularHours, "25 9 * * *", "@open", "", nil),
		Entry("Daily 1 hour after market open, regular hours", "@open 0 1 * * *", tradecron.RegularHours, "30 10 * * *", "@open", "", nil),
		Entry("Daily 90 minutes after market open, regular hours", "@open 90 0 * * *", tradecron.RegularHours, "0 11 * * *", "@open", "", nil),
		Entry("Daily 1 hour before market open, regular hours", "@open 0 -1 * * *", tradecron.RegularHours, "30 8 * * *", "@open", "", nil),
		Entry("Daily 15 hours after market open, regular hours", "@open 0 15 * * *", tradecron.RegularHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 10 hours before market open, regular hours", "@open 0 -10 * * *", tradecron.RegularHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 5 minutes after market close, regular hours", "@close 5 0 * * *", tradecron.RegularHours, "5 16 * * *", "@close", "", nil),
		Entry("Daily 5 minutes before market close, regular hours", "@close -5 0 * * *", tradecron.RegularHours, "55 15 * * *", "@close", "", nil),
		Entry("Daily 1 hour after market close, regular hours", "@close 0 1 * * *", tradecron.RegularHours, "0 17 * * *", "@close", "", nil),
		Entry("Daily 1 hour before market close, regular hours", "@close 0 -1 * * *", tradecron.RegularHours, "0 15 * * *", "@close", "", nil),
		Entry("Daily 8 hours after market close, regular hours", "@close 0 8 * * *", tradecron.RegularHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 17 hours before market close, regular hours", "@close 0 -17 * * *", tradecron.RegularHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 5 minutes after market open, extended hours", "@open 5 0 * * *", tradecron.ExtendedHours, "5 7 * * *", "@open", "", nil),
		Entry("Daily 5 minutes before market open, extended hours", "@open -5 0 * * *", tradecron.ExtendedHours, "55 6 * * *", "@open", "", nil),
		Entry("Daily 1 hour after market open, extended hours", "@open 0 1 * * *", tradecron.ExtendedHours, "0 8 * * *", "@open", "", nil),
		Entry("Daily 1 hour before market open, extended hours", "@open 0 -1 * * *", tradecron.ExtendedHours, "0 6 * * *", "@open", "", nil),
		Entry("Daily 17 hours after market open, extended hours", "@open 0 17 * * *", tradecron.ExtendedHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 8 hours before market open, extended hours", "@open 0 -8 * * *", tradecron.ExtendedHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 5 minutes after market close, extended hours", "@close 5 0 * * *", tradecron.ExtendedHours, "5 20 * * *", "@close", "", nil),
		Entry("Daily 5 minutes before market close, extended hours", "@close -5 0 * * *", tradecron.ExtendedHours, "55 19 * * *", "@close", "", nil),
		Entry("Daily 1 hour after market close, extended hours", "@close 0 1 * * *", tradecron.ExtendedHours, "0 21 * * *", "@close", "", nil),
		Entry("Daily 1 hour before market close, extended hours", "@close 0 -1 * * *", tradecron.ExtendedHours, "0 19 * * *", "@close", "", nil),
		Entry("Daily 8 hours after market close, extended hours", "@close 0 8 * * *", tradecron.ExtendedHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Daily 21 hours before market close, extended hours", "@close 0 -21 * * *", tradecron.ExtendedHours, "", "", "", tradecron.ErrFieldOutOfBounds),
		Entry("Annually, regular hours", "@monthend * * * 12 *", tradecron.RegularHours, "* * * 12 *", "", "@monthend", nil),
		Entry("Both @open @close specified", "@open @close", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("Both @weekbegin @weekend specified", "@weekbegin @weekend", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("Both @weekbegin @monthbegin specified", "@weekbegin @monthbegin", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("Both @weekbegin @monthend specified", "@weekbegin @monthend", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("Both @weekend @monthbegin specified", "@weekend @monthbegin", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("Both @weekend @monthend specified", "@weekend @monthend", tradecron.RegularHours, "", "", "", tradecron.ErrConflictingModifiers),
		Entry("@weekbegin test", "@weekbegin */5", tradecron.RegularHours, "*/5 * * * *", "", "@weekbegin", nil),
		Entry("Unknown modifier", "@modifier", tradecron.RegularHours, "", "", "", tradecron.ErrUnknownModifier),
	)

	DescribeTable("when evaluating next trade day",
		func(spec string, hours tradecron.MarketHours, given time.Time, expected time.Time) {
			cron, err := tradecron.New(spec, hours)
			Expect(err).To(BeNil())
			next := cron.Next(given)
			Expect(next).To(Equal(expected))
		},
		Entry("every 5 minutes starting on saturday", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 16, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 18, 9, 30, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting on monday at market open", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 18, 9, 30, 0, 0, GetTimezone()), time.Date(2022, 7, 18, 9, 35, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting on monday at market close", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 18, 16, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 19, 9, 30, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting on monday, extended hours", "*/5 * * * *", tradecron.ExtendedHours, time.Date(2022, 7, 18, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 18, 7, 0, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting on monday at market close, extended hours", "*/5 * * * *", tradecron.ExtendedHours, time.Date(2022, 7, 18, 20, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 19, 7, 0, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting on July 4th holiday", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 4, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 5, 9, 30, 0, 0, GetTimezone())),
		Entry("every 5 minutes starting at early close", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 11, 25, 13, 0, 0, 0, GetTimezone()), time.Date(2022, 11, 28, 9, 30, 0, 0, GetTimezone())),
		Entry("Annually, regular hours", "@monthend * * * 12 *", tradecron.RegularHours, time.Date(2022, 6, 25, 13, 0, 0, 0, GetTimezone()), time.Date(2022, 12, 30, 9, 30, 0, 0, GetTimezone())),
		Entry("month begin, only 7th month of the year, 10am, regular hours", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 10, 0, 0, 0, GetTimezone()), time.Date(2023, 7, 3, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is 1 minute before required", "@monthbegin */5 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 10, 9, 0, 0, GetTimezone()), time.Date(2022, 7, 1, 10, 10, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is 1 minute before required", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 9, 59, 0, 0, GetTimezone()), time.Date(2022, 7, 1, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is 1 hour before required", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 9, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 1, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is 1 hour after required", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 11, 0, 0, 0, GetTimezone()), time.Date(2023, 7, 3, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is month begin but time is end of day", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 1, 23, 0, 0, 0, GetTimezone()), time.Date(2023, 7, 3, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin, check when for date is 1 trading day after required", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2022, 7, 4, 11, 0, 0, 0, GetTimezone()), time.Date(2023, 7, 3, 10, 0, 0, 0, GetTimezone())),
		Entry("month begin", "@monthbegin", tradecron.RegularHours, time.Date(2022, 6, 25, 13, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 1, 9, 30, 0, 0, GetTimezone())),
		Entry("month end", "@monthend", tradecron.RegularHours, time.Date(2022, 7, 1, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 29, 9, 30, 0, 0, GetTimezone())),
		Entry("month end at 10am", "@monthend 0 10 * *", tradecron.RegularHours, time.Date(2022, 7, 29, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 29, 10, 0, 0, 0, GetTimezone())),
		Entry("month end at 10am, 1 min before", "@monthend 0 10 * *", tradecron.RegularHours, time.Date(2022, 7, 29, 9, 59, 0, 0, GetTimezone()), time.Date(2022, 7, 29, 10, 0, 0, 0, GetTimezone())),
		Entry("month end at 10am, on the dot", "@monthend 0 10 * *", tradecron.RegularHours, time.Date(2022, 7, 29, 10, 0, 0, 0, GetTimezone()), time.Date(2022, 8, 31, 10, 0, 0, 0, GetTimezone())),
		Entry("month end at 10am, 1 min after", "@monthend 0 10 * *", tradecron.RegularHours, time.Date(2022, 7, 29, 10, 1, 0, 0, GetTimezone()), time.Date(2022, 8, 31, 10, 0, 0, 0, GetTimezone())),
		Entry("month end at 10am, every 5 min, 1 min after", "@monthend */5 10 * *", tradecron.RegularHours, time.Date(2022, 7, 29, 10, 1, 0, 0, GetTimezone()), time.Date(2022, 7, 29, 10, 5, 0, 0, GetTimezone())),
		Entry("month end at 10am, date after last trading day of month", "@monthend 0 10 * *", tradecron.RegularHours, time.Date(2022, 7, 30, 10, 0, 0, 0, GetTimezone()), time.Date(2022, 8, 31, 10, 0, 0, 0, GetTimezone())),
		Entry("week begin, july 4th holiday", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 4, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 5, 9, 30, 0, 0, GetTimezone())),
		Entry("week begin, regular", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 11, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 11, 9, 30, 0, 0, GetTimezone())),
		Entry("week begin, at open", " @open @weekbegin", tradecron.RegularHours, time.Date(2022, 7, 11, 9, 30, 0, 0, GetTimezone()), time.Date(2022, 7, 18, 9, 30, 0, 0, GetTimezone())),
		Entry("week begin, at open, 5 min", "@weekbegin */5", tradecron.RegularHours, time.Date(2022, 7, 11, 9, 30, 0, 0, GetTimezone()), time.Date(2022, 7, 11, 9, 35, 0, 0, GetTimezone())),
		Entry("week end", "@weekend", tradecron.RegularHours, time.Date(2022, 7, 4, 0, 0, 0, 0, GetTimezone()), time.Date(2022, 7, 8, 9, 30, 0, 0, GetTimezone())),
		Entry("week end, at open", "@open @weekend", tradecron.RegularHours, time.Date(2022, 7, 8, 9, 30, 0, 0, GetTimezone()), time.Date(2022, 7, 15, 9, 30, 0, 0, GetTimezone())),
		Entry("Jan 27 2023", "@monthbegin 0 10 * 7", tradecron.RegularHours, time.Date(2023, 1, 27, 0, 0, 0, 0, GetTimezone()), time.Date(2023, 7, 3, 10, 0, 0, 0, GetTimezone())),
	)

	DescribeTable("when evaluating IsTradeDay",
		func(spec string, hours tradecron.MarketHours, given time.Time, expected bool) {
			cron, err := tradecron.New(spec, hours)
			Expect(err).To(BeNil())
			tradeDay := cron.IsTradeDay(given)
			Expect(tradeDay).To(Equal(expected))
		},
		Entry("every 5 minutes starting on saturday", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 16, 0, 0, 0, 0, GetTimezone()), false),
		Entry("every 5 minutes starting on monday at market open", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 18, 9, 30, 0, 0, GetTimezone()), true),
		Entry("every 5 minutes starting on monday at market close", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 18, 16, 0, 0, 0, GetTimezone()), true),
		Entry("every 5 minutes starting on monday, extended hours", "*/5 * * * *", tradecron.ExtendedHours, time.Date(2022, 7, 18, 0, 0, 0, 0, GetTimezone()), true),
		Entry("every 5 minutes starting on monday at market close, extended hours", "*/5 * * * *", tradecron.ExtendedHours, time.Date(2022, 7, 18, 20, 0, 0, 0, GetTimezone()), true),
		Entry("every 5 minutes starting on July 4th holiday", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 7, 4, 0, 0, 0, 0, GetTimezone()), false),
		Entry("every 5 minutes starting at early close", "*/5 * * * *", tradecron.RegularHours, time.Date(2022, 11, 25, 13, 0, 0, 0, GetTimezone()), true),
		Entry("month begin, date given not month begin", "@monthbegin", tradecron.RegularHours, time.Date(2022, 6, 25, 13, 0, 0, 0, GetTimezone()), false),
		Entry("month begin, date given is month begin", "@monthbegin", tradecron.RegularHours, time.Date(2022, 7, 1, 13, 0, 0, 0, GetTimezone()), true),
		Entry("month end, date given not month end", "@monthend", tradecron.RegularHours, time.Date(2022, 7, 1, 0, 0, 0, 0, GetTimezone()), false),
		Entry("month end, date given is month end", "@monthend", tradecron.RegularHours, time.Date(2022, 7, 29, 9, 30, 0, 0, GetTimezone()), true),
		Entry("week begin, date given is not week begin (holiday)", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 4, 0, 0, 0, 0, GetTimezone()), false),
		Entry("week begin, date given is week begin (holiday)", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 5, 9, 30, 0, 0, GetTimezone()), true),
		Entry("week begin, date given is not week begin", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 1, 0, 0, 0, 0, GetTimezone()), false),
		Entry("week begin, date given is week begin", "@weekbegin", tradecron.RegularHours, time.Date(2022, 7, 11, 9, 30, 0, 0, GetTimezone()), true),
		Entry("week end, date given is not week end", "@weekend", tradecron.RegularHours, time.Date(2022, 7, 6, 0, 0, 0, 0, GetTimezone()), false),
		Entry("week end, date given is week end", "@weekend", tradecron.RegularHours, time.Date(2022, 7, 8, 9, 30, 0, 0, GetTimezone()), true),
	)
})
