## The references type slice

When you create a slice, go internally is creating two separate data structures for you. The one is the array data. The
another is what we refer to as the slice, the slices of data structure that has three elements inside of it.

It has a pointer, a capacity number and a length number:

1. The pointer is a pointer over to the underlying array that represents the actual list of items.
2. The capacity is how many elements it can contain at present.
3. The length is represents how many elements currently exist inside the slice.

| Value Types | Reference Types |
|  --------   | --------------- |
|     int     |     slices      |
|    float    |     maps        |
|    string   |     channels    |
|     bool    |     pointers    |
|    structs  |     functions   |
