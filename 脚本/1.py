
#!/usr/bin/env python
# -*- coding: utf8 -*-

import json
import sys

count = 0
with open("/tmp/os_result.txt") as f:
    text = f.readlines()
    for line in text:
        log = json.loads(line)
        if not log['status']:
            print(line)
            count += 1
    print("total check item:",len(text),"   false item:", count)

if count > 0:
    sys.exit(1)
