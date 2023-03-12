namespace Kup50

open System

module Date = 
    let getFirstAndLastMonthDay(now: DateTime) =
        let first = DateOnly(now.Year, now.Month, 1)
        let last = first.AddMonths(1).AddDays(-1)
        struct (first, last)

    let formatFirstAndLastMonthDay(struct (first: DateOnly, last: DateOnly)) =
        let res = struct (first.ToString("MM/dd/yyyy"), last.ToString("MM/dd/yyyy"))
        res