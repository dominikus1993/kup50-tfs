import { assertEquals } from "https://deno.land/std@0.178.0/testing/asserts.ts";
import {getFirstAndLastDayInMonth} from "./index.ts";

Deno.test(function addTest() {
    const result = getFirstAndLastDayInMonth(new Date(2022, 4, 14)) // 14.05.2022
    assertEquals(result, { firstDay:  "05/01/2022", lastDay: "05/31/2022", });
});
