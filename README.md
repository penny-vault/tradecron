# Tradecron

TradeCron enables market aware scheduling. It supports schedules via the standard
CRON format of: Minutes(Min) Hours(H) DayOfMonth(DoM) Month(M) DayOfWeek(DoW)
See: https://en.wikipedia.org/wiki/Cron

'*' wildcards only execute during market open hours

Additional market-aware modifiers are supported:

| Modifier    | Description                                                           |
|-------------|-----------------------------------------------------------------------|
| @open       | Run at market open; replaces Minute and Hour field. e.g., @open * * * |
| @close      | Run at market close; replaces Minute and Hour field                   |
| @weekbegin  | Run on first trading day of week; replaces DayOfMonth field           |
| @weekend    | Run on last trading day of week; replaces DayOfMonth field            |
| @monthbegin | Run at market open or timespec on first trading day of month          |
| @monthend   | Run at market close or timespec on last trading day of month          |

## Examples:

| Time Spec         | Description                              |
|-------------------|------------------------------------------|
| `*/5 * * * *`     | every 5 minutes                          |
| `@open * * 2`     | market open on tuesdays                  |
| `15 @open * * *`  | 15 minutes after market open             |
| `@weekbegin`      | market open on first trading day of week |
| `@open @monthend` | market open on last trading day of month |