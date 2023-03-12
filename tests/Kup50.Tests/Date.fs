module Tests

open System
open Kup50
open Xunit

[<Fact>]
let ``Test First And Last Day In Month`` () =
    let now = DateTime(2023, 3, 12)
    let struct (first, last) = Date.getFirstAndLastMonthDay(now)
    Assert.Equal(DateOnly(2023, 3, 1), first)
    Assert.Equal(DateOnly(2023, 3, 31), last)

[<Fact>]
let ``Test Format First And Last Day In Month`` () =
    let now = DateTime(2022, 5, 14)
    let struct (first, last) = Date.getFirstAndLastMonthDay(now) |> Date.formatFirstAndLastMonthDay
    Assert.Equal("05/01/2022", first)
    Assert.Equal("05/31/2022", last)