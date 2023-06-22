from datetime import datetime
from kup.name import random_kup_folder_name
import unittest

class TestDate(unittest.TestCase):
    def test_first_day(self):
        subject = random_kup_folder_name()
        self.assertTrue(len(subject) > 0)