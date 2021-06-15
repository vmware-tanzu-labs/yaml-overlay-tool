[Back to Overlay actions](overlayActions.md#3-merge)  
[Back to Overlay qualifiers](overlayQualifiers.md)  
[Back to Table of contents](../documentation.md) 

# Format markers

Format markers allow a Yot user the ability to access the original values returned from the JSONPath query.  This can be used to add some additional text to the original value, and they can be used more than once within the overlay's `value` key.  

| Marker | Description | Action where available | status |
| --- | --- | --- | --- |
| %k | Original key's value returned from the JSONPath `query`. If querying for a key by appending the `~` character, use `%v` | combine (comments only), merge, replace (comments only), and currently only available on scalar types | experimental in (v0.1.0) |
| %v | Original value returned from the JSONPath `query` | combine, merge, replace | stable |
| %h | Original value of the head comment (comment above original value) returned from the JSONPath `query` | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |
| %l | Original value of the line comment (comment on the same line as the value) | combine (comments only), merge, replace (comments only) | stable |
| %f | Original value of the foot comment (comment below original value) | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |

# Marker Manipulation
Format marker values can be further customized using sed by providing an optional {sedCommand} suffix

the sedCommand suffix supports the full feature set of sed curtious of the [Go-Sed](https://github.com/rwtodd/Go.Sed) package.

## Examples

`%v{s/hello/world/g}` would search the original value and replace the word hello with world

`%v{ r ~/myfile.yaml}` would insert the contents of myfile.yaml on a new line after the existing value

## Differences from Standard Sed
Regexps: uses the google/re2 package. Therefore, you have to use that syntax for the regular expressions. The main differences are:

command syntax	Traditional RE	Notes
s/a(bc*)d/$1/g	s/a\(bc*\)d/\1/g	Don't escape (); Use $1, $2, etc.
s/(?s).//	s/.//	If you want dot to match \n, use (?s) flag.
Go's regexps have many rich options, which you can see [here](https://github.com/google/re2/wiki/Syntax).

There are a few niceties though, such as interpreting '\t' and '\n' in replacement strings:

s/\w+\s*/\t$0\n/g
You can also escape the newline like in a typical sed, if you want.

Slightly Friendlier Syntax: the command syntax is a little more user-friendly when it comes to syntax. In a normal sed, you have to use one (and ONLY one) space between a r or w and the filename. the command syntax eats whitespace until it sees the filename.



[Back to Overlay actions](overlayActions.md#3-merge)  
[Back to Overlay qualifiers](overlayQualifiers.md)  
[Back to Table of contents](../documentation.md)  
