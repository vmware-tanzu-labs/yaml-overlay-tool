[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md) 

# Format Markers

Format markers allow a Yot user the ability to access the original values returned from the JSONPath query.  This can be used to add some additional text to the original value, and they can be used more than once in the overlay's `value` key.  

| Marker | Description | Action where available |
| --- | --- | --- |
| %v | Original value returned from the JSONPath `query` | combine, merge, replace |
| %h | Original value of the head comment (comment above original value) returned from the JSONPath `query` | combine, merge |
| %l | Original value of the line comment (comment on the same line as the value) | combine, merge |
| %f | Original value of the foot comment (comment below original value) | combine, merge |

[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md)  
