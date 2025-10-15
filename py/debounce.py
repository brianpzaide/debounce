import asyncio
import time
from functools import wraps

async def _debounced_execution(delay_seconds, func, args, kwargs):
    await asyncio.sleep(delay_seconds)
    print(f"finally implemented at {time.time()}")
    
    if asyncio.iscoroutinefunction(func):
        return await func(*args, **kwargs)
    else:
        return func(*args, **kwargs) 

def debounce(delay_seconds: int = 1):
    def decorator(func):
        _pending_task: asyncio.Task | None = None
        _lock = asyncio.Lock()
        
        @wraps(func)
        async def wrapper(*args, **kwargs):
            nonlocal _pending_task
            async with _lock:
                if (_pending_task is not None) and (not _pending_task.done()):
                    print(f"_pending_task cancelled with args: {_pending_task._debounce_args}")
                    _pending_task.cancel()
                    _pending_task = None
                    
                _pending_task = asyncio.create_task(
                    _debounced_execution(delay_seconds, func, args, kwargs)
                )
                _pending_task._debounce_args = args
        return wrapper
    return decorator


@debounce(delay_seconds=2)
def greet(name = ""):
    print(f"Hi {name}, have a nice day!!")

async def main():
    print(f"Start time: {time.time():.2f}")

    asyncio.create_task(greet("Alice"))

    await asyncio.sleep(0.5)
    asyncio.create_task(greet("Bob"))

    await asyncio.sleep(0.5)
    asyncio.create_task(greet("Charlie"))

    await asyncio.sleep(2.1)
    
    print(f"End time: {time.time():.2f}")

if __name__ == '__main__':
    asyncio.run(main())
