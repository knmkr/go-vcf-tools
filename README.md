# go-vcf-tools

## vcf-filter

```
$ cat a.vcf| vcf-filter -keep-ids ids.txt -keep-pos pos.txt
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
$ cat a.vcf| vcf-fill-rsids -bucket b142_GRCh37
```
