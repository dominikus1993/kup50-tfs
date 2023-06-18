from typing import Iterable


def stream_to_unicode(stream: Iterable[bytes]) -> Iterable[str] | None:
    try:
        lines = ''.join(chunk.decode("utf-8") for chunk in stream).split('\n')
        return (line + '\n' for line in lines)
    except Exception:
        return None