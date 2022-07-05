/**
 * Independent time operation tool to facilitate subsequent switch to dayjs
 */
import dayjs from 'dayjs';
import { format } from 'date-fns';

type ParsableDateType = Parameters<typeof dayjs>[0];

// Format Rule: https://date-fns.org/v2.28.0/docs/format
const FORMATTERS = {
  normal: 'yyyy-MM-dd',
  datetime: 'yyyy-MM-dd HH:mm:ss',
  datetimeWithoutSec: 'yyyy-MM-dd HH:mm',
  monthOnly: 'MM-dd',
};

const dayjsParseDate = (date: ParsableDateType) => {
  if (!date) return;
  if (date instanceof Date) return date;

  return dayjs(date).toDate();
};

export function formatDate(
  date: ParsableDateType = undefined,
  formatter: string | ((f: typeof FORMATTERS) => string) = FORMATTERS.normal,
): string {
  if (typeof formatter === 'function') formatter = formatter(FORMATTERS);
  const parsedDate = dayjsParseDate(date);
  return parsedDate ? format(parsedDate, formatter) : '';
}

export function formatToDateTime(
  date: ParsableDateType = undefined,
  formatter: string | ((f: typeof FORMATTERS) => string) = FORMATTERS.datetime,
) {
  return formatDate(date, formatter);
}

export const dateUtil = dayjs;
