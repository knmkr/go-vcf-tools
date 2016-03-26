# go-vcf-tools

## filter

```
$ cat a.vcf| go-vcf filter --keep-ids ids.txt --keep-pos pos.txt
```

## freq

```
$ cat a.vcf| go-vcf freq
```

## subset

```
$ cat a.vcf| go-vcf subset --keep-index 0
```

## to-tab

```
$ cat a.vcf| go-vcf to-tab
```

## update

```
$ cat a.vcf| go-vcf update --rs-merge-arch RsMergeArch.bcp.gz
```

## WIP: fill-rsids

```
$ cat a.vcf| go-vcf fill-rsids --bucket b142_GRCh37
```
