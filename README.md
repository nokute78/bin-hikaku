# bin-hikaku

A tool to compare binary files.

## Simple mode
Compare each bytes of two input files. It is similar to `cmp -l`.


### Option

|Option|Description|
|--|---|
|-r uint|Compare only uint bytes. |
|-s uint|Skip uint bytes at start of file.|

### Example
```
$ ./bin-hikaku file1 file2
0x0000000000000000 4d aa

0x00000000008bad4f 48 61
0x00000000008bad50 86 64
0x00000000008bad51 f7 62
0x00000000008bad52 0d 63
0x00000000008bad53 01 64
0x00000000008bad54 01 6f
0x00000000008bad55 0b 65

0x00000000008bb060 32 ff
0x00000000008bb061 73 ff
```


## Histogram Mode (-H)
Construct histogram. Lower value means less different.

### Option
|Option|Description|
|--|---|
|-r uint|Compare only uint bytes. |
|-s uint|Skip uint bytes at start of file.|
|-u uint|The range of values. Default is 4096.|

### Example
```
$ /bin-hikaku -H file1 file2
0x0000000000000000-0x0000000000001000   1% #
0x0000000000001000-0x0000000000002000   0% 
0x0000000000002000-0x0000000000003000   0% 
0x0000000000003000-0x0000000000004000   1% #
0x0000000000004000-0x0000000000005000   8% #
0x0000000000005000-0x0000000000006000  58% ######
0x0000000000006000-0x0000000000007000  99% ##########
0x0000000000007000-0x0000000000008000  99% ##########
0x0000000000008000-0x0000000000009000  99% ##########
0x0000000000009000-0x000000000000a000  99% ##########
```



## License

[Apache License v2.0](https://www.apache.org/licenses/LICENSE-2.0)
