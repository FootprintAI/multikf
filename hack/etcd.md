#### ETCD Performance

kube-controller and kube-scheduler` keep crashing as it got lease timeout due to slow etcd response.

etcd is heavily relies on the disk performance, so we did some benchmark on disk usage.

short conclusion: vm running on harvester has slow disk performance than expected. so we probably use [tmpfs for etcd](https://github.com/kubernetes-sigs/kind/issues/845) or [disable-fsync](https://github.com/etcd-io/etcd/pull/11946)


* disable-fsync

```
vim /etc/kubernetes/manifests/etcd.yaml

spec:
  containers:
  - command:
    - etcd
+    - --unsafe-no-fsync

```


reference:

[1][使用Fio测试存储性能是否满足etcd要求](https://blog.happyhack.io/2021/08/05/fio-and-etcd/)

[2]
[[KB]NVMe PCIe - Slow Virtual Machine Performance  · Issue #3356 · harvester/harvester](https://github.com/harvester/harvester/issues/3356)

```
kubectl -n kube-system get pods                               1 ↵
NAME                                     READY   STATUS    RESTARTS          AGE
coredns-bd6b6df9f-j9px5                  1/1     Running   1                 2d
coredns-bd6b6df9f-skhsv                  1/1     Running   1                 42h
etcd-dev01-master01                      1/1     Running   0                 8h
kube-apiserver-dev01-master01            1/1     Running   5 (8h ago)        41h
kube-controller-manager-dev01-master01   1/1     Running   345 (2m53s ago)   39h
kube-proxy-6qgv8                         1/1     Running   1                 2d
kube-proxy-8m7hc                         1/1     Running   2 (9h ago)        2d
kube-proxy-wb8nz                         1/1     Running   0                 8h
kube-proxy-x6nhf                         1/1     Running   1                 2d
kube-scheduler-dev01-master01            1/1     Running   336 (2m53s ago)   39h
```

experiment setup

1. running with harvester nvme with 3 replica

```
# running on a single vm
# 4vCore, 12G
ubuntu@dev01-etcd01:/tmp$ fio --rw=write --ioengine=sync --fdatasync=1 --directory=test-data --size=22m --bs=2300 --name=mytest
mytest: (g=0): rw=write, bs=(R) 2300B-2300B, (W) 2300B-2300B, (T) 2300B-2300B, ioengine=sync, iodepth=1
fio-3.16
Starting 1 process
Jobs: 1 (f=1): [W(1)][100.0%][w=13KiB/s][w=6 IOPS][eta 00m:00s]
mytest: (groupid=0, jobs=1): err= 0: pid=2589: Fri Sep  8 07:53:08 2023
  write: IOPS=58, BW=132KiB/s (135kB/s)(21.0MiB/170299msec); 0 zone resets
    clat (usec): min=2, max=167031, avg=53.34, stdev=1687.68
     lat (usec): min=2, max=167032, avg=53.64, stdev=1687.69
    clat percentiles (usec):
     |  1.00th=[    4],  5.00th=[    5], 10.00th=[    6], 20.00th=[    7],
     | 30.00th=[    8], 40.00th=[    9], 50.00th=[   10], 60.00th=[   10],
     | 70.00th=[   11], 80.00th=[   12], 90.00th=[   15], 95.00th=[   18],
     | 99.00th=[ 1352], 99.50th=[ 2343], 99.90th=[ 3654], 99.95th=[ 3785],
     | 99.99th=[ 5145]
   bw (  KiB/s): min=    4, max= 1037, per=100.00%, avg=192.10, stdev=220.93, samples=233
   iops        : min=    1, max=  462, avg=85.76, stdev=98.33, samples=233
  lat (usec)   : 4=2.22%, 10=62.42%, 20=31.81%, 50=2.16%, 100=0.03%
  lat (usec)   : 250=0.07%, 500=0.02%, 750=0.12%, 1000=0.03%
  lat (msec)   : 2=0.31%, 4=0.78%, 10=0.02%, 250=0.01%
  fsync/fdatasync/sync_file_range:
    sync (usec): min=399, max=5729.7k, avg=16924.38, stdev=135471.87
    sync percentiles (usec):
     |  1.00th=[    537],  5.00th=[    652], 10.00th=[    734],
     | 20.00th=[    971], 30.00th=[   1516], 40.00th=[   1876],
     | 50.00th=[   2147], 60.00th=[   2507], 70.00th=[   3228],
     | 80.00th=[   5932], 90.00th=[  10028], 95.00th=[  33424],
     | 99.00th=[ 299893], 99.50th=[ 624952], 99.90th=[2021655],
     | 99.95th=[2801796], 99.99th=[4244636]
  cpu          : usr=0.05%, sys=0.32%, ctx=24729, majf=0, minf=15
  IO depths    : 1=200.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,10029,0,0 short=10029,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=132KiB/s (135kB/s), 132KiB/s-132KiB/s (135kB/s-135kB/s), io=21.0MiB (23.1MB), run=170299-170299msec

