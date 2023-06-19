from typing import Iterable, Sequence


def stream_to_unicode(stream: Iterable[bytes]) -> Sequence[str] | None:
    try:
        lines = ''.join(chunk.decode("utf-8") for chunk in stream).split('\n')
        return list((line + '\n' for line in lines))
    except Exception:
        return None