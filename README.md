# go-vcf-tools

## vcf-filter

```
$ cat a.vcf| vcf-filter rslist.txt
```

## vcf-freq

```
$ cat a.vcf| vcf-freq
```

## vcf-subset

```
$ cat a.vcf| vcf-subset --keep-index 0
```

## vcf-to-tab

```
$ cat a.vcf| vcf-to-tab
```

## vcf-update

```
$ cat a.vcf| vcf-update RsMergeArch.bcp.gz
```

## WIP: vcf-fill-rsids

```
$ vcf-fill-rsids -bucket head100000.bcp.gz -setup
$ cat a.vcf| vcf-fill-rsids -bucket head100000.bcp.gz
```
