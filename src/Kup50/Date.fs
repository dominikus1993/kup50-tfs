namespace Kup50

open System
open System.Globalization

module Date = 
    let getFirstAndLastMonthDay(now: DateTime)=
        let first: DateOnly = DateOnly(now.Year, now.Month, 1)
        let last: DateOnly = first.AddMonths(1).AddDays(-1)
        struct (first, last)
    
    let format (date: DateOnly) =
        date.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture)
        
    let formatFirstAndLastMonthDay(struct (first: DateOnly, last: DateOnly)) =
        let res = struct (first.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture), last.ToString("yyyy-MM-dd", CultureInfo.InvariantCulture))
        res