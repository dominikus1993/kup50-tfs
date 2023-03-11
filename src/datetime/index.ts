import { format } from "../deps.ts";

export function getFirstAndLastDayInMonth(date: Date) : { firstDay: string, lastDay: string } {
    const year = date.getFullYear()
    const month = date.getMonth()
    const firstDay = new Date(year, month, 1);
    const lastDay = new Date(year, month + 1, 0);
    return { firstDay: format(firstDay, "MM/dd/yyyy"), lastDay: format(lastDay, "MM/dd/yyyy"),}
}