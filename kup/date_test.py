from datetime import datetime
from kup.date import get_first_day_of_month, get_first_day_of_month_when_none, get_last_day_of_month, get_last_day_of_month_when_none
import unittest

class TestDate(unittest.TestCase):
        def test_first_day(self):
                subject = get_first_day_of_month(datetime.strptime("2021-09-30", '%Y-%m-%d'))
                self.assertEqual(subject, "2021-09-01")
        def test_last_day(self):
                subject = get_last_day_of_month(datetime.strptime("2021-09-15", '%Y-%m-%d'))
                self.assertEqual(subject, "2021-09-30")

        def test_first_day_when_none(self):
                subject = get_first_day_of_month_when_none(None)
                self.assertIsNotNone(subject)
        def test_last_day_when_none(self):
                subject = get_last_day_of_month_when_none(None)
                self.assertIsNotNone(subject)

        def test_first_day_when_not_none(self):
                subject = get_first_day_of_month_when_none("2021-09-02")
                self.assertEqual(subject, "2021-09-02")
        def test_last_day_when_not_none(self):
                subject = get_last_day_of_month_when_none("2021-09-16")
                self.assertEqual(subject, "2021-09-16")