from datetime import datetime
from kup.dir import path_to_html_file_name
import unittest

class TestDir(unittest.TestCase):
    def test_first_day(self):
        subject = path_to_html_file_name("test")
        self.assertTrue(len(subject) > 0)
        self.assertEqual("test.html", subject)