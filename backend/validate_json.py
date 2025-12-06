import json

file_path = r'C:\Users\16611\.gemini\antigravity\brain\01e2fd21-408f-419f-99a8-430744b866a1\apifox_swagger.json'

with open(file_path, 'r', encoding='utf-8') as f:
    # Read raw content to manually fix structure if JSON load fails
    content = f.read()

# Strategy: 
# 1. Provide a known good structure for the end of the file.
# 2. We suspect 'SubscribeRequest' is outside 'schemas'.
# We will regenerate the file content ensuring proper nesting.

# Let's try to load it as a dict first (likely to fail)
try:
    data = json.loads(content)
    print("JSON is valid! No fix needed.")
except json.JSONDecodeError as e:
    print(f"JSON Error: {e}")
    # Manual Fix Strategy:
    # The file likely has: ... }, "SubscribeRequest": ...
    # It SHOULD be: ... , "SubscribeRequest": ... (inside schemas)
    
    # Let's locate "SubscribeRequest"
    idx = content.find('"SubscribeRequest":')
    if idx != -1:
        # Look backwards for the closing brace of 'schemas'
        # It's likely the first '}' before "SubscribeRequest"
        pre_content = content[:idx]
        post_content = content[idx:]
        
        # Find last '}' in pre_content
        last_brace_idx = pre_content.rfind('}')
        if last_brace_idx != -1:
             # Remove that brace (and potentially comma/whitespace) to merge into previous object
             # We effectively replace `}` with `,` or just remove it if there's already a comma
             
             # Check if there is a comma after the brace
             segment = pre_content[last_brace_idx:]
             print(f"Segment around split: {segment}")
             
             # Heuristic: Replace the last occurrence of '}' with ',' 
             # But we need to be careful. The previous object "TeamWrapper" ends with '}'. 
             # We want: ... "TeamWrapper": { ... } , "SubscribeRequest": ...
             # Currently likely: ... "TeamWrapper": { ... } } , "SubscribeRequest": ...
             
             # Let's just rewrite the end of the file programmatically if we can't parse it.
             pass

# Alternative: Parse manually up to components -> schemas
# Since we know the file content, let's just do a string replacement to fix the nesting.
# The error is that "SubscribeRequest" starts AFTER "schemas" closes.
# We need to remove the closing brace before "SubscribeRequest".

# Search for the specific pattern we saw in view_file
# ... "name": { "type": "string" } } } , "SubscribeRequest": ...
# We want ... "name": { "type": "string" } } , "SubscribeRequest": ...

# We will read the file, find the lines, and rewrite.

lines = content.splitlines()
new_lines = []
skip_next_brace = False

for i, line in enumerate(lines):
    # Previous view showed:
    # 954:         },
    # 955:         "SubscribeRequest": {
    # It should mean line 954 was closing 'schemas'. We want to remove line 954.
    if '"SubscribeRequest": {' in line:
        # Check previous line
        if new_lines and new_lines[-1].strip() == '},':
             # This meant we closed 'schemas'. 
             # We should KEEP the comma but ensure we are inside schemas.
             pass
        elif new_lines and new_lines[-1].strip() == '}':
             # Replace '}' with ',' to continue the object
             new_lines[-1] = '            },'
    
    if '"SubscribeResponse": {' in line:
         # Ensure indentation matches schemas elements (12 spaces)
         line = '            "SubscribeResponse": {'
    
    if '"APIResponse": {' in line:
         line = '            "APIResponse": {'
         
    new_lines.append(line)

# Ensure checking for double closing braces at the end
# The last few lines should be:
#        }  <-- closes APIResponse
#    }      <-- closes schemas
# }         <-- closes components
# }         <-- closes root

# Re-assemble and try to parse
fixed_content = '\n'.join(new_lines)

# This approach is risky. Let's just use Python to construct the correct structure if we can salvage the parts.
# Plan B: Just output the corrected JSON structure by removing the specific closing brace we saw earlier.

# We saw in previous turn:
# 956:            }
# 957:        },   <-- This closes schemas (line 568 in view file)
# 958:        "SubscribeRequest": {

# We want to remove line 957 (and ensure line 956 has a comma if needed, or put comma on 957)

fixed_lines = []
lines = content.splitlines()
for i, line in enumerate(lines):
    if '"SubscribeRequest":' in line:
        # Look at previous line
        prev = lines[i-1].strip()
        if prev == '},':
             # This closed schemas? No, schemas starts at indent 12 usually in this file?
             # Let's just look at line 950-960 area
             pass
             
    # Specific patch based on line numbers from previous tool view
    # Line 950: "studentId",
    # ...
    # 953: ]
    # 954: }, 
    # 955: "SubscribeResponse": {
    
    # Wait, my previous view_file output in step 569 showed:
    # 954: },
    # 955: "SubscribeResponse": {
    # This looks like SubscribeRequest was DELETED? Or I am misreading.
    
    # In step 526 view_file, "SubscribeRequest" was at line 961.
    # In step 549 inner diff, I see I removed "SubscribeRequest" accidentally?
    # No, diff showed "- }, + }," then "SubscribeRequest"
    
    # Let's just use Python to load the partial JSON and re-dump it correctly.
    pass

# We will simply overwrite the file with a version where we manipulate the string directly to ensure correctness.
# We know the schemas are: User, RegisterRequest ... TeamWrapper
# And then SubscribeRequest, SubscribeResponse, APIResponse.
# We will check if SubscribeRequest is nested in schemas.

import re
# Regex to find:  } \s* ,? \s* "SubscribeRequest"
# And replace with: , "SubscribeRequest"

fixed_content = re.sub(r'}\s*,\s*"SubscribeRequest"', '},\n            "SubscribeRequest"', content)
# Verify if that was the issue. 
# Also check for end of file.

with open('temp_fixed.json', 'w', encoding='utf-8') as f:
    f.write(fixed_content)

try:
    json.load(open('temp_fixed.json', 'r', encoding='utf-8'))
    print("Fixed JSON is valid.")
    # If valid, overwrite original
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(fixed_content)
except Exception as e:
    print(f"Still invalid: {e}")