Disk stats (read/write):
  vda: ios=129/25993, merge=0/11792, ticks=271/375906, in_queue=348508, util=98.08%
```

2. AWS median instance with gp2

```
ubuntu@ip-172-31-21-84:/tmp$ fio --rw=write --ioengine=sync --fdatasync=1 --directory=test-data --size=22m --bs=2300 --name=mytest
mytest: (g=0): rw=write, bs=(R) 2300B-2300B, (W) 2300B-2300B, (T) 2300B-2300B, ioengine=sync, iodepth=1
fio-3.28
Starting 1 process
mytest: Laying out IO file (1 file / 22MiB)
Jobs: 1 (f=1): [W(1)][100.0%][w=1558KiB/s][w=693 IOPS][eta 00m:00s]
mytest: (groupid=0, jobs=1): err= 0: pid=2274: Fri Sep  8 14:46:38 2023
  write: IOPS=662, BW=1488KiB/s (1523kB/s)(22.0MiB/15142msec); 0 zone resets
    clat (usec): min=3, max=2052, avg=10.37, stdev=25.16
     lat (usec): min=4, max=2053, avg=11.34, stdev=25.18
    clat percentiles (usec):
     |  1.00th=[    5],  5.00th=[    5], 10.00th=[    6], 20.00th=[    8],
     | 30.00th=[    9], 40.00th=[   10], 50.00th=[   10], 60.00th=[   11],
     | 70.00th=[   11], 80.00th=[   12], 90.00th=[   15], 95.00th=[   17],
     | 99.00th=[   25], 99.50th=[   31], 99.90th=[   38], 99.95th=[   43],
     | 99.99th=[ 1434]
   bw (  KiB/s): min= 1096, max= 1581, per=99.96%, avg=1487.93, stdev=86.36, samples=30
   iops        : min=  488, max=  704, avg=662.73, stdev=38.49, samples=30
  lat (usec)   : 4=0.04%, 10=57.65%, 20=40.65%, 50=1.63%, 100=0.01%
  lat (msec)   : 2=0.01%, 4=0.01%
  fsync/fdatasync/sync_file_range:
    sync (usec): min=451, max=24225, avg=1491.47, stdev=1036.43
    sync percentiles (usec):
     |  1.00th=[  506],  5.00th=[  545], 10.00th=[  570], 20.00th=[  611],
     | 30.00th=[  652], 40.00th=[  758], 50.00th=[ 1844], 60.00th=[ 1942],
     | 70.00th=[ 2008], 80.00th=[ 2073], 90.00th=[ 2180], 95.00th=[ 2278],
     | 99.00th=[ 5145], 99.50th=[ 5997], 99.90th=[12387], 99.95th=[15533],
     | 99.99th=[23725]
  cpu          : usr=0.92%, sys=3.29%, ctx=18518, majf=0, minf=17
  IO depths    : 1=200.0%, 2=0.0%, 4=0.0%, 8=0.0%, 16=0.0%, 32=0.0%, >=64=0.0%
     submit    : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     complete  : 0=0.0%, 4=100.0%, 8=0.0%, 16=0.0%, 32=0.0%, 64=0.0%, >=64=0.0%
     issued rwts: total=0,10029,0,0 short=10029,0,0,0 dropped=0,0,0,0
     latency   : target=0, window=0, percentile=100.00%, depth=1

Run status group 0 (all jobs):
  WRITE: bw=1488KiB/s (1523kB/s), 1488KiB/s-1488KiB/s (1523kB/s-1523kB/s), io=22.0MiB (23.1MB), run=15142-15142msec

Disk stats (read/write):
  xvda: ios=0/21336, merge=0/11113, ticks=0/17388, in_queue=17387, util=99.42%
```
