# Size estimates

## Storage

### Plain text data

WWM requires a single location to be able to process at least 200 patients / week. From this we can then calculate storage requirements for next five years:

`200 [patients/week] *  54 [weeks/year] * 5 [year] = 5400 [examinations]`

If we assume the worst case scenario this will result in 5400 initial screening entries and 5400 encounter logs which would then result in:

`5400 [examinations] * (XKB [average initial screening] + YKB [average encounter log]) = ZMB`

Even if not supported right now, we need to assume it will be possible to appended multimedia files to those storage entries soon after launch. If we average out the size of those pictures to 2MB and assume it they will be attached to ~70% of storage entries we come up with:

`5400 [examinations] * 70% [rate of attachment] * 2MB [average file size] = 7.4GB`

