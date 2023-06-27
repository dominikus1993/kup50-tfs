
import random
import string


def random_kup_folder_name() -> str:
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for _ in range(10))
    
    return f'kup_{result_str}'