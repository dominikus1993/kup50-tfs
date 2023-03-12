namespace Kup50

open System
open System.Globalization

module Date = 
    let getFirstAndLastMonthDay(now: DateTime)=
        let first: DateOnly = DateOnly(now.Year, now.Month, 1)
        let last: DateOnly = first.AddMonths(1).AddDays(-1)
        struct (first, last)

    let formatFirstAndLastMonthDay(struct (first: DateOnly, last: DateOnly)) =
        let res = struct (first.ToString(), last.ToString())
        res