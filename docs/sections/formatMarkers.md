[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md) 

# Format Markers

Format markers allow a Yot user the ability to access the original values returned from the JSONPath query.  This can be used to add some additional text to the original value, and they can be used more than once in the overlay's `value` key.  

| Marker | Description | Action where available | status |
| --- | --- | --- | --- |
| %k | Original key's value returned from the JSONPath `query`. If querying for a key by appending the `~` character, use `%v` | combine (comments only), merge, replace (comments only), and currently only available on scalar types | experimental in (v0.1.0) |
| %v | Original value returned from the JSONPath `query` | combine, merge, replace | stable |
| %h | Original value of the head comment (comment above original value) returned from the JSONPath `query` | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |
| %l | Original value of the line comment (comment on the same line as the value) | combine (comments only), merge, replace (comments only) | stable |
| %f | Original value of the foot comment (comment below original value) | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |

[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md)  
